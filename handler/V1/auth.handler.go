package v1

import (
	"encoding/json"
	"lapcart/common/response"
	"lapcart/model"
	"lapcart/service"
	"lapcart/utils"
	"net/http"
)

type AuthHandler interface {
	AdminLogin() http.HandlerFunc
}

type authHandler struct {
	jwtService   service.JWTService
	authService  service.AuthService
	adminService service.AdminService
}

func NewAdminHandler(
	jwtService service.JWTService,
	authService service.AuthService,
	adminService service.AdminService,
) AuthHandler {
	return &authHandler{
		jwtService:   jwtService,
		authService:  authService,
		adminService: adminService,
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
			utils.ResponseJSON(w, response)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		//getting admin values
		admin, _ := c.adminService.FindAdminByEmail(loginRequest.Email)
		token := c.jwtService.GenerateToken(admin.ID, admin.Email, "admin")
		admin.Password = ""
		admin.Token = token
		response := response.BuildResponse(true, "OK!", admin)
		utils.ResponseJSON(w, response)
	}

}
