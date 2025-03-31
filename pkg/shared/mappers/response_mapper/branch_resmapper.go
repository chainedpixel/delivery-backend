package response_mapper

import (
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/dto"
)

// BranchToResponseDTO mapea una entidad de sucursal a su DTO de respuesta
func BranchToResponseDTO(branch *entities.Branch, includeDetails bool) dto.BranchResponse {
	response := dto.BranchResponse{
		ID:             branch.ID,
		CompanyID:      branch.CompanyID,
		Name:           branch.Name,
		Code:           branch.Code,
		ContactName:    branch.ContactName,
		ContactPhone:   branch.ContactPhone,
		ContactEmail:   branch.ContactEmail,
		IsActive:       branch.IsActive,
		ZoneID:         branch.ZoneID,
		OperatingHours: branch.OperatingHours,
		CreatedAt:      branch.CreatedAt,
		UpdatedAt:      branch.UpdatedAt,
	}

	// Incluir nombre de la compañía si está disponible
	if branch.Company != nil {
		response.CompanyName = branch.Company.Name
	}

	// Incluir nombre de la zona si está disponible
	if branch.Zone != nil {
		response.ZoneName = branch.Zone.Name
	}

	return response
}

// MapBranchesToResponse mapea un conjunto de sucursales a una respuesta paginada
func MapBranchesToResponse(branches []entities.Branch, params *entities.BranchQueryParams, total int64) *dto.PaginatedResponse {
	responseItems := make([]dto.BranchListResponse, len(branches))

	for i, branch := range branches {
		companyName := ""
		if branch.Company != nil {
			companyName = branch.Company.Name
		}

		zoneName := ""
		if branch.Zone != nil {
			zoneName = branch.Zone.Name
		}

		// Aquí podrías obtener el número de órdenes activas para cada sucursal,
		// pero por simplicidad usamos un valor estático como ejemplo
		activeOrders := 0

		responseItems[i] = dto.BranchListResponse{
			ID:           branch.ID,
			CompanyID:    branch.CompanyID,
			CompanyName:  companyName,
			Name:         branch.Name,
			Code:         branch.Code,
			ContactName:  branch.ContactName,
			ContactEmail: branch.ContactEmail,
			IsActive:     branch.IsActive,
			ZoneName:     zoneName,
			ActiveOrders: activeOrders,
			CreatedAt:    branch.CreatedAt,
		}
	}

	return &dto.PaginatedResponse{
		Data:       responseItems,
		TotalItems: total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: calculateTotalPages(total, params.PageSize),
	}
}

// BranchToMetricsDTO mapea las métricas de una sucursal a su DTO de respuesta
func BranchToMetricsDTO(metrics *entities.BranchMetrics) *dto.BranchMetricsResponse {
	if metrics == nil {
		return nil
	}

	return &dto.BranchMetricsResponse{
		TotalOrders:         metrics.TotalOrders,
		CompletedOrders:     metrics.CompletedOrders,
		CancelledOrders:     metrics.CancelledOrders,
		DeliverySuccessRate: metrics.DeliverySuccessRate,
		AverageDeliveryTime: metrics.AverageDeliveryTime,
		TotalRevenue:        metrics.TotalRevenue,
		ActiveDrivers:       metrics.ActiveDrivers,
		UniqueCustomers:     metrics.UniqueCustomers,
		PeakHourOrderRate:   metrics.PeakHourOrderRate,
	}
}
