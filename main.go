package main

import (
	"avaloqpoc/internal/api"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Serve OpenAPI Document
	r.HandleFunc("/api/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./api/openapi.yaml")
	})

	// Serve Swagger UI
	r.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./third_party/swaggerui/"))))
	r.HandleFunc("/api/v1/execute", api.ExecuteCommandHandler).Methods("POST")
	r.HandleFunc("/api/login", api.Login).Methods("POST")
	r.HandleFunc("/api/whoami", api.WhoAmI).Methods("GET")
	// Redirect root to Swagger UI
	
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swaggerui/index.html?url=/api/openapi.yaml", http.StatusSeeOther)
	})

	fmt.Println("My App is Starting")
	http.HandleFunc("/api/health", Health)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Application Health OK!")
}
