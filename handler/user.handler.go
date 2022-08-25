package handler

import (
	"encoding/json"
	"lapcart/common/response"
	"lapcart/model"
	"lapcart/service"
	"lapcart/utils"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/go-chi/chi"
	"github.com/razorpay/razorpay-go"
)

type UserHandler interface {
	AddAddress() http.HandlerFunc
	ViewAddress() http.HandlerFunc
	DeleteAddress() http.HandlerFunc
	GetAllProducts() http.HandlerFunc
	ProceedToCheckout() http.HandlerFunc
	Payment() http.HandlerFunc
	SendVerificationEmail() http.HandlerFunc
	VerifyEmail() http.HandlerFunc
	PaymentSuccess() http.HandlerFunc
	Success() http.HandlerFunc
	GetAllOrders() http.HandlerFunc
	CancelOrder() http.HandlerFunc
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
		w.Header().Add("Content-Type", "application/json")
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
		w.Header().Add("Content-Type", "application/json")
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
		w.Header().Add("Content-Type", "application/json")
		utils.ResponseJSON(w, response)
	}
}

func (c *userHandler) GetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var filter model.Filter

		user_id, _ := strconv.Atoi(r.Header.Get("user_id"))

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pagesize"))

		r.ParseForm()

		filter.Category = r.Form["category"]
		filter.Brand = r.Form["brand"]
		filter.Color = r.Form["color"]
		filter.Name = r.Form["name"]
		filter.Processor = r.Form["processor"]
		filter.ProductCode = r.Form["product_code"]
		filter.Sort.Name = r.URL.Query().Get("sortbyname")
		filter.Sort.Price = r.URL.Query().Get("sortbyprice")
		filter.Sort.Latest = r.URL.Query().Get("sortbylatest")
		priceRange := r.URL.Query().Get("range")
		if priceRange != "" {
			price := strings.Split(priceRange, "-")
			filter.PriceRange.Min, _ = strconv.Atoi(price[0])
			filter.PriceRange.Max, _ = strconv.Atoi(price[1])
		}

		pagenation := utils.Filter{
			Page:     page,
			PageSize: pageSize,
		}

		products, metadata, err := c.userService.GetAllProducts(filter, user_id, pagenation)

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
		w.Header().Add("Content-Type", "application/json")
		utils.ResponseJSON(w, response)

	}
}

func (c *userHandler) ProceedToCheckout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user_id, _ := strconv.Atoi(r.Header.Get("user_id"))

		err := c.userService.ProceedToCheckout(user_id)

		if err != nil {
			response := response.BuildErrorResponse("error processing request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK!", "ready for payment")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)

	}
}

func (c *userHandler) Payment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type PageVariables struct {
			UserID           int
			OrderIdCreated   string
			TotalPrice       float64
			Name             string
			Email            string
			Phone_Number     int
			AmountInSubUnits float64
			Coupon           uint
			OrderId          uint
		}

		var requestData model.Payment

		client := razorpay.NewClient("rzp_test_cp9c3hbN2Icv3X", "K7YZRpFEGhNg4QXeIs2gHZvA")

		requestData.User_ID, _ = strconv.Atoi(chi.URLParam(r, "user_id"))
		requestData.Coupon_Code = r.URL.Query().Get("coupon")

		paymentData, err := c.userService.ProcessingPayment(requestData)

		if err != nil {
			response := response.BuildErrorResponse("error processing request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		amountInSubUnits := paymentData.Amount * 100

		data := map[string]interface{}{
			"amount":          amountInSubUnits,
			"currency":        "INR",
			"receipt":         "some_receipt_id",
			"payment_capture": 1,
		}
		body, err := client.Order.Create(data, nil)

		if err != nil {
			response := response.BuildErrorResponse("error processing request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		val := body["id"]

		orderIDCreated := val.(string)

		pageVariables := PageVariables{
			UserID:           requestData.User_ID,
			OrderIdCreated:   orderIDCreated,
			TotalPrice:       paymentData.Amount,
			AmountInSubUnits: amountInSubUnits,
			Name:             paymentData.Full_Name,
			Email:            paymentData.Email,
			Phone_Number:     paymentData.Phone_Number,
			OrderId:          paymentData.Order_ID,
		}

		parsedTemplate, err := template.ParseFiles("template/app.html")
		parsedTemplate.Execute(w, pageVariables)

		if err != nil {
			response := response.BuildErrorResponse("error processing request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}
	}
}

func (c *userHandler) PaymentSuccess() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var data model.Payment

		// json.NewDecoder(r.Body).Decode(&data)

		data.Razorpay_payment_id = r.URL.Query().Get("payment_id")
		data.Razorpay_order_id = r.URL.Query().Get("order_id")
		data.Razorpay_signature = r.URL.Query().Get("signature")
		Id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		data.Order_ID = uint(Id)
		Id, _ = strconv.Atoi(r.URL.Query().Get("user_id"))
		data.User_ID = Id
		data.Amount, _ = strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
		data.PaymentType = "razor_pay"

		data.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		data.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		err := c.userService.AddPayment(data)

		if err != nil {
			response := response.BuildErrorResponse("error processing request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK!", "payment successful")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}

func (c *userHandler) SendVerificationEmail() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user_id, _ := strconv.Atoi(r.Header.Get("user_id"))

		email, err := c.userService.SendVerificationEmail(user_id)

		if err != nil {
			response := response.BuildErrorResponse("error processing request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK!", "please check "+*email+" for verification code")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)

	}

}

func (c *userHandler) VerifyEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user model.User

		json.NewDecoder(r.Body).Decode(&user)

		user.ID, _ = strconv.Atoi(r.Header.Get("user_id"))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		err := c.userService.VerifyEmail(user)

		if err != nil {
			response := response.BuildErrorResponse("error verifying email", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK!", "your email has been verified")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}

func (c *userHandler) GetAllOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user_ID, _ := strconv.Atoi(r.Header.Get("user_id"))

		Orders, err := c.userService.GetAllOrders(user_ID)

		if err != nil {
			response := response.BuildErrorResponse("error getting products", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK!", Orders)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)

	}
}

func (c *userHandler) CancelOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		order_id, _ := strconv.Atoi(r.URL.Query().Get("order_id"))
		user_id, _ := strconv.Atoi(r.Header.Get("user_id"))

		err := c.userService.CancelOrder(order_id, user_id)

		if err != nil {
			response := response.BuildErrorResponse("error cancelling orders", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK!", "order cancelled successfully")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)

	}
}

func (c *userHandler) Success() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {

		parsedTemplate, _ := template.ParseFiles("template/success.html")
		parsedTemplate.Execute(w, nil)

	}
}
