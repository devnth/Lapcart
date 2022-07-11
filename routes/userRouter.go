package routes

import (
	"devnth/controllers"
	"devnth/middleware"

	"github.com/go-chi/chi"
)

// var Controller controllers.Controller

func UserRoute(routes chi.Router, Controller controllers.Controller) {
	routes.Get("/user/signup", Controller.UserSignUpIndex)
	routes.Post("/user/signup", Controller.UserSignUp())
	routes.Get("/user/login", Controller.UserLoginIndex)
	routes.Post("/user/login", Controller.UserLogin())
	routes.Get("/user/logout", Controller.UserLogout)
	routes.Get("/", Controller.HomePage)

	routes.Group(func(r chi.Router) {
		r.Use(middleware.TokenVerifyMiddleware)
		r.Get("/homepage", Controller.UserHomePage())
		r.Post("/user/add/cart", Controller.AddToCart())
		r.Post("/user/view/cart", Controller.ViewCart())
	})

}
