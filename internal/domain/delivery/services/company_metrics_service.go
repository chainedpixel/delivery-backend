package services

import (
	"context"
	"errors"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/interfaces"
	"gorm.io/gorm"
	"time"

	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/ports"
	errPackage "github.com/MarlonG1/delivery-backend/internal/domain/error"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
)

// CompanyMetricsService extiende el CompanyService con funcionalidad de métricas
type CompanyMetricsService struct {
	companyRepo ports.CompanyRepository
	metricsRepo ports.MetricsRepository
}

func NewCompanyMetricsService(companyRepo ports.CompanyRepository, metricsRepo ports.MetricsRepository) interfaces.MetricsService {
	return &CompanyMetricsService{
		companyRepo: companyRepo,
		metricsRepo: metricsRepo,
	}
}

// GetCompanyMetrics obtiene las métricas reales de una empresa
func (s *CompanyMetricsService) GetCompanyMetrics(ctx context.Context, companyID string) (*entities.CompanyMetrics, error) {
	// 1. Verificar que la empresa existe
	_, err := s.companyRepo.GetCompanyByID(ctx, companyID)
	if err != nil {
		logs.Error("Failed to get company by ID", map[string]interface{}{
			"error": err,
			"id":    companyID,
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errPackage.NewDomainErrorWithCause("CompanyMetricsService", "GetCompanyMetrics", "Company not found", errPackage.ErrCompanyNotFound)
		}

		return nil, errPackage.NewDomainErrorWithCause("CompanyMetricsService", "GetCompanyMetrics", "Error getting company by ID", err)
	}

	// 2. Definir el rango de fechas para las métricas (último mes)
	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0) // Un mes atrás

	// 3. Inicializar el objeto de métricas
	metrics := &entities.CompanyMetrics{}

	// 4. Obtener métricas reales usando el repositorio de métricas
	// 4.1 Conteos de órdenes
	total, completed, cancelled, err := s.metricsRepo.GetOrderCountByCompany(ctx, companyID, startDate, endDate)
	if err != nil {
		logs.Error("Failed to get order counts", map[string]interface{}{
			"error":      err,
			"company_id": companyID,
		})
		// Continuamos con las demás métricas incluso si una falla
		// Si las fallas son consistentes, podemos manejar escenarios de degradación aquí
	} else {
		metrics.TotalOrders = total
		metrics.CompletedOrders = completed
		metrics.CancelledOrders = cancelled

		// Calcular tasa de éxito si hay órdenes totales
		if total > 0 {
			metrics.DeliverySuccessRate = float64(completed) / float64(total) * 100
		}
	}

	// 4.2 Tiempo promedio de entrega
	avgTime, err := s.metricsRepo.GetAverageDeliveryTimeByCompany(ctx, companyID, startDate, endDate)
	if err != nil {
		logs.Error("Failed to get average delivery time", map[string]interface{}{
			"error":      err,
			"company_id": companyID,
		})
	} else {
		metrics.AverageDeliveryTime = avgTime
	}

	// 4.3 Ingresos totales
	revenue, err := s.metricsRepo.GetTotalRevenueByCompany(ctx, companyID, startDate, endDate)
	if err != nil {
		logs.Error("Failed to get total revenue", map[string]interface{}{
			"error":      err,
			"company_id": companyID,
		})
	} else {
		metrics.TotalRevenue = revenue
	}

	// 4.4 Sucursales activas
	activeBranches, err := s.metricsRepo.GetActiveBranchesCountByCompany(ctx, companyID)
	if err != nil {
		logs.Error("Failed to get active branches count", map[string]interface{}{
			"error":      err,
			"company_id": companyID,
		})
	} else {
		metrics.ActiveBranches = activeBranches
	}

	// 4.5 Clientes únicos
	uniqueCustomers, err := s.metricsRepo.GetUniqueCustomersByCompany(ctx, companyID, startDate, endDate)
	if err != nil {
		logs.Error("Failed to get unique customers", map[string]interface{}{
			"error":      err,
			"company_id": companyID,
		})
	} else {
		metrics.UniqueCustomers = uniqueCustomers
	}

	return metrics, nil
}

// GetBranchMetrics obtiene las métricas reales de una sucursal
func (s *CompanyMetricsService) GetBranchMetrics(ctx context.Context, branchID, companyID string) (*entities.BranchMetrics, error) {
	// 1. Verificar que la sucursal existe
	_, err := s.companyRepo.GetBranchByID(ctx, branchID)
	if err != nil {
		logs.Error("Failed to get branch by ID", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errPackage.NewDomainErrorWithCause("CompanyMetricsService", "GetBranchMetrics", "Branch not found", errPackage.ErrBranchNotFound)
		}

		return nil, errPackage.NewDomainErrorWithCause("CompanyMetricsService", "GetBranchMetrics", "Error getting branch by ID", err)
	}

	// 2. Definir el rango de fechas para las métricas (último mes)
	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0) // Un mes atrás

	// 3. Inicializar el objeto de métricas
	metrics := &entities.BranchMetrics{}

	// 4. Obtener métricas reales usando el repositorio de métricas
	// 4.1 Conteos de órdenes
	total, completed, cancelled, err := s.metricsRepo.GetOrderCountByBranch(ctx, branchID, startDate, endDate)
	if err != nil {
		logs.Error("Failed to get order counts", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
		})
		// Continuamos con las demás métricas incluso si una falla
	} else {
		metrics.TotalOrders = total
		metrics.CompletedOrders = completed
		metrics.CancelledOrders = cancelled

		// Calcular tasa de éxito si hay órdenes totales
		if total > 0 {
			metrics.DeliverySuccessRate = float64(completed) / float64(total) * 100
		}
	}

	// 4.2 Tiempo promedio de entrega
	avgTime, err := s.metricsRepo.GetAverageDeliveryTimeByBranch(ctx, branchID, startDate, endDate)
	if err != nil {
		logs.Error("Failed to get average delivery time", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
		})
	} else {
		metrics.AverageDeliveryTime = avgTime
	}

	// 4.3 Ingresos totales
	revenue, err := s.metricsRepo.GetTotalRevenueByBranch(ctx, branchID, startDate, endDate)
	if err != nil {
		logs.Error("Failed to get total revenue", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
		})
	} else {
		metrics.TotalRevenue = revenue
	}

	// 4.4 Repartidores activos
	activeDrivers, err := s.metricsRepo.GetActiveDriversCountByBranch(ctx, branchID)
	if err != nil {
		logs.Error("Failed to get active drivers count", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
		})
	} else {
		metrics.ActiveDrivers = activeDrivers
	}

	// 4.5 Clientes únicos
	uniqueCustomers, err := s.metricsRepo.GetUniqueCustomersByBranch(ctx, branchID, startDate, endDate)
	if err != nil {
		logs.Error("Failed to get unique customers", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
		})
	} else {
		metrics.UniqueCustomers = uniqueCustomers
	}

	// 4.6 Tasa de pedidos en hora pico
	peakRate, err := s.metricsRepo.GetPeakHourOrderRateByBranch(ctx, branchID, endDate)
	if err != nil {
		logs.Error("Failed to get peak hour order rate", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
		})
	} else {
		metrics.PeakHourOrderRate = peakRate
	}

	return metrics, nil
}
