package controllers

import (
	"database/sql"
	"devnth/models"
	"devnth/token"
	"devnth/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (c Controller) UserSignUpIndex(w http.ResponseWriter, r *http.Request) {

}

func (c Controller) UserSignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		//getting values from json
		json.NewDecoder(r.Body).Decode(&user)

		//Checking whether email exists
		log.Println("Checking whether user exists.")
		err := c.UserRepo.DoesUserExists(user)

		if err {
			w.WriteHeader(http.StatusUnauthorized)
			message := models.ResponseError{
				Status:  "error",
				Message: "user already exists.",
			}
			utils.ResponseJSON(w, message)
			return
		}

		// hashing password using bcrypt
		hashedPassword := utils.HashPassword(user.Password)

		// assigning hashed password to usermodels
		user.Password = hashedPassword

		log.Println("Hashed Password: ", user.Password)

		//creating update time
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		//Registering user to the database
		user = c.UserRepo.UserSignup(user)

		//writing userdata to the response
		log.Println("Signed Up Successfully")

		message := models.ResponseSuccess{
			Status:  "success",
			Message: "signed up successfully",
			Result:  &user,
		}

		utils.ResponseJSON(w, message)
	}
}

func (c Controller) UserLoginIndex(w http.ResponseWriter, r *http.Request) {

}

func (c Controller) UserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var founduser models.User
		var jwt models.JWT

		//Decoding from request body
		json.NewDecoder(r.Body).Decode(&founduser)

		requestPassword := founduser.Password

		log.Println("Checking whether User exists.")
		founduser, err := c.UserRepo.UserLogin(founduser)

		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				message := models.ResponseError{
					Status:  "error",
					Message: "check email",
				}
				utils.ResponseJSON(w, message)
				return
			} else {
				log.Fatal(err)
			}
		}

		if !founduser.IsActive {
			log.Println("user been blocked.")
			message := models.ResponseError{
				Status:  "error",
				Message: "user been blocked.",
			}
			w.WriteHeader(http.StatusUnauthorized)
			utils.ResponseJSON(w, message)
			return
		}

		//getting hashed password from database
		dbPassword := founduser.Password
		//verifying password
		passwordMatch := utils.VerifyPassword(requestPassword, dbPassword)

		if !passwordMatch {
			message := models.ResponseError{
				Status:  "error",
				Message: "invalid password",
			}
			w.WriteHeader(http.StatusUnauthorized)
			utils.ResponseJSON(w, message)
			return
		}

		token, refresh_token := token.GenerateToken(founduser.First_Name, founduser.Last_Name, founduser.Email, founduser.Phone_Number, founduser.User_ID)

		jwt.Token = token
		jwt.Refresh_Token = refresh_token

		message := models.ResponseSuccess{
			Status:  "success",
			Message: "login success",
			Result:  &jwt.Token,
		}
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, message)

	}
}

func (c Controller) UserHomePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// var products []models.Product
		products, err := c.ProductRepo.ViewProduct()

		if err != nil {
			log.Println("Error in Executing Query for Product View:", err)
			w.WriteHeader(http.StatusNotImplemented)
			message := models.ResponseError{
				Status:  "error",
				Message: err.Error(),
			}
			utils.ResponseJSON(w, message)
			return
		}

		log.Println("Success")
		message := models.ResponseSuccess{
			Status:  "success",
			Message: "Home page Success",
			Result:  &products,
		}
		utils.ResponseJSON(w, message)

	}
}

func (c Controller) AddToCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cart models.Cart

		//fetching data
		json.NewDecoder(r.Body).Decode(&cart)

		cart, err := c.UserRepo.AddToCart(cart)

		if err != nil {
			log.Println(err)
			message := models.ResponseError{
				Status:  "error",
				Message: err,
			}
			w.WriteHeader(http.StatusNotImplemented)
			utils.ResponseJSON(w, message)
			return
		}

		message := models.ResponseSuccess{
			Status:  "success",
			Message: "product added to cart",
			Result:  cart,
		}
		utils.ResponseJSON(w, message)
	}
}

func (c Controller) ViewCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cart models.Cart

		json.NewDecoder(r.Body).Decode(&cart)

		carts, err := c.UserRepo.ViewCart(cart)

		if err != nil {
			log.Println(err)
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
			Message: "View Cart",
			Result:  carts,
		}
		utils.ResponseJSON(w, message)
	}

}

func (c Controller) UserLogout(w http.ResponseWriter, r *http.Request) {

}

func (c Controller) HomePage(w http.ResponseWriter, r *http.Request) {

}
