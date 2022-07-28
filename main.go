package main

import (
	"database/sql"
	"lapcart/config"
	v1 "lapcart/handler/v1"
	m "lapcart/middleware"
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
		db              *sql.DB                = config.ConnectDB()
		adminRepo       repo.AdminRepository   = repo.NewAdminRepo(db)
		userRepo        repo.UserRepository    = repo.NewUserRepo(db)
		productRepo     repo.ProductRepository = repo.NewProductRepo(db)
		jwtAdminService service.JWTService     = service.NewJWTAdminService()
		jwtUserService  service.JWTService     = service.NewJWTUserService()
		authService     service.AuthService    = service.NewAuthService(adminRepo, userRepo)
		adminService    service.AdminService   = service.NewAdminService(adminRepo, userRepo)
		userService     service.UserService    = service.NewUserService(userRepo)
		productService  service.ProductService = service.NewProductService(productRepo)
		authHandler     v1.AuthHandler         = v1.NewAuthHandler(jwtAdminService,
			jwtUserService, authService,
			adminService,
			userService)
		adminMiddleware m.Middleware      = m.NewMiddlewareAdmin(jwtAdminService)
		userMiddleware  m.Middleware      = m.NewMiddlewareUser(jwtUserService)
		adminHandler    v1.AdminHandler   = v1.NewAdminHandler(adminService, userService)
		userHandler     v1.UserHandler    = v1.NewUserHandler(userService)
		productHandler  v1.ProductHandler = v1.NewProductHandler(productService)
		adminRoute      routes.AdminRoute = routes.NewAdminRoute()
		userRoute       routes.UserRoute  = routes.NewUserRoute()
	)

	//routing
	adminRoute.AdminRouter(router,
		authHandler,
		adminHandler,
		adminMiddleware,
		productHandler)
	userRoute.UserRouter(router,
		authHandler,
		userMiddleware,
		userHandler)

	log.Println("Api is listening on port:", port)
	// Starting server
	http.ListenAndServe(":"+port, router)

}
