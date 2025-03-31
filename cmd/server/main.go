package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Zam83-AZE/logistics_system/internal/auth"
	"github.com/Zam83-AZE/logistics_system/internal/server"
)

func main() {
	// Initialize router and server
	router := server.NewRouter()

	// Register routes
	auth.RegisterRoutes(router)

	// Start server
	port := ":9099"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
