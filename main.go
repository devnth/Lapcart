package main

import (
	"database/sql"
	"lapcart/config"
	v1 "lapcart/handler"
	m "lapcart/middleware"
	"lapcart/repo"
	"lapcart/routes"
	"lapcart/service"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/subosito/gotenv"
)

// to call functions before main functions
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
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	var (
		db              *sql.DB                 = config.ConnectDB()
		mailConfig      config.MailConfig       = config.NewMailConfig()
		validate        *validator.Validate     = validator.New()
		adminRepo       repo.AdminRepository    = repo.NewAdminRepo(db)
		userRepo        repo.UserRepository     = repo.NewUserRepo(db)
		productRepo     repo.ProductRepository  = repo.NewProductRepo(db)
		wishListRepo    repo.WishListRepository = repo.NewWishListRepo(db)
		cartRepo        repo.CartRepository     = repo.NewCartRepository(db)
		jwtAdminService service.JWTService      = service.NewJWTAdminService()
		jwtUserService  service.JWTService      = service.NewJWTUserService()
		authService     service.AuthService     = service.NewAuthService(adminRepo, userRepo)
		adminService    service.AdminService    = service.NewAdminService(adminRepo, userRepo, productRepo)
		userService     service.UserService     = service.NewUserService(userRepo, productRepo, cartRepo, adminRepo, mailConfig)
		productService  service.ProductService  = service.NewProductService(productRepo)
		wishListService service.WishListService = service.NewWishListService(wishListRepo)
		cartService     service.CartService     = service.NewCartService(cartRepo, productRepo)
		authHandler     v1.AuthHandler          = v1.NewAuthHandler(jwtAdminService,
			jwtUserService, authService,
			adminService,
			userService,
			validate)
		adminMiddleware m.Middleware         = m.NewMiddlewareAdmin(jwtAdminService)
		userMiddleware  m.Middleware         = m.NewMiddlewareUser(jwtUserService)
		adminHandler    v1.AdminHandler      = v1.NewAdminHandler(adminService, userService)
		userHandler     v1.UserHandler       = v1.NewUserHandler(userService)
		productHandler  v1.ProductHandler    = v1.NewProductHandler(productService)
		wishListHandler v1.WishListHandler   = v1.NewWishListHandler(wishListService)
		cartHandler     v1.CartHandler       = v1.NewCartHandler(cartService)
		adminRoute      routes.AdminRoute    = routes.NewAdminRoute()
		userRoute       routes.UserRoute     = routes.NewUserRoute()
		wishListRoute   routes.WishListRoute = routes.NewWishListRoute()
		cartRoute       routes.CartRoute     = routes.NewCartRoute()
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

	wishListRoute.WishListRouter(router,
		userMiddleware,
		wishListHandler)

	cartRoute.CartRouter(
		router,
		userMiddleware,
		cartHandler,
	)

	log.Println("Api is listening on port:", port)
	// Starting server
	http.ListenAndServe(":"+port, router)

}
