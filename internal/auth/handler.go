package auth

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// Parse template
	tmplPath := filepath.Join("templates", "login.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Şablon yüklənə bilmədi", http.StatusInternalServerError)
		return
	}

	// Execute template
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Şablon göstərilə bilmədi", http.StatusInternalServerError)
		return
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Yalnız POST sorğuları qəbul olunur", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Sorğu formatı yanlışdır", http.StatusBadRequest)
		return
	}

	// Authenticate user
	success, err := h.service.Authenticate(req.Username, req.Password)

	// Prepare response
	response := LoginResponse{
		Success: success,
	}

	if !success {
		response.Message = "İstifadəçi adı və ya şifrə yanlışdır"
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Parse templates with inheritance
	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		"templates/dashboard.html",
	)
	if err != nil {
		http.Error(w, "Şablon yüklənə bilmədi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute template with layout
	err = tmpl.ExecuteTemplate(w, "layout.html", nil)
	if err != nil {
		http.Error(w, "Şablon göstərilə bilmədi: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
