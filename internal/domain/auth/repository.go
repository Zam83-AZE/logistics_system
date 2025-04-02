package auth

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Repository istifadəçi məlumatları əməliyyatlarını müəyyən edir
type Repository interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
}

// PostgresRepository Repository interfeysini həyata keçirir
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository yeni PostgresRepository yaradır
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// GetByUsername istifadəçini username-ə görə əldə edir
func (r *PostgresRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT id, username, password, email, full_name, is_active, created_at, updated_at
		FROM users
		WHERE username = $1 AND is_active = true
	`

	user := &User{}
	err := r.db.GetContext(ctx, user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // İstifadəçi tapılmadı
		}
		return nil, err
	}

	return user, nil
}
