package v1

import (
	"encoding/json"
	"lapcart/common/response"
	"lapcart/model"
	"lapcart/service"
	"lapcart/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type AuthHandler interface {
	AdminLogin() http.HandlerFunc
	UserLogin() http.HandlerFunc
	UserRegister() http.HandlerFunc
	UserRefreshToken() http.HandlerFunc
	AdminRefreshToken() http.HandlerFunc
}

type authHandler struct {
	jwtAdminService service.JWTService
	jwtUserService  service.JWTService
	authService     service.AuthService
	adminService    service.AdminService
	userService     service.UserService
	validate        *validator.Validate
}

func NewAuthHandler(
	jwtAdminService service.JWTService,
	jwtUserService service.JWTService,
	authService service.AuthService,
	adminService service.AdminService,
	userService service.UserService,
	validate *validator.Validate,

) AuthHandler {
	return &authHandler{
		jwtAdminService: jwtAdminService,
		jwtUserService:  jwtUserService,
		authService:     authService,
		adminService:    adminService,
		userService:     userService,
		validate:        validate,
	}
}

func (c *authHandler) AdminLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var loginRequest model.Admin

		json.NewDecoder(r.Body).Decode(&loginRequest)

		//verifying  admin credentials
		err := c.authService.VerifyAdminCredential(loginRequest.Email, loginRequest.Password)

		if err != nil {
			response := response.BuildErrorResponse("Failed to login", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			utils.ResponseJSON(w, response)
			return
		}

		//getting admin values
		admin, _ := c.adminService.FindAdminByEmail(loginRequest.Email)
		token := c.jwtAdminService.GenerateToken(admin.ID, admin.Email, "admin")
		admin.Password = ""
		admin.Token = token
		response := response.BuildResponse(true, "OK!", admin.Token)
		utils.ResponseJSON(w, response)
	}

}

func (c *authHandler) UserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var loginRequest model.User

		json.NewDecoder(r.Body).Decode(&loginRequest)

		//verify Admin Credentials
		err := c.authService.VerifyUserCredential(loginRequest.Email, loginRequest.Password)

		if err != nil {
			response := response.BuildErrorResponse("Failed to login", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			utils.ResponseJSON(w, response)
			return
		}

		//fetching user details
		user, _ := c.userService.FindUserByEmail(loginRequest.Email)
		token := c.jwtUserService.GenerateToken(user.ID, user.Email, "user")
		user.Password = ""
		user.Token = token
		response := response.BuildResponse(true, "OK", user.Token)
		utils.ResponseJSON(w, response)
	}
}

func (c *authHandler) UserRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var input model.UserRequest

		//fetching data
		json.NewDecoder(r.Body).Decode(&input)

		log.Println(input)

		if err := c.validate.Struct(input); err != nil {
			response := response.BuildErrorResponse("validation error", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		registerRequest := model.User{
			First_Name:   input.First_Name,
			Last_Name:    input.Last_Name,
			Password:     input.Password,
			Phone_Number: input.Phone_Number,
			Email:        input.Email,
			Created_At:   time.Now(),
		}

		err := c.userService.CreateUser(registerRequest)

		if err != nil {
			response := response.BuildErrorResponse("Failed to create user", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		user, _ := c.userService.FindUserByEmail(registerRequest.Email)
		user.Password = ""
		response := response.BuildResponse(true, "OK", user)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}

func (c *authHandler) UserRefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		autheader := r.Header.Get("Authorization")
		bearerToken := strings.Split(autheader, " ")
		token := bearerToken[1]

		refreshToken, err := c.jwtUserService.GenerateRefreshToken(token)

		if err != nil {
			response := response.BuildErrorResponse("error generating refresh token", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK!", refreshToken)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)

	}
}

func (c *authHandler) AdminRefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		autheader := r.Header.Get("Authorization")
		bearerToken := strings.Split(autheader, " ")
		token := bearerToken[1]

		refreshToken, err := c.jwtAdminService.GenerateRefreshToken(token)

		if err != nil {
			response := response.BuildErrorResponse("error generating refresh token", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.BuildResponse(true, "OK!", refreshToken)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)

	}
}
