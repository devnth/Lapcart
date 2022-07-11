package main

import (
	"devnth/controllers"
	"devnth/database"
	"devnth/repository"

	"devnth/routes"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/subosito/gotenv"
)

//to call functions before main functions
func init() {
	gotenv.Load()
}

func main() {

	//Loading value from env file
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	//For making log file
	file, err := os.OpenFile("Logging Details", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("Logging in File not done")
	}
	log.SetOutput(file)

	//Connecting to database
	dataBase := database.ConnectDB()

	//initializing interface in controllers
	controller := &controllers.Controller{
		ProductRepo: repository.Repository{
			DB: dataBase,
		},
		UserRepo: repository.Repository{
			DB: dataBase,
		},
	}

	// creating an instance of chi r
	router := chi.NewRouter()

	// using logger to display each request
	router.Use(middleware.Logger)

	// database injection
	routes.UserRoute(router, *controller)
	routes.AdminRoute(router, *controller)

	log.Println("Api is listening on port:", port)
	// Starting server
	http.ListenAndServe(":"+port, router)

}
