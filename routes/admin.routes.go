package routes

import (
	v1 "lapcart/handler"
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
		r.Post("/admin/product", productHandler.AddProduct())
		r.Patch("/admin/product", productHandler.UpdateProduct())
		r.Delete("/admin/product", productHandler.DeleteProducts())
		r.Get("/admin/product", productHandler.ViewProducts())
		r.Post("/admin/add/discount", adminHandler.AddDiscount())
		r.Post("/admin/add/coupon", adminHandler.AddCoupon())
		r.Post("/admin/manage/order", adminHandler.ManageOrder())
		r.Get("/admin/refresh/token", authHandler.AdminRefreshToken())
		r.Get("/admin/manage/order", adminHandler.GetAllOrders())

	})

}
