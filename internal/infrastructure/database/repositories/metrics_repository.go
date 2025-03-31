package repositories

import (
	"context"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/constants"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/ports"
	"gorm.io/gorm"
	"time"
)

type MetricsRepository struct {
	db *gorm.DB
}

func NewMetricsRepository(db *gorm.DB) ports.MetricsRepository {
	return &MetricsRepository{
		db: db,
	}
}

// Implementación de métricas para Empresas

// GetOrderCountByCompany obtiene el conteo de órdenes de una empresa por estado
func (r *MetricsRepository) GetOrderCountByCompany(ctx context.Context, companyID string, startDate, endDate time.Time) (total, completed, cancelled int64, err error) {
	// Total de órdenes
	err = r.db.WithContext(ctx).Model(&struct{}{}).
		Table("orders").
		Where("company_id = ? AND created_at BETWEEN ? AND ?", companyID, startDate, endDate).
		Count(&total).Error
	if err != nil {
		return 0, 0, 0, err
	}

	// Órdenes completadas
	err = r.db.WithContext(ctx).Model(&struct{}{}).
		Table("orders").
		Where("company_id = ? AND status = ? AND created_at BETWEEN ? AND ?",
			companyID, constants.OrderStatusDelivered, startDate, endDate).
		Count(&completed).Error
	if err != nil {
		return 0, 0, 0, err
	}

	// Órdenes canceladas
	err = r.db.WithContext(ctx).Model(&struct{}{}).
		Table("orders").
		Where("company_id = ? AND status = ? AND created_at BETWEEN ? AND ?",
			companyID, constants.OrderStatusCancelled, startDate, endDate).
		Count(&cancelled).Error
	if err != nil {
		return 0, 0, 0, err
	}

	return total, completed, cancelled, nil
}

// GetAverageDeliveryTimeByCompany calcula el tiempo promedio de entrega para una empresa
func (r *MetricsRepository) GetAverageDeliveryTimeByCompany(ctx context.Context, companyID string, startDate, endDate time.Time) (float64, error) {
	var result struct {
		AvgDeliveryTime float64
	}

	err := r.db.WithContext(ctx).
		Table("orders o").
		Joins("INNER JOIN order_details od ON o.id = od.order_id").
		Where("o.company_id = ? AND o.status = ? AND o.created_at BETWEEN ? AND ? AND od.delivered_at IS NOT NULL",
			companyID, constants.OrderStatusDelivered, startDate, endDate).
		Select("AVG(TIMESTAMPDIFF(MINUTE, o.created_at, od.delivered_at)) as avg_delivery_time").
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return result.AvgDeliveryTime, nil
}

// GetTotalRevenueByCompany calcula los ingresos totales para una empresa
func (r *MetricsRepository) GetTotalRevenueByCompany(ctx context.Context, companyID string, startDate, endDate time.Time) (float64, error) {
	var result struct {
		TotalRevenue float64
	}

	err := r.db.WithContext(ctx).
		Table("orders o").
		Joins("INNER JOIN order_details od ON o.id = od.order_id").
		Where("o.company_id = ? AND o.created_at BETWEEN ? AND ? AND o.status != ?",
			companyID, startDate, endDate, constants.OrderStatusCancelled).
		Select("SUM(od.price) as total_revenue").
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return result.TotalRevenue, nil
}

// GetActiveBranchesCountByCompany cuenta las sucursales activas de una empresa
func (r *MetricsRepository) GetActiveBranchesCountByCompany(ctx context.Context, companyID string) (int, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&struct{}{}).
		Table("company_branches").
		Where("company_id = ? AND is_active = ?", companyID, true).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// GetUniqueCustomersByCompany cuenta clientes únicos para una empresa
func (r *MetricsRepository) GetUniqueCustomersByCompany(ctx context.Context, companyID string, startDate, endDate time.Time) (int, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&struct{}{}).
		Table("orders").
		Select("COUNT(DISTINCT client_id)").
		Where("company_id = ? AND created_at BETWEEN ? AND ?", companyID, startDate, endDate).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// Implementación de métricas para Sucursales

