package dashboard

import (
	"net/http"

	"html/template"

	"github.com/Zam83-AZE/logistics_system/pkg/session"
)

// Handler dashboard HTTP sorğularını işləyir
type Handler struct {
	service        Service
	tmpl           *template.Template
	sessionManager *session.Manager
}

// NewHandler yeni dashboard işləyicisi yaradır
func NewHandler(service Service, tmpl *template.Template, sessionManager *session.Manager) *Handler {
	return &Handler{
		service:        service,
		tmpl:           tmpl,
		sessionManager: sessionManager,
	}
}

// Index dashboard ana səhifəsini göstərir
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Sessiyadan istifadəçi adını əldə et
	username := h.sessionManager.GetUsername(r)

	data, err := h.service.GetDashboardData(ctx, username)
	if err != nil {
		http.Error(w, "Dashboard məlumatları əldə edərkən xəta baş verdi", http.StatusInternalServerError)
		return
	}

	h.tmpl.ExecuteTemplate(w, "dashboard/index.html", data)
}
