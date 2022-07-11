package controllers

import (
	"database/sql"
	"devnth/models"
	"devnth/repository"

	"devnth/token"
	"devnth/utils"
	"encoding/json"
	"log"
	"net/http"
)

// var ErrMesssage models.ResponseError
// var SuccesMessage models.ResponseSuccess

type Controller struct {
	ProductRepo repository.ProductRepository
	UserRepo    repository.UserRepository
}

func (c Controller) AdminLoginIndex(w http.ResponseWriter, r *http.Request) {

}

func (c Controller) AdminLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var admin models.Admin
		var jwt models.JWT

		json.NewDecoder(r.Body).Decode(&admin)

		requestPassword := admin.Password

		// checking whether the Login Credentials is of admin
		log.Println("Checking whether Admin exists.")
		admin, err := c.UserRepo.AdminLogin(admin)

		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("Please enter correct admin details")
				w.WriteHeader(http.StatusUnauthorized)
				message := models.ResponseError{
					Status:  "error",
					Message: "invalid email adress",
				}
				utils.ResponseJSON(w, message)
				return
			} else {
				message := models.ResponseError{
					Status:  "error",
					Message: err,
				}
				utils.ResponseJSON(w, message)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

		}
		//getting hashed password from database
		dbPassword := admin.Password
		//verifying password
		passwordMatch := utils.VerifyPassword(requestPassword, dbPassword)

		if !passwordMatch {
			log.Println("Invalid Password.")
			w.WriteHeader(http.StatusUnauthorized)
			message := models.ResponseError{
				Status:  "error",
				Message: "invalid password",
			}
			json.NewEncoder(w).Encode(message)
			return
		}

		token, refresh_token := token.GenerateToken(admin.First_Name, admin.Last_Name, admin.Email, admin.Phone_Number, admin.Admin_ID)

		jwt.Token = token
		jwt.Refresh_Token = refresh_token
		message := models.ResponseSuccess{
			Status:  "success",
			Message: "Login Successfully",
			Result:  jwt.Token,
		}
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, message)
	}
}

func (c Controller) AdminProductView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// var products []models.Product
		//fetching products from database
		products, err := c.ProductRepo.ViewProduct()

		if err != nil {
			log.Println("Error in Executing Query for Product View:", err)
			w.WriteHeader(http.StatusNotImplemented)
			message := models.ResponseError{
				Status:  "error",
				Message: err,
			}
			utils.ResponseJSON(w, message)
			return
		}

		message := models.ResponseSuccess{
			Status:  "success",
			Message: "successsfully fetched product",
			Result:  &products,
		}

		utils.ResponseJSON(w, message)
	}
}

func (c Controller) AdminProductAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//declaring variable product
		var product models.Product
		//fetching data from fron end into product
		json.NewDecoder(r.Body).Decode(&product)

		//calling function to add the product
		product, err := c.ProductRepo.Addproduct(product)

		if err != nil {
			w.WriteHeader(http.StatusNotImplemented)
			message := models.ResponseError{
				Status:  "error",
				Message: "Failed to add product",
			}
			utils.ResponseJSON(w, message)
			return
		}
		log.Println("Product added.")
		message := models.ResponseSuccess{
			Status:  "success",
			Message: "Product added.",
			Result:  &product,
		}
		utils.ResponseJSON(w, message)

	}
}

func (c Controller) AdminViewUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//fetching users from database
		users, err := c.UserRepo.AdminViewUsers()

		if err != nil {
			w.WriteHeader(http.StatusNotImplemented)
			message := models.ResponseError{
				Status:  "error",
				Message: err,
			}
			utils.ResponseJSON(w, message)
			return
		}

		message := models.ResponseSuccess{
			Status:  "Succes",
			Message: "users fetched",
			Result:  &users,
		}

		utils.ResponseJSON(w, message)

	}
}

