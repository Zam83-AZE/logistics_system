package auth

import (
	"time"
)

// User verilənlər bazasından gələn istifadəçi məlumatlarını təmsil edir
type User struct {
	ID        int       `db:"id" json:"id"`
	Username  string    `db:"username" json:"username" validate:"required"`
	Password  string    `db:"password" json:"-" validate:"required"`
	Email     string    `db:"email" json:"email" validate:"required,email"`
	FullName  string    `db:"full_name" json:"fullName" validate:"required"`
	IsActive  bool      `db:"is_active" json:"isActive"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// LoginForm istifadəçi giriş formunu təmsil edir
type LoginForm struct {
	Username string
	Password string
	Error    string
}
