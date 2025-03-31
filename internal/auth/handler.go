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
	// This would typically check for authentication tokens
	// For simplicity, we'll just serve a simple HTML response

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Dashboard</title>
		<style>
			body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
			.message { padding: 20px; background-color: #f8f9fa; border-radius: 5px; }
		</style>
	</head>
	<body>
		<div class="message">
			<h1>Xoş gəlmisiniz!</h1>
			<p>Sistemə uğurla <strong>daxil oldunuz</strong>.</p>
		</div>
	</body>
	</html>
	`))
}
