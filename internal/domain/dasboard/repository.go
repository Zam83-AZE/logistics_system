package dashboard

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Repository dashboard məlumatları əməliyyatlarını müəyyən edir
type Repository interface {
	GetSummary(ctx context.Context) (*Summary, error)
}

// PostgresRepository Repository interfeysini həyata keçirir
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository yeni PostgresRepository yaradır
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// GetSummary dashboard üçün əsas statistikaları əldə edir
func (r *PostgresRepository) GetSummary(ctx context.Context) (*Summary, error) {
	query := `
		SELECT
			(SELECT COUNT(*) FROM users) AS total_customers,
			0 AS total_containers,
			0 AS active_shipments,
			0 AS pending_invoices
	`

	summary := &Summary{}
	err := r.db.GetContext(ctx, summary, query)
	if err != nil {
		return nil, err
	}

	return summary, nil
}
