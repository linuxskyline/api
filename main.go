package main

import (
	"apt-api/app"
	"apt-api/controllers"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8000" //localhost
	}

	router.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusOK)
		},
	)

	// Host routes
	router.HandleFunc("/api/hosts", controllers.CreateHost).Methods("POST")
	router.HandleFunc("/api/hosts", controllers.ListHosts).Methods("GET")

	// Update routes
	router.HandleFunc("/api/updates", controllers.CreateUpdate).Methods("POST")
	router.HandleFunc("/api/updates", controllers.ListUpdates).Methods("GET")

	// Account routes
	router.HandleFunc("/api/user", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")

	fmt.Println(port)

	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "content-type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	err := http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(router)) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
