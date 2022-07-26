package routes

import (
	v1 "lapcart/handler/v1"

	"github.com/go-chi/chi"
)

type AdminRoute interface {
	AdminRouter(routes chi.Router, authHandler v1.AuthHandler)
}

type adminRoute struct{}

func NewAdminRoute() AdminRoute {
	return &adminRoute{}
}

// to handle admin routes
func (r adminRoute) AdminRouter(routes chi.Router, authHandler v1.AuthHandler) {
	// routes.Use(middleware.Logger)

	routes.Post("/admin/login", authHandler.AdminLogin())
	// routes.Get("/admin/logout", Controller.AdminLogout)
	// routes.Group(func(r chi.Router) {
	// 	r.Use(middleware.AdminTokenVerifyMiddleware)
	// 	r.Post("/admin/add/product", Controller.AdminProductAdd())
	// 	r.Get("/admin/view/product/page-{page}", Controller.AdminProductView())
	// 	r.Get("/admin/view/users", Controller.AdminViewUsers())
	// 	r.Post("/admin/block/users", Controller.BlockUsers())
	// 	r.Post("/admin/add/category", Controller.AddCategory())
	// 	r.Post("/admin/add/brand", Controller.AddBranding())
	// 	r.Post("/admin/add/processor", Controller.AddProcessor())
	// 	r.Get("/admin/view/brands", Controller.ViewBranding())
	// 	r.Get("/admin/view/category", Controller.ViewCategory())
	// 	r.Get("/admin/view/processor", Controller.ViewProcessor())
	// })

}
