package company

import (
	"context"
	"github.com/MarlonG1/delivery-backend/internal/application/ports"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/constants"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/interfaces"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/auth"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	errPackage "github.com/MarlonG1/delivery-backend/internal/domain/error"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
	"net/http"
	"strconv"
	"time"
)

type CompanyUseCase struct {
	companyService interfaces.Companyrer
}

func NewCompanyUseCase(companyService interfaces.Companyrer) ports.CompanyUseCase {
	return &CompanyUseCase{
		companyService: companyService,
	}
}

// GetCompanyByID obtiene una empresa por su ID
func (uc *CompanyUseCase) GetCompanyByID(ctx context.Context) (*entities.Company, error) {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return nil, errPackage.NewDomainErrorWithCause("CompanyUseCase", "GetCompanyByID", "Failed to get claims from context", nil)
	}

	company, err := uc.companyService.GetCompanyByID(ctx, claims.CompanyID)
	if err != nil {
		logs.Error("Failed to get company by ID", map[string]interface{}{
			"error":      err.Error(),
			"company_id": claims.CompanyID,
		})
		return nil, err
	}

	return company, nil
}

// CreateCompany crea una nueva empresa
func (uc *CompanyUseCase) CreateCompany(ctx context.Context, company *entities.Company) error {
	// Validamos que la empresa sea válida según reglas de negocio
	if err := uc.companyService.ValidateCompany(ctx, company); err != nil {
		logs.Error("Failed to validate company", map[string]interface{}{
			"error":        err.Error(),
			"company_name": company.Name,
		})
		return err
	}

	// Creamos la empresa usando el servicio de dominio
	err := uc.companyService.CreateCompany(ctx, company)
	if err != nil {
		logs.Error("Failed to create company", map[string]interface{}{
			"error":        err.Error(),
			"company_name": company.Name,
		})
		return err
	}

	return nil
}

// UpdateCompany actualiza una empresa existente
func (uc *CompanyUseCase) UpdateCompany(ctx context.Context, company *entities.Company) error {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return errPackage.NewDomainErrorWithCause("CompanyUseCase", "UpdateCompany", "Failed to get claims from context", nil)
	}
	company.ID = claims.CompanyID

	// Actualizar la empresa
	err := uc.companyService.UpdateCompany(ctx, company)
	if err != nil {
		logs.Error("Failed to update company", map[string]interface{}{
			"error":      err.Error(),
			"company_id": claims.CompanyID,
		})
		return err
	}

	return nil
}

// DeactivateCompany desactiva una empresa
func (uc *CompanyUseCase) DeactivateCompany(ctx context.Context, companyID string) error {
	err := uc.companyService.DeactivateCompany(ctx, companyID)
	if err != nil {
		logs.Error("Failed to deactivate company", map[string]interface{}{
			"error":      err.Error(),
			"company_id": companyID,
		})
		return err
	}

	return nil
}

// ReactivateCompany reactiva una empresa
func (uc *CompanyUseCase) ReactivateCompany(ctx context.Context, companyID string) error {
	err := uc.companyService.ReactivateCompany(ctx, companyID)
	if err != nil {
		logs.Error("Failed to reactivate company", map[string]interface{}{
			"error":      err.Error(),
			"company_id": companyID,
		})
		return err
	}

	return nil
}

// GetCompanyAddresses obtiene las direcciones de una empresa
func (uc *CompanyUseCase) GetCompanyAddresses(ctx context.Context) ([]entities.CompanyAddress, error) {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return nil, errPackage.NewDomainErrorWithCause("CompanyUseCase", "GetCompanyAddresses", "Failed to get claims from context", nil)
	}

	// Obtener las direcciones de la empresa
	addresses, err := uc.companyService.GetAddresses(ctx, claims.CompanyID)
	if err != nil {
		logs.Error("Failed to get company addresses", map[string]interface{}{
			"error":      err.Error(),
			"company_id": claims.CompanyID,
		})
		return nil, errPackage.NewDomainErrorWithCause("CompanyUseCase", "GetCompanyAddresses", "Error getting company addresses", err)
	}

	return addresses, nil
}

