package auth

import (
	"net/http"

	"html/template"

	"github.com/Zam83-AZE/logistics_system/pkg/session"
)

// Handler istifadəçi autentifikasiyası HTTP sorğularını işləyir
type Handler struct {
	service        Service
	tmpl           *template.Template
	sessionManager *session.Manager
}

// NewHandler yeni autentifikasiya işləyicisi yaradır
func NewHandler(service Service, tmpl *template.Template, sessionManager *session.Manager) *Handler {
	return &Handler{
		service:        service,
		tmpl:           tmpl,
		sessionManager: sessionManager,
	}
}

// LoginPage login səhifəsini göstərir
func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// Əgər istifadəçi artıq giriş edibsə, dashboard-a yönləndir
	if h.sessionManager.IsAuthenticated(r) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}

	data := LoginForm{}
	h.tmpl.ExecuteTemplate(w, "auth/login.html", data)
}

// Login giriş əməliyyatını icra edir
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	ctx := r.Context()

	// Form məlumatlarını əldə edin
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.service.Login(ctx, username, password)
	if err != nil {
		data := LoginForm{
			Username: username,
			Error:    err.Error(),
		}
		h.tmpl.ExecuteTemplate(w, "auth/login.html", data)
		return
	}

	// Sessiyada istifadəçi məlumatlarını saxla
	err = h.sessionManager.Login(w, r, user.ID, user.Username)
	if err != nil {
		data := LoginForm{
			Username: username,
			Error:    "Giriş zamanı xəta baş verdi",
		}
		h.tmpl.ExecuteTemplate(w, "auth/login.html", data)
		return
	}

	// Dashboard səhifəsinə yönləndir
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Logout çıxış əməliyyatını icra edir
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	h.sessionManager.Logout(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
