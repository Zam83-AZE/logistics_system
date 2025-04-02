package dashboard

import (
	"context"
)

// Service dashboard biznes məntiqini müəyyən edir
type Service interface {
	GetDashboardData(ctx context.Context, username string) (*DashboardData, error)
}

// DashboardService Service interfeysini həyata keçirir
type DashboardService struct {
	repo Repository
}

// NewDashboardService yeni DashboardService yaradır
func NewDashboardService(repo Repository) *DashboardService {
	return &DashboardService{repo: repo}
}

// GetDashboardData dashboard üçün lazım olan bütün məlumatları əldə edir
func (s *DashboardService) GetDashboardData(ctx context.Context, username string) (*DashboardData, error) {
	summary, err := s.repo.GetSummary(ctx)
	if err != nil {
		return nil, err
	}

	dashboardData := &DashboardData{
		Summary:     *summary,
		UserName:    username,
		CurrentPage: "dashboard",
	}

	return dashboardData, nil
}