func (c Controller) BlockUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		//fetching datas from body
		json.NewDecoder(r.Body).Decode(&user)

		//checking if user exists and managing the user
		user, err := c.UserRepo.BlockUsers(user)

		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("user does not exist")
				message := models.ResponseError{
					Status:  "error",
					Message: "such users does not exist",
				}
				w.WriteHeader(http.StatusUnauthorized)
				utils.ResponseJSON(w, message)
				return
			} else {
				log.Println(err)
				message := models.ResponseError{
					Status:  "error",
					Message: err,
				}
				w.WriteHeader(http.StatusUnauthorized)
				utils.ResponseJSON(w, message)
				return
			}
		}

		//writing success response message
		message := models.ResponseSuccess{
			Status:  "success",
			Message: "successfully block/unblocked user",
			Result:  &user,
		}
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, message)
	}
}

func (c Controller) AddCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var category models.Product_Category

		// fetching datas
		json.NewDecoder(r.Body).Decode(&category)

		//adding category to database
		category, err := c.ProductRepo.AddCategory(category)

		if err != nil {
			w.WriteHeader(http.StatusNotImplemented)
			message := models.ResponseError{
				Status:  "error",
				Message: "failed to add category",
			}
			utils.ResponseJSON(w, message)
			return
		}
		log.Println("Product added.")
		message := models.ResponseSuccess{
			Status:  "success",
			Message: "category added.",
			Result:  &category,
		}
		utils.ResponseJSON(w, message)

	}
}

func (c Controller) AddBranding() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var branding models.Product_Branding

		//fetching datas
		json.NewDecoder(r.Body).Decode(&branding)
		//adding branding to database
		branding, err := c.ProductRepo.AddBranding(branding)

		if err != nil {
			message := models.ResponseError{
				Status:  "error",
				Message: err,
			}
			w.WriteHeader(http.StatusUnauthorized)
			utils.ResponseJSON(w, message)
			return
		}

		message := models.ResponseSuccess{
			Status:  "success",
			Message: "product branding added",
			Result:  &branding,
		}
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, message)
	}
}

func (c Controller) AddProcessor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var processor models.Product_Processor

		//fetching data
		json.NewDecoder(r.Body).Decode(&processor)

		processor, err := c.ProductRepo.AddProcessor(processor)

		if err != nil {
			message := models.ResponseError{
				Status:  "error",
				Message: err,
			}
			utils.ResponseJSON(w, message)
		}

		message := models.ResponseSuccess{
			Status:  "success",
			Message: "processor added",
			Result:  &processor,
		}

		utils.ResponseJSON(w, message)
	}
}

func (c Controller) ViewCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//fetching data from database
		categories, err := c.ProductRepo.ViewCategory()

		if err != nil {
			message := models.ResponseError{
				Status:  "error",
				Message: err,
			}
			utils.ResponseJSON(w, message)
		}

		message := models.ResponseSuccess{
			Status:  "success",
			Message: "view category",
			Result:  categories,
		}
		utils.ResponseJSON(w, message)
	}
}

func (c Controller) ViewProcessor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//fetching data from database
		processors, err := c.ProductRepo.ViewProcessor()

		if err != nil {
			message := models.ResponseError{
				Status:  "error",
				Message: err,
			}
			utils.ResponseJSON(w, message)
		}

		message := models.ResponseSuccess{
			Status:  "success",
			Message: "view processors",
			Result:  processors,
		}
		utils.ResponseJSON(w, message)
	}
}

func (c Controller) ViewBranding() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//fetching data from database
		brands, err := c.ProductRepo.ViewBranding()

		if err != nil {
			message := models.ResponseError{
				Status:  "error",
				Message: err,
			}
			w.WriteHeader(http.StatusUnauthorized)
			utils.ResponseJSON(w, message)
			return
		}

		message := models.ResponseSuccess{
			Status:  "success",
			Message: "view brands",
			Result:  brands,
		}
		utils.ResponseJSON(w, message)
	}
}

func (c Controller) AddColor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (c Controller) AdminLogout(w http.ResponseWriter, r *http.Request) {

}
