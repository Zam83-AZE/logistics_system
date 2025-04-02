package session

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	sessionName      = "logistics-session"
	userIDKey        = "user_id"
	usernameKey      = "username"
	authenticatedKey = "authenticated"
)

// Manager sessiya idarəsini təmin edir
type Manager struct {
	store sessions.Store
}

var manager *Manager

// NewManager yeni sessiya meneceri yaradır
func NewManager(store sessions.Store) *Manager {
	manager = &Manager{
		store: store,
	}
	return manager
}

// GetManager mövcud sessiya menecerini qaytarır (singleton)
func GetManager() *Manager {
	return manager
}

// Middleware session məlumatlarını hər HTTP sorğusuna əlavə edir
func (m *Manager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := m.store.Get(r, sessionName)

		ctx := context.WithValue(r.Context(), "user_id", session.Values[userIDKey])
		ctx = context.WithValue(ctx, "username", session.Values[usernameKey])
		ctx = context.WithValue(ctx, "authenticated", session.Values[authenticatedKey])

		// Yeni kontekstlə davam et
		next.ServeHTTP(w, r.WithContext(ctx))

		next.ServeHTTP(w, r)
	})
}

// Login istifadəçi sessiyasını yaradır
func (m *Manager) Login(w http.ResponseWriter, r *http.Request, userID int, username string) error {
	session, _ := m.store.Get(r, sessionName)

	session.Values[userIDKey] = userID
	session.Values[usernameKey] = username
	session.Values[authenticatedKey] = true

	return session.Save(r, w)
}

// Logout istifadəçi sessiyasını silir
func (m *Manager) Logout(w http.ResponseWriter, r *http.Request) error {
	session, _ := m.store.Get(r, sessionName)

	delete(session.Values, userIDKey)
	delete(session.Values, usernameKey)
	delete(session.Values, authenticatedKey)

	return session.Save(r, w)
}

// IsAuthenticated istifadəçinin giriş etdiyini yoxlayır
func (m *Manager) IsAuthenticated(r *http.Request) bool {
	session, _ := m.store.Get(r, sessionName)

	auth, ok := session.Values[authenticatedKey].(bool)
	return ok && auth
}

// GetUserID sessiyadan istifadəçi ID-sini əldə edir
func (m *Manager) GetUserID(r *http.Request) int {
	session, _ := m.store.Get(r, sessionName)

	if id, ok := session.Values[userIDKey].(int); ok {
		return id
	}

	return 0
}

// GetUsername sessiyadan istifadəçi adını əldə edir
func (m *Manager) GetUsername(r *http.Request) string {
	session, _ := m.store.Get(r, sessionName)

	if username, ok := session.Values[usernameKey].(string); ok {
		return username
	}

	return ""
}
