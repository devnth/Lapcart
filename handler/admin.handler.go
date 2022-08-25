package handler

import (
	"encoding/json"
	"lapcart/common/response"
	"lapcart/model"
	"lapcart/service"
	"lapcart/utils"
	"net/http"
	"strconv"
	"time"
)

type AdminHandler interface {
	ViewUsers() http.HandlerFunc
	ManageUsers() http.HandlerFunc
	AddDiscount() http.HandlerFunc
	AddCoupon() http.HandlerFunc
	ManageOrder() http.HandlerFunc
	GetAllOrders() http.HandlerFunc
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

func (c *adminHandler) AddDiscount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var discount model.Discount

		json.NewDecoder(r.Body).Decode(&discount)

		discount.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		discount.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		err := c.adminService.AddDiscount(discount)

		if err != nil {
			response := response.BuildErrorResponse("error processing request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK!", "new discount added")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)

	}
}

func (c *adminHandler) AddCoupon() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var coupon model.Coupon

		json.NewDecoder(r.Body).Decode(&coupon)

		coupon.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		err := c.adminService.AddCoupon(coupon)

		if err != nil {
			response := response.BuildErrorResponse("error processing request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK!", "new coupon added")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}

func (c *adminHandler) ManageOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var data model.ManageOrder

		json.NewDecoder(r.Body).Decode(&data)

		data.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		err := c.adminService.ManageOrders(data)

		if err != nil {
			response := response.BuildErrorResponse("error processing request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK!", "order status updated")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}

func (c *adminHandler) GetAllOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pagesize"))

		pagenation := utils.Filter{
			Page:     page,
			PageSize: pageSize,
		}

		orders, metadata, err := c.adminService.GetAllOrders(pagenation)

		result := struct {
			Orders *[]model.GetOrders
			Meta   *utils.Metadata
		}{
			Orders: orders,
			Meta:   metadata,
		}

		if err != nil {
			response := response.BuildErrorResponse("Failed to fetch orders", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK!", result)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}
