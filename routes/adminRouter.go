package routes

import (
	"devnth/controllers"
	"devnth/middleware"

	"github.com/go-chi/chi"
)

var Controller controllers.Controller

// to handle admin routes
func AdminRoute(routes chi.Router, Controller controllers.Controller) {

	routes.Get("/admin/login", Controller.AdminLoginIndex)
	routes.Post("/admin/login", Controller.AdminLogin())
	routes.Get("/admin/logout", Controller.AdminLogout)

	routes.Group(func(r chi.Router) {
		r.Use(middleware.TokenVerifyMiddleware)
		r.Post("/admin/add/product", Controller.AdminProductAdd())
		r.Get("/admin/view/product", Controller.AdminProductView())
		r.Get("/admin/view/users", Controller.AdminViewUsers())
		r.Post("/admin/block/users", Controller.BlockUsers())
		r.Post("/admin/add/category", Controller.AddCategory())
		r.Post("/admin/add/brand", Controller.AddBranding())
		r.Post("/admin/add/processor", Controller.AddProcessor())
		r.Get("/admin/view/brands", Controller.ViewBranding())
		r.Get("/admin/view/category", Controller.ViewCategory())
		r.Get("/admin/view/processor", Controller.ViewProcessor())
	})

}
