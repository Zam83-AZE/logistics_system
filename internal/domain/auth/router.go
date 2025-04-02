package auth

import (
	"html/template"
	"net/http"

	"github.com/Zam83-AZE/logistics_system/pkg/session"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// RegisterRoutes autentifikasiya marşrutlarını qeydə alır
func RegisterRoutes(router *mux.Router, db *sqlx.DB, tmpl *template.Template, sessionManager *session.Manager) {
	repo := NewPostgresRepository(db)
	service := NewAuthService(repo)
	handler := NewHandler(service, tmpl, sessionManager)

	// Login səhifəsi
	router.HandleFunc("/login", handler.LoginPage).Methods("GET")
	router.HandleFunc("/login", handler.Login).Methods("POST")

	// Logout
	router.HandleFunc("/logout", handler.Logout).Methods("GET")

	// Root path-i login-ə yönləndirir
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}).Methods("GET")
}
