package response_mapper

import (
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/dto"
)

// CompanyToResponseDTO mapea una entidad de compañía a su DTO de respuesta
func CompanyToResponseDTO(company *entities.Company, includeDetails bool) *dto.CompanyResponse {
	response := &dto.CompanyResponse{
		ID:                company.ID,
		Name:              company.Name,
		LegalName:         company.LegalName,
		TaxID:             company.TaxID,
		ContactEmail:      company.ContactEmail,
		ContactPhone:      company.ContactPhone,
		Website:           company.Website,
		IsActive:          company.IsActive,
		ContractDetails:   company.ContractDetails,
		DeliveryRate:      company.DeliveryRate,
		LogoURL:           company.LogoURL,
		ContractStartDate: company.ContractStartDate,
		ContractEndDate:   company.ContractEndDate,
		CreatedAt:         company.CreatedAt,
		UpdatedAt:         company.UpdatedAt,
	}

	// Incluir direcciones si existen y si se solicitan detalles
	if includeDetails && company.Address != nil {
		response.Addresses = []dto.CompanyAddressResponse{
			CompanyAddressToResponseDTO(*company.Address),
		}
	}

	// Incluir sucursales si existen y si se solicitan detalles
	if includeDetails && company.Branches != nil {
		response.Branches = make([]dto.BranchResponse, len(company.Branches))
		for i, branch := range company.Branches {
			response.Branches[i] = BranchToResponseDTO(&branch, false)
		}
	}

	return response
}

// CompanyAddressToResponseDTO mapea una entidad de dirección de compañía a su DTO de respuesta
func CompanyAddressToResponseDTO(address entities.CompanyAddress) dto.CompanyAddressResponse {
	return dto.CompanyAddressResponse{
		ID:           address.ID,
		CompanyID:    address.CompanyID,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		State:        address.State,
		PostalCode:   address.PostalCode,
		Latitude:     address.Latitude,
		Longitude:    address.Longitude,
		IsMain:       address.IsMain,
		CreatedAt:    address.CreatedAt,
	}
}

// MapCompaniesToResponse mapea un conjunto de compañías a una respuesta paginada
func MapCompaniesToResponse(companies []entities.Company, params *entities.CompanyQueryParams, total int64) *dto.PaginatedResponse {
	responseItems := make([]dto.CompanyListResponse, len(companies))

	for i, company := range companies {
		branchCount := 0
		if company.Branches != nil {
			branchCount = len(company.Branches)
		}

		responseItems[i] = dto.CompanyListResponse{
			ID:                company.ID,
			Name:              company.Name,
			LegalName:         company.LegalName,
			TaxID:             company.TaxID,
			ContactEmail:      company.ContactEmail,
			ContactPhone:      company.ContactPhone,
			IsActive:          company.IsActive,
			BranchCount:       branchCount,
			ContractStartDate: company.ContractStartDate,
			ContractEndDate:   company.ContractEndDate,
			CreatedAt:         company.CreatedAt,
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

// CompanyToMetricsDTO mapea las métricas de una compañía a su DTO de respuesta
func CompanyToMetricsDTO(metrics *entities.CompanyMetrics) *dto.CompanyMetricsResponse {
	if metrics == nil {
		return nil
	}

	return &dto.CompanyMetricsResponse{
		TotalOrders:         metrics.TotalOrders,
		CompletedOrders:     metrics.CompletedOrders,
		CancelledOrders:     metrics.CancelledOrders,
		DeliverySuccessRate: metrics.DeliverySuccessRate,
		AverageDeliveryTime: metrics.AverageDeliveryTime,
		TotalRevenue:        metrics.TotalRevenue,
		ActiveBranches:      metrics.ActiveBranches,
		UniqueCustomers:     metrics.UniqueCustomers,
	}
}
