package auth

// AuthService implements the business logic for authentication
type AuthService struct {
	// In a real application, this would have a repository layer
	// for user data storage
}

// NewAuthService creates a new authentication service
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Authenticate validates the username and password
func (s *AuthService) Authenticate(username, password string) (bool, error) {
	// For demo purposes, hardcoded credentials
	// In a real application, this would check against a database
	return username == "demo" && password == "demo", nil
}