// AddCompanyAddress añade una dirección a una empresa
func (uc *CompanyUseCase) AddCompanyAddress(ctx context.Context, address *entities.CompanyAddress) error {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return errPackage.NewDomainErrorWithCause("CompanyUseCase", "AddCompanyAddress", "Failed to get claims from context", nil)
	}

	// Asignar el ID de la empresa a la dirección
	address.CompanyID = claims.CompanyID

	// Añadir la dirección
	err := uc.companyService.AddCompanyAddress(ctx, claims.CompanyID, address)
	if err != nil {
		logs.Error("Failed to add company address", map[string]interface{}{
			"error":      err.Error(),
			"company_id": claims.CompanyID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyUseCase", "AddCompanyAddress", "Error adding company address", err)
	}

	return nil
}

// UpdateCompanyAddress actualiza una dirección existente
func (uc *CompanyUseCase) UpdateCompanyAddress(ctx context.Context, addressID string, address *entities.CompanyAddress) error {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return errPackage.NewDomainErrorWithCause("CompanyUseCase", "UpdateCompanyAddress", "Failed to get claims from context", nil)
	}

	// Verificar que la dirección existe y sea de la compañia del usuario
	existingAddress, err := uc.companyService.GetAddressByID(ctx, addressID, claims.CompanyID)
	if err != nil {
		return err
	}

	// Mantener el ID y la empresa de la dirección original
	address.ID = addressID
	address.CompanyID = existingAddress.CompanyID

	// Actualizar la dirección
	err = uc.companyService.UpdateCompanyAddress(ctx, claims.CompanyID, address)
	if err != nil {
		logs.Error("Failed to update company address", map[string]interface{}{
			"error":      err.Error(),
			"address_id": addressID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyUseCase", "UpdateCompanyAddress", "Error updating company address", err)
	}

	return nil
}

// DeleteCompanyAddress elimina una dirección
func (uc *CompanyUseCase) DeleteCompanyAddress(ctx context.Context, addressID string) error {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return errPackage.NewDomainErrorWithCause("CompanyUseCase", "DeleteCompanyAddress", "Failed to get claims from context", nil)
	}

	// Eliminar la dirección
	err := uc.companyService.DeleteCompanyAddress(ctx, addressID, claims.CompanyID)
	if err != nil {
		logs.Error("Failed to delete company address", map[string]interface{}{
			"error":      err.Error(),
			"address_id": addressID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyUseCase", "DeleteCompanyAddress", "Error deleting company address", err)
	}

	return nil
}

// GetCompanyMetrics obtiene las métricas de una empresa
func (uc *CompanyUseCase) GetCompanyMetrics(ctx context.Context) (*entities.CompanyMetrics, error) {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return nil, errPackage.NewDomainErrorWithCause("CompanyUseCase", "GetCompanyMetrics", "Failed to get claims from context", nil)
	}

	metrics, err := uc.companyService.GetCompanyMetrics(ctx, claims.CompanyID)
	if err != nil {
		logs.Error("Failed to get company metrics", map[string]interface{}{
			"error":      err.Error(),
			"company_id": claims.CompanyID,
		})
		return nil, err
	}

	return metrics, nil
}

// GetCompanies obtiene todas las empresas con filtros y paginación
func (uc *CompanyUseCase) GetCompanies(ctx context.Context, request *http.Request) ([]entities.Company, *entities.CompanyQueryParams, int64, error) {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return nil, nil, 0, errPackage.NewDomainErrorWithCause("CompanyUseCase", "GetCompanies", "Failed to get claims from context", nil)
	}

	// 2. Verificar permisos de acceso (opcional, según requerimientos)
	if claims.Role != constants.AdminRole {
		logs.Error("User does not have admin permissions", map[string]interface{}{
			"user_id": claims.UserID,
			"role":    claims.Role,
		})
		return nil, nil, 0, errPackage.NewDomainError("CompanyUseCase", "GetCompanies", "User does not have sufficient permissions")
	}

	// 3. Parsear parámetros de consulta
	params := uc.parseCompanyQueryParams(request)

	// 4. Obtener las empresas según los parámetros
	companies, total, err := uc.companyService.GetCompanies(ctx, params)
	if err != nil {
		logs.Error("Failed to get companies", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, nil, 0, err
	}

	return companies, params, total, nil
}

// Función auxiliar para parsear los parámetros de consulta
func (uc *CompanyUseCase) parseCompanyQueryParams(r *http.Request) *entities.CompanyQueryParams {
	params := &entities.CompanyQueryParams{}

	// Filtros específicos para compañías
	params.Name = r.URL.Query().Get("name")
	params.ContactEmail = r.URL.Query().Get("contact_email")
	params.ContactPhone = r.URL.Query().Get("contact_phone")

	// Estado activo/inactivo
	isActive := r.URL.Query().Get("is_active")
	if isActive != "" {
		active := isActive == "true" || isActive == "1"
		params.IsActive = &active
	}

	// Fechas de contrato
	contractStartDateStr := r.URL.Query().Get("contract_start_date")
	if contractStartDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, contractStartDateStr); err == nil {
			params.ContractStartDate = &startDate
		}
	}

	contractEndDateStr := r.URL.Query().Get("contract_end_date")
	if contractEndDateStr != "" {
		if endDate, err := time.Parse(time.RFC3339, contractEndDateStr); err == nil {
			params.ContractEndDate = &endDate
		}
	}

	// Incluir eliminados
	includeDeletedStr := r.URL.Query().Get("include_deleted")
	params.IncludeDeleted = includeDeletedStr == "true" || includeDeletedStr == "1"

	// Parsear parámetros de paginación
	uc.parsePaginationQueryParams(r, &params.PaginationQueryParams)

	return params
}

// Función auxiliar para parsear los parámetros de paginación
func (uc *CompanyUseCase) parsePaginationQueryParams(r *http.Request, params *entities.PaginationQueryParams) {
	// Paginación
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		} else {
			params.Page = 1 // Default
		}
	} else {
		params.Page = 1 // Default
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			params.PageSize = pageSize
		} else {
			params.PageSize = 10 // Default
		}
	} else {
		params.PageSize = 10 // Default
	}

	// Ordenamiento
	params.SortBy = r.URL.Query().Get("sort_by")
	params.SortDirection = r.URL.Query().Get("sort_direction")
	if params.SortDirection != "asc" && params.SortDirection != "desc" {
		params.SortDirection = "desc" // Default
	}
}
