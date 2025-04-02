package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Service istifadəçi autentifikasiyası biznes məntiqini müəyyən edir
type Service interface {
	Login(ctx context.Context, username, password string) (*User, error)
}

// AuthService Service interfeysini həyata keçirir
type AuthService struct {
	repo Repository
}

// NewAuthService yeni AuthService yaradır
func NewAuthService(repo Repository) *AuthService {
	return &AuthService{repo: repo}
}

// Login istifadəçi adı və şifrəyə görə istifadəçini yoxlayır
func (s *AuthService) Login(ctx context.Context, username, password string) (*User, error) {
	if username == "" || password == "" {
		return nil, errors.New("istifadəçi adı və şifrə tələb olunur")
	}

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("istifadəçi adı və ya şifrə yanlışdır")
	}

	// Şifrəni yoxla
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("istifadəçi adı və ya şifrə yanlışdır")
	}

	// Şifrəni silmək (təhlükəsizlik üçün)
	user.Password = ""

	return user, nil
}
