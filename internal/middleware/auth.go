package middleware

import (
	"net/http"

	"github.com/Zam83-AZE/logistics_system/pkg/session"
)

// RequireAuth istifadəçi girişini tələb edən middleware
func RequireAuth(sessionManager *session.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// İstifadəçi girişini yoxla
			if !sessionManager.IsAuthenticated(r) {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
