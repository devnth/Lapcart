package routes

import (
	v1 "lapcart/handler/v1"

	"github.com/go-chi/chi"
)

type UserRoute interface {
	UserRouter(router chi.Router, authHandler v1.AuthHandler)
}

type userRoute struct{}

func NewUserRoute() UserRoute {
	return &userRoute{}
}

func (r *userRoute) UserRouter(routes chi.Router, authHandler v1.AuthHandler) {

	routes.Post("/user/register", authHandler.UserRegister())
	routes.Post("/user/login", authHandler.UserLogin())

}
