package v1

import (
	"lapcart/service"
	"net/http"
)

type UserHandler interface {
	AddAddress() http.HandlerFunc
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (c *userHandler) AddAddress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
