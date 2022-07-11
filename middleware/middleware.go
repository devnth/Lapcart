package middleware

import (
	"devnth/models"
	"devnth/token"
	"devnth/utils"
	"log"
	"net/http"
	"strings"
)

func TokenVerifyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//getting token from header
		autheader := r.Header.Get("Authorization")
		bearerToken := strings.Split(autheader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			ok, err := token.ValidateToken(authToken)

			if err != nil {
				message := models.ResponseError{
					Status:  "error",
					Message: err,
				}
				w.WriteHeader(http.StatusUnauthorized)
				utils.ResponseJSON(w, message)
				return
			}

			if !ok {
				log.Println("Token Invalid.")
				message := models.ResponseError{
					Status:  "error",
					Message: "Token Invalid",
				}
				w.WriteHeader(http.StatusUnauthorized)
				utils.ResponseJSON(w, message)
				return

			} else if ok {
				next.ServeHTTP(w, r)
			}

		} else {
			log.Println("No token")
			w.WriteHeader(http.StatusUnauthorized)
			message := models.ResponseError{
				Status:  "error",
				Message: "No token",
			}
			utils.ResponseJSON(w, message)
		}
	})
}