// GetOrderCountByBranch obtiene el conteo de órdenes de una sucursal por estado
func (r *MetricsRepository) GetOrderCountByBranch(ctx context.Context, branchID string, startDate, endDate time.Time) (total, completed, cancelled int64, err error) {
	// Total de órdenes
	err = r.db.WithContext(ctx).Model(&struct{}{}).
		Table("orders").
		Where("branch_id = ? AND created_at BETWEEN ? AND ?", branchID, startDate, endDate).
		Count(&total).Error
	if err != nil {
		return 0, 0, 0, err
	}

	// Órdenes completadas
	err = r.db.WithContext(ctx).Model(&struct{}{}).
		Table("orders").
		Where("branch_id = ? AND status = ? AND created_at BETWEEN ? AND ?",
			branchID, constants.OrderStatusDelivered, startDate, endDate).
		Count(&completed).Error
	if err != nil {
		return 0, 0, 0, err
	}

	// Órdenes canceladas
	err = r.db.WithContext(ctx).Model(&struct{}{}).
		Table("orders").
		Where("branch_id = ? AND status = ? AND created_at BETWEEN ? AND ?",
			branchID, constants.OrderStatusCancelled, startDate, endDate).
		Count(&cancelled).Error
	if err != nil {
		return 0, 0, 0, err
	}

	return total, completed, cancelled, nil
}

// GetAverageDeliveryTimeByBranch calcula el tiempo promedio de entrega para una sucursal
func (r *MetricsRepository) GetAverageDeliveryTimeByBranch(ctx context.Context, branchID string, startDate, endDate time.Time) (float64, error) {
	var result struct {
		AvgDeliveryTime float64
	}

	err := r.db.WithContext(ctx).
		Table("orders o").
		Joins("INNER JOIN order_details od ON o.id = od.order_id").
		Where("o.branch_id = ? AND o.status = ? AND o.created_at BETWEEN ? AND ? AND od.delivered_at IS NOT NULL",
			branchID, constants.OrderStatusDelivered, startDate, endDate).
		Select("AVG(TIMESTAMPDIFF(MINUTE, o.created_at, od.delivered_at)) as avg_delivery_time").
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return result.AvgDeliveryTime, nil
}

// GetTotalRevenueByBranch calcula los ingresos totales para una sucursal
func (r *MetricsRepository) GetTotalRevenueByBranch(ctx context.Context, branchID string, startDate, endDate time.Time) (float64, error) {
	var result struct {
		TotalRevenue float64
	}

	err := r.db.WithContext(ctx).
		Table("orders o").
		Joins("INNER JOIN order_details od ON o.id = od.order_id").
		Where("o.branch_id = ? AND o.created_at BETWEEN ? AND ? AND o.status != ?",
			branchID, startDate, endDate, constants.OrderStatusCancelled).
		Select("SUM(od.price) as total_revenue").
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return result.TotalRevenue, nil
}

// GetActiveDriversCountByBranch cuenta los repartidores activos asignados a una sucursal
func (r *MetricsRepository) GetActiveDriversCountByBranch(ctx context.Context, branchID string) (int, error) {
	// Esta es una consulta compleja que depende de cómo esté modelada la relación entre drivers y sucursales
	// Aquí usamos una consulta simplificada basada en las órdenes asignadas
	var count int64

	// Contar drivers únicos que tienen órdenes activas en esta sucursal
	err := r.db.WithContext(ctx).Model(&struct{}{}).
		Table("orders").
		Select("COUNT(DISTINCT driver_id)").
		Where("branch_id = ? AND driver_id IS NOT NULL AND status IN (?, ?, ?)",
			branchID, constants.OrderStatusAccepted, constants.OrderStatusPickedUp, constants.OrderStatusInTransit).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// GetUniqueCustomersByBranch cuenta clientes únicos para una sucursal
func (r *MetricsRepository) GetUniqueCustomersByBranch(ctx context.Context, branchID string, startDate, endDate time.Time) (int, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&struct{}{}).
		Table("orders").
		Select("COUNT(DISTINCT client_id)").
		Where("branch_id = ? AND created_at BETWEEN ? AND ?", branchID, startDate, endDate).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// GetPeakHourOrderRateByBranch calcula la tasa de pedidos en hora pico para una sucursal
func (r *MetricsRepository) GetPeakHourOrderRateByBranch(ctx context.Context, branchID string, date time.Time) (float64, error) {
	// Primero encontramos la hora con más pedidos
	type HourlyCount struct {
		Hour  int
		Count int64
	}

	var hourlyCounts []HourlyCount

	// Ajustar fecha para obtener inicio y fin del día
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.Add(24 * time.Hour)

	err := r.db.WithContext(ctx).
		Table("orders").
		Select("HOUR(created_at) as hour, COUNT(*) as count").
		Where("branch_id = ? AND created_at BETWEEN ? AND ?", branchID, startDate, endDate).
		Group("HOUR(created_at)").
		Order("count DESC").
		Scan(&hourlyCounts).Error

	if err != nil {
		return 0, err
	}

	// Si no hay datos, devolvemos 0
	if len(hourlyCounts) == 0 {
		return 0, nil
	}

	// La tasa de pedidos en hora pico es el número máximo de pedidos en una hora
	return float64(hourlyCounts[0].Count), nil
}
