package routes

import (
	v1 "lapcart/handler/v1"
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
	routes.Group(func(r chi.Router) {
		r.Use(middleware.AuthorizeJwt)
		r.Post("/user/add/address", userHandler.AddAddress())
		r.Get("/user/view/address", userHandler.ViewAddress())
		r.Delete("/user/delete/address-id-{addressid}", userHandler.DeleteAddress())
		r.Get("/homepage/page-{page}", userHandler.GetAllProductUser())
		r.Post("/user/filter/page-{page}", userHandler.SearchByFilter())
		r.Post("/user/proceedtocheckout", userHandler.ProceedToCheckout())
		r.Post("/user/payment", userHandler.Payment())
	})

}
