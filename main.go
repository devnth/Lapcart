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

	var (
		db              *sql.DB              = config.ConnectDB()
		adminRepo       repo.AdminRepository = repo.NewAdminRepo(db)
		jwtAdminService service.JWTService   = service.NewJWTAdminService()
		authService     service.AuthService  = service.NewAuthService(adminRepo)
		adminService    service.AdminService = service.NewAdminService(adminRepo)
		authHandler     v1.AuthHandler       = v1.NewAdminHandler(jwtAdminService, authService, adminService)
		adminRoute      routes.AdminRoute    = routes.NewAdminRoute()
	)
	// creating an instance of chi r
	router := chi.NewRouter()

	// using logger to display each request

	router.Use(middleware.Logger)
	// database injection
	// routes.UserRoute(router, *controller)
	adminRoute.AdminRouter(router, authHandler)
	// router.Post("/admin/login", authHandler.AdminLogin())

	log.Println("Api is listening on port:", port)
	// Starting server
	http.ListenAndServe(":"+port, router)

}
