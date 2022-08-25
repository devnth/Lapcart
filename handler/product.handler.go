package handler

import (
	"encoding/json"
	"lapcart/common/response"
	"lapcart/model"
	"lapcart/service"
	"lapcart/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ProductHandler interface {
	AddProduct() http.HandlerFunc
	ViewProducts() http.HandlerFunc
	UpdateProduct() http.HandlerFunc
	DeleteProducts() http.HandlerFunc
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

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))

		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pagesize"))

		log.Println(page, "   ", pageSize)

		pagenation := utils.Filter{
			Page:     page,
			PageSize: pageSize,
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

func (c *productHandler) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var requestData model.UpdateProduct

		json.NewDecoder(r.Body).Decode(&requestData)

		requestData.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		err := c.productService.UpdateProduct(requestData)

		if err != nil {
			response := response.BuildErrorResponse("Failed to update product", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK!", "product has been updated successfully")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)

	}
}

func (c *productHandler) DeleteProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var data model.DeleteProduct

		data.ProductId, _ = strconv.Atoi(r.URL.Query().Get("productId"))
		data.Product_Code = r.URL.Query().Get("productCode")

		data.Deleted_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		err := c.productService.DeleteProduct(data)

		if err != nil {
			response := response.BuildErrorResponse("Failed to delete product", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK!", "product has been deleted successfully")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)

	}
}
