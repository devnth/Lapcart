package v1

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
			response := response.BuildErrorResponse("product could not add to wishlist", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK", message)
		utils.ResponseJSON(w, response)

	}
}
