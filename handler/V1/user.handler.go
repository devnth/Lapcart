package v1

import (
	"encoding/json"
	"lapcart/common/response"
	"lapcart/model"
	"lapcart/service"
	"lapcart/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type UserHandler interface {
	AddAddress() http.HandlerFunc
	ViewAddress() http.HandlerFunc
	DeleteAddress() http.HandlerFunc
	GetAllProductUser() http.HandlerFunc
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

		var address model.Address

		//getting user id
		address.User_id, _ = strconv.Atoi(r.Header.Get("user_id"))

		json.NewDecoder(r.Body).Decode(&address)

		err := c.userService.AddAddress(address)

		if err != nil {
			response := response.BuildErrorResponse("address not added", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK", "Address added successfully")
		utils.ResponseJSON(w, response)

	}
}

func (c *userHandler) ViewAddress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//fetching user_id
		user_id, _ := strconv.Atoi(r.Header.Get("user_id"))

		address, err := c.userService.GetAddressByUserID(user_id)

		if err != nil {
			response := response.BuildErrorResponse("unable to fetch address", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK", address)
		utils.ResponseJSON(w, response)

	}
}

func (c *userHandler) DeleteAddress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user_id, _ := strconv.Atoi(r.Header.Get("user_id"))
		address_id, _ := strconv.Atoi(chi.URLParam(r, "addressid"))

		err := c.userService.DeleteAddress(user_id, address_id)

		if err != nil {
			response := response.BuildErrorResponse("could not make the request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK", "address deleted")
		utils.ResponseJSON(w, response)
	}
}

func (c *userHandler) GetAllProductUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user_id, _ := strconv.Atoi(r.Header.Get("user_id"))
		page, _ := strconv.Atoi(chi.URLParam(r, "page"))

		pagenation := utils.Filter{
			Page:     page,
			PageSize: 3,
		}

		products, metadata, err := c.userService.GetAllProductsUser(user_id, pagenation)

		result := struct {
			Products *[]model.GetProduct
			Meta     *utils.Metadata
		}{
			Products: products,
			Meta:     metadata,
		}

		if err != nil {

			response := response.BuildErrorResponse("could not process the request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK", result)
		utils.ResponseJSON(w, response)

	}
}
