package v1

import (
	"encoding/json"
	"lapcart/common/response"
	"lapcart/model"
	"lapcart/service"
	"lapcart/utils"
	"net/http"
)

type AdminHandler interface {
	ViewUsers() http.HandlerFunc
	ManageUsers() http.HandlerFunc
}

type adminHandler struct {
	adminService service.AdminService
	userService  service.UserService
}

func NewAdminHandler(
	adminService service.AdminService,
	userService service.UserService,
) AdminHandler {
	return &adminHandler{
		adminService: adminService,
		userService:  userService,
	}
}

func (c *adminHandler) ViewUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {

		users, err := c.adminService.AllUsers()

		if err != nil {
			response := response.BuildErrorResponse("error getting users from database", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK", users)
		utils.ResponseJSON(w, response)

	}
}

func (c *adminHandler) ManageUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var requestData model.User

		json.NewDecoder(r.Body).Decode(&requestData)

		err := c.adminService.ManageUsers(requestData.Email, requestData.IsActive)

		if err != nil {
			response := response.BuildErrorResponse("error managing users", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		user, _ := c.userService.FindUserByEmail(requestData.Email)
		user.Password = ""

		response := response.BuildResponse(true, "OK!", user)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)

	}
}
