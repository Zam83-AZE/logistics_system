package dashboard

import (
	"html/template"

	"github.com/Zam83-AZE/logistics_system/pkg/session"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// RegisterRoutes dashboard marşrutlarını qeydə alır
func RegisterRoutes(router *mux.Router, db *sqlx.DB, tmpl *template.Template) {
	sessionManager := session.GetManager() // Singleton pattern ilə əldə et

	repo := NewPostgresRepository(db)
	service := NewDashboardService(repo)
	handler := NewHandler(service, tmpl, sessionManager)

	// Dashboard ana səhifəsi
	router.HandleFunc("/dashboard", handler.Index).Methods("GET")
}
