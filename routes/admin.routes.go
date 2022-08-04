package routes

import (
	v1 "lapcart/handler/v1"
	m "lapcart/middleware"

	"github.com/go-chi/chi"
)

type AdminRoute interface {
	AdminRouter(routes chi.Router,
		authHandler v1.AuthHandler,
		adminHandler v1.AdminHandler,
		middleware m.Middleware,
		productHandler v1.ProductHandler)
}

type adminRoute struct{}

func NewAdminRoute() AdminRoute {
	return &adminRoute{}
}

// to handle admin routes
func (r *adminRoute) AdminRouter(routes chi.Router,
	authHandler v1.AuthHandler,
	adminHandler v1.AdminHandler,
	middleware m.Middleware,
	productHandler v1.ProductHandler) {

	routes.Post("/admin/login", authHandler.AdminLogin())

	routes.Group(func(r chi.Router) {
		r.Use(middleware.AuthorizeJwt)

		r.Get("/admin/view/users", adminHandler.ViewUsers())
		r.Put("/admin/block/users", adminHandler.ManageUsers())
		r.Post("/admin/add/product", productHandler.AddProduct())
		r.Get("/admin/view/product/page-{page}", productHandler.ViewProducts())
		r.Post("/admin/add/discount", adminHandler.AddDiscount())
		r.Post("/admin/add/coupon", adminHandler.AddCoupon())
	})

}