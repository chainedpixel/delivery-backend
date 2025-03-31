package ports

import (
	"context"
	"time"
)

// MetricsRepository define los métodos para obtener métricas empresariales y de sucursales
type MetricsRepository interface {
	// Métricas de Empresa
	GetOrderCountByCompany(ctx context.Context, companyID string, startDate, endDate time.Time) (total, completed, cancelled int64, err error)
	GetAverageDeliveryTimeByCompany(ctx context.Context, companyID string, startDate, endDate time.Time) (float64, error)
	GetTotalRevenueByCompany(ctx context.Context, companyID string, startDate, endDate time.Time) (float64, error)
	GetActiveBranchesCountByCompany(ctx context.Context, companyID string) (int, error)
	GetUniqueCustomersByCompany(ctx context.Context, companyID string, startDate, endDate time.Time) (int, error)

	// Métricas de Sucursal
	GetOrderCountByBranch(ctx context.Context, branchID string, startDate, endDate time.Time) (total, completed, cancelled int64, err error)
	GetAverageDeliveryTimeByBranch(ctx context.Context, branchID string, startDate, endDate time.Time) (float64, error)
	GetTotalRevenueByBranch(ctx context.Context, branchID string, startDate, endDate time.Time) (float64, error)
	GetActiveDriversCountByBranch(ctx context.Context, branchID string) (int, error)
	GetUniqueCustomersByBranch(ctx context.Context, branchID string, startDate, endDate time.Time) (int, error)
	GetPeakHourOrderRateByBranch(ctx context.Context, branchID string, date time.Time) (float64, error)
}
