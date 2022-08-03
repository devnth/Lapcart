package v1

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

type CartHandler interface {
	AddToCart() http.HandlerFunc
}

type cartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) CartHandler {
	return &cartHandler{
		cartService: cartService,
	}
}

func (c *cartHandler) AddToCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var cart model.Cart

		json.NewDecoder(r.Body).Decode(&cart)

		cart.User_Id, _ = strconv.Atoi(r.Header.Get("user_id"))

		cart.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		cart.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		message, err := c.cartService.AddToCart(cart)

		if err != nil {
			response := response.BuildErrorResponse("Failed to add to cart", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK!", message)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)

	}
}
