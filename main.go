package main

import (
	"apt-api/app"
	"apt-api/controllers"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	level := os.Getenv("DEBUG_LEVEL")
	if level == "debug" {
		log.Info("Log level set to debug")
		log.SetLevel(log.DebugLevel)
	} else if level == "trace" {
		log.Info("Log level set to trace")
		log.SetLevel(log.TraceLevel)
	}

	router := mux.NewRouter()
	router.Use(app.Logging)
	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "80" //localhost
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
	router.HandleFunc("/api/hosts/{id:[0-9]+}", controllers.DeleteHost).Methods("DELETE")

	// Update routes
	router.HandleFunc("/api/updates", controllers.CreateUpdate).Methods("POST")
	router.HandleFunc("/api/updates", controllers.ListUpdates).Methods("GET")
	router.HandleFunc("/api/updates/{id:[0-9]+}", controllers.DeleteUpdate).Methods("DELETE")

	// Account routes
	router.HandleFunc("/api/user", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")

	// Because we use JWT and not cookies, cors doesn't add any security.
	//  Especially since this is designed to be a public API.
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "content-type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	err := http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(router))
	if err != nil {
		fmt.Print(err)
	}
}
