package v1

// type AdminHandler interface {
// 	AdminLogin() http.HandlerFunc
// }

// type adminHandler struct {
// 	jwtService  service.JWTService
// 	authService service.AuthService
// }

// func NewAdminHandler(
// 	jwtService service.JWTService,
// 	authService service.AuthService,
// ) AdminHandler {
// 	return &adminHandler{
// 		jwtService:  jwtService,
// 		authService: authService,
// 	}
// }

// func (c *adminHandler) AdminLogin() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		var admin model.Admin

// 		json.NewDecoder(r.Body).Decode(&admin)

// 		requestPassword := admin.Password

// 		//checking whether admin exists
// 		adminResponse, err := c.authService.VerifyAdminCredential(admin.Email, )

// 		if err != nil {
// 			if err == sql.ErrNoRows {
// 				response := response.BuildErrorResponse("Please enter correct admin details", err.Error(), nil)
// 				utils.ResponseJSON(w, response)
// 				w.WriteHeader(http.StatusUnauthorized)
// 				return
// 			}
// 		} else {
// 			response := response.BuildErrorResponse("Failed to process request", err.Error(), nil)
// 			utils.ResponseJSON(w, response)
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		//getting hashed password from database
// 		dbPassword := admin.Password

// 		//verifying password
// 		passwordMatch := VerifyPassword(requestPassword, dbPassword)
// 	}

// }
