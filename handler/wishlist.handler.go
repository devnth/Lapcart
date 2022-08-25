package handler

import (
	"encoding/json"
	"lapcart/common/response"
	"lapcart/model"
	"lapcart/service"
	"lapcart/utils"
	"net/http"
	"strconv"
)

type WishListHandler interface {
	AddOrDeleteWishList() http.HandlerFunc
	ViewWishList() http.HandlerFunc
}

type wishListHandler struct {
	wishListService service.WishListService
}

func NewWishListHandler(
	wishListService service.WishListService) WishListHandler {

	return &wishListHandler{
		wishListService: wishListService,
	}
}

func (c *wishListHandler) AddOrDeleteWishList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var userRequest model.WishList

		json.NewDecoder(r.Body).Decode(&userRequest)

		userRequest.User_Id, _ = strconv.Atoi(r.Header.Get("user_id"))

		message, err := c.wishListService.AddOrDeleteWishList(userRequest)

		if err != nil {
			response := response.BuildErrorResponse("product not add/deleted to wishlist", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK", message)
		utils.ResponseJSON(w, response)

	}
}

func (c *wishListHandler) ViewWishList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user_id, _ := strconv.Atoi(r.Header.Get("user_id"))

		products, err := c.wishListService.GetWishList(user_id)

		if err != nil {
			response := response.BuildErrorResponse("could not reach the request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}
		response := response.BuildResponse(true, "OK", products)
		utils.ResponseJSON(w, response)

	}
}
