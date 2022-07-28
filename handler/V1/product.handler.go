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

type ProductHandler interface {
	AddProduct() http.HandlerFunc
	ViewProducts() http.HandlerFunc
}

type productHandler struct {
	productService service.ProductService
}

func NewProductHandler(
	productService service.ProductService,
) ProductHandler {
	return &productHandler{
		productService: productService,
	}
}

func (c *productHandler) AddProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var requestData model.Product

		json.NewDecoder(r.Body).Decode(&requestData)

		_, err := c.productService.AddProduct(requestData)

		if err != nil {
			response := response.BuildErrorResponse("Failed to add product", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		products, _ := c.productService.GetProductByCode(requestData.Code)
		response := response.BuildResponse(true, "OK!", products)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)

	}
}

func (c *productHandler) ViewProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		page, _ := strconv.Atoi(chi.URLParam(r, "page"))

		pagenation := utils.Filter{
			Page:     page,
			PageSize: 3,
		}

		products, metadata, err := c.productService.ViewProducts(pagenation)

		result := struct {
			Products *[]model.ProductResponse
			Meta     *utils.Metadata
		}{
			Products: products,
			Meta:     metadata,
		}

		if err != nil {
			response := response.BuildErrorResponse("Failed to fetch product", err.Error(), nil)
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
