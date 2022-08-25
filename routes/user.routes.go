package routes

import (
	v1 "lapcart/handler"
	m "lapcart/middleware"

	"github.com/go-chi/chi"
)

type UserRoute interface {
	UserRouter(router chi.Router,
		authHandler v1.AuthHandler,
		middleware m.Middleware,
		userHandler v1.UserHandler,
	)
}

type userRoute struct{}

func NewUserRoute() UserRoute {
	return &userRoute{}
}

func (r *userRoute) UserRouter(routes chi.Router,
	authHandler v1.AuthHandler,
	middleware m.Middleware,
	userHandler v1.UserHandler) {

	routes.Post("/user/register", authHandler.UserRegister())
	routes.Post("/user/login", authHandler.UserLogin())
	routes.Get("/user/payment/{user_id}", userHandler.Payment())
	routes.Get("/payment-success", userHandler.PaymentSuccess())
	routes.Get("/success", userHandler.Success())
	routes.Group(func(r chi.Router) {
		r.Use(middleware.AuthorizeJwt)
		r.Get("/user/email/verification", userHandler.SendVerificationEmail())
		r.Post("/user/email/verification", userHandler.VerifyEmail())
		r.Post("/user/add/address", userHandler.AddAddress())
		r.Get("/user/view/address", userHandler.ViewAddress())
		r.Delete("/user/delete/address-id-{addressid}", userHandler.DeleteAddress())
		r.Get("/user/products", userHandler.GetAllProducts())
		r.Post("/user/proceedtocheckout", userHandler.ProceedToCheckout())
		r.Get("/user/orders", userHandler.GetAllOrders())
		r.Patch("/user/orders", userHandler.CancelOrder())
		r.Get("/user/refresh/token", authHandler.UserRefreshToken())
	})

}
