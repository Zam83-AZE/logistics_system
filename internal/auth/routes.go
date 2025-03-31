package auth

import (
	"github.com/gorilla/mux"
)

// RegisterRoutes sets up the authentication related routes
func RegisterRoutes(router *mux.Router) {
	// Create layers following Clean Architecture
	service := NewAuthService()
	handler := NewAuthHandler(service)

	// Register routes - note the different syntax
	router.HandleFunc("/", handler.LoginPage).Methods("GET")
	router.HandleFunc("/api/login", handler.Login).Methods("POST")
	router.HandleFunc("/dashboard", handler.Dashboard).Methods("GET")
}
