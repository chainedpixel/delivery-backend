package interfaces

import (
	"context"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
)

// MetricsService define los métodos para obtener métricas
type MetricsService interface {
	// Métricas de empresa
	GetCompanyMetrics(ctx context.Context, companyID string) (*entities.CompanyMetrics, error)

	// Métricas de sucursal
	GetBranchMetrics(ctx context.Context, branchID, companyID string) (*entities.BranchMetrics, error)
}
