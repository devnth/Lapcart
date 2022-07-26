package main

import (
	"database/sql"
	"lapcart/config"
	v1 "lapcart/handler/v1"
	"lapcart/repo"
	"lapcart/routes"
	"lapcart/service"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/subosito/gotenv"
)

//to call functions before main functions
func init() {
	gotenv.Load()
}

func main() {

	//Loading value from env file
	port := os.Getenv("PORT")

	//For making log file
	file, err := os.OpenFile("Logging Details", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("Logging in File not done")
	}
	log.SetOutput(file)

	// creating an instance of chi r
	router := chi.NewRouter()

	// using logger to display each request

	router.Use(middleware.Logger)

	var (
		db              *sql.DB              = config.ConnectDB()
		adminRepo       repo.AdminRepository = repo.NewAdminRepo(db)
		userRepo        repo.UserRepository  = repo.NewUserRepo(db)
		jwtAdminService service.JWTService   = service.NewJWTAdminService()
		jwtUserService  service.JWTService   = service.NewJWTUserService()
		authService     service.AuthService  = service.NewAuthService(adminRepo, userRepo)
		adminService    service.AdminService = service.NewAdminService(adminRepo)
		userService     service.UserService  = service.NewUserService(userRepo)
		authHandler     v1.AuthHandler       = v1.NewAdminHandler(jwtAdminService,
			jwtUserService, authService,
			adminService,
			userService)
		adminRoute routes.AdminRoute = routes.NewAdminRoute()
		userRoute  routes.UserRoute  = routes.NewUserRoute()
	)

	//routing
	adminRoute.AdminRouter(router, authHandler)
	userRoute.UserRouter(router, authHandler)

	log.Println("Api is listening on port:", port)
	// Starting server
	http.ListenAndServe(":"+port, router)

}
