package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter creates and configures a new HTTP router
func NewRouter() *mux.Router {
	// Create a new router instead of ServeMux
	router := mux.NewRouter()

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return router
}
