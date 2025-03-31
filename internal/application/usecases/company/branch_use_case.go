package company

import (
	"context"
	"github.com/MarlonG1/delivery-backend/internal/application/ports"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/interfaces"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/auth"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	errPackage "github.com/MarlonG1/delivery-backend/internal/domain/error"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
	"net/http"
	"strconv"
)

type BranchUseCase struct {
	companyService interfaces.Companyrer
}

func NewBranchUseCase(companyService interfaces.Companyrer) ports.BranchUseCase {
	return &BranchUseCase{
		companyService: companyService,
	}
}

// GetBranches obtiene todas las sucursales con filtros y paginación
func (uc *BranchUseCase) GetBranches(ctx context.Context, request *http.Request) ([]entities.Branch, *entities.BranchQueryParams, int64, error) {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return nil, nil, 0, errPackage.NewDomainErrorWithCause("BranchUseCase", "GetBranches", "Failed to get claims from context", nil)
	}

	// 2. Parsear parámetros de consulta
	params := uc.parseBranchQueryParams(request)

	// 3. Asignar el companyID del usuario autenticado
	params.CompanyID = claims.CompanyID

	// 4. Obtener las sucursales según los parámetros
	branches, err := uc.companyService.GetCompanyBranches(ctx, claims.CompanyID)
	if err != nil {
		logs.Error("Failed to get branches", map[string]interface{}{
			"error":      err.Error(),
			"company_id": claims.CompanyID,
		})
		return nil, nil, 0, err
	}

	// 5. Filtrar y aplicar paginación a las sucursales (en un caso real esto debería hacerse en la BD)
	filteredBranches, total := uc.filterAndPaginateBranches(branches, params)

	return filteredBranches, params, total, nil
}

// GetBranchByID obtiene una sucursal por su ID
func (uc *BranchUseCase) GetBranchByID(ctx context.Context, branchID string) (*entities.Branch, error) {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return nil, errPackage.NewDomainErrorWithCause("BranchUseCase", "GetBranchByID", "Failed to get claims from context", nil)
	}

	// 2. Obtener la sucursal
	branch, err := uc.companyService.GetBranchByID(ctx, branchID)
	if err != nil {
		logs.Error("Failed to get branch", map[string]interface{}{
			"error":     err.Error(),
			"branch_id": branchID,
		})
		return nil, err
	}

	// 3. Verificar que la sucursal pertenezca a la empresa del usuario
	if branch.CompanyID != claims.CompanyID {
		logs.Error("Branch does not belong to user's company", map[string]interface{}{
			"branch_id":  branchID,
			"company_id": claims.CompanyID,
		})
		return nil, errPackage.NewDomainError("BranchUseCase", "GetBranchByID", "Branch does not belong to user's company")
	}

	return branch, nil
}

// GetBranchesByCompany obtiene todas las sucursales de una empresa
func (uc *BranchUseCase) GetBranchesByCompany(ctx context.Context, companyID string) ([]entities.Branch, error) {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return nil, errPackage.NewDomainErrorWithCause("BranchUseCase", "GetBranchesByCompany", "Failed to get claims from context", nil)
	}

	// 2. Si se especifica un companyID diferente, verificar permisos (sólo admins pueden ver otras empresas)
	if companyID != "" && companyID != claims.CompanyID && claims.Role != "ADMIN" {
		logs.Error("User cannot access branches from another company", map[string]interface{}{
			"user_company_id":      claims.CompanyID,
			"requested_company_id": companyID,
		})
		return nil, errPackage.NewDomainError("BranchUseCase", "GetBranchesByCompany", "User cannot access branches from another company")
	}

	// 3. Si no se especifica companyID, usar el del usuario autenticado
	if companyID == "" {
		companyID = claims.CompanyID
	}

	// 4. Obtener las sucursales
	branches, err := uc.companyService.GetCompanyBranches(ctx, companyID)
	if err != nil {
		logs.Error("Failed to get branches", map[string]interface{}{
			"error":      err.Error(),
			"company_id": companyID,
		})
		return nil, err
	}

	return branches, nil
}

// CreateBranch crea una nueva sucursal para una empresa
func (uc *BranchUseCase) CreateBranch(ctx context.Context, companyID string, branch *entities.Branch) error {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return errPackage.NewDomainErrorWithCause("BranchUseCase", "CreateBranch", "Failed to get claims from context", nil)
	}

	// 2. Si se especifica un companyID diferente, verificar permisos (sólo admins pueden crear para otras empresas)
	if companyID != "" && companyID != claims.CompanyID && claims.Role != "ADMIN" {
		logs.Error("User cannot create branch for another company", map[string]interface{}{
			"user_company_id":      claims.CompanyID,
			"requested_company_id": companyID,
		})
		return errPackage.NewDomainError("BranchUseCase", "CreateBranch", "User cannot create branch for another company")
	}

	// 3. Si no se especifica companyID, usar el del usuario autenticado
	if companyID == "" {
		companyID = claims.CompanyID
	}

	// 4. Crear la sucursal
	err := uc.companyService.AddBranchToCompany(ctx, companyID, branch)
	if err != nil {
		logs.Error("Failed to create branch", map[string]interface{}{
			"error":      err.Error(),
			"company_id": companyID,
		})
		return err
	}

	return nil
}

// UpdateBranch actualiza una sucursal existente
func (uc *BranchUseCase) UpdateBranch(ctx context.Context, branchID string, branch *entities.Branch) error {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return errPackage.NewDomainErrorWithCause("BranchUseCase", "UpdateBranch", "Failed to get claims from context", nil)
	}

	// 2. Obtener la sucursal actual para verificar propiedad
	existingBranch, err := uc.companyService.GetBranchByID(ctx, branchID)
	if err != nil {
		logs.Error("Failed to get branch", map[string]interface{}{
			"error":     err.Error(),
			"branch_id": branchID,
		})
		return err
	}

	// 3. Verificar que la sucursal pertenezca a la empresa del usuario
	if existingBranch.CompanyID != claims.CompanyID && claims.Role != "ADMIN" {
		logs.Error("Branch does not belong to user's company", map[string]interface{}{
			"branch_id":  branchID,
			"company_id": claims.CompanyID,
		})
		return errPackage.NewDomainError("BranchUseCase", "UpdateBranch", "Branch does not belong to user's company")
	}

	// 4. Asignar el ID a la sucursal a actualizar
	branch.ID = branchID

	// 5. Mantener la empresa original
	branch.CompanyID = existingBranch.CompanyID

	// 6. Actualizar la sucursal
	err = uc.companyService.UpdateBranch(ctx, branch)
	if err != nil {
		logs.Error("Failed to update branch", map[string]interface{}{
			"error":     err.Error(),
			"branch_id": branchID,
		})
		return err
	}

	return nil
}

// DeactivateBranch desactiva una sucursal
func (uc *BranchUseCase) DeactivateBranch(ctx context.Context, branchID string) error {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return errPackage.NewDomainErrorWithCause("BranchUseCase", "DeactivateBranch", "Failed to get claims from context", nil)
	}

	// 2. Obtener la sucursal para verificar propiedad
	branch, err := uc.companyService.GetBranchByID(ctx, branchID)
	if err != nil {
		logs.Error("Failed to get branch", map[string]interface{}{
			"error":     err.Error(),
			"branch_id": branchID,
		})
		return err
	}

	// 3. Verificar que la sucursal pertenezca a la empresa del usuario
	if branch.CompanyID != claims.CompanyID && claims.Role != "ADMIN" {
		logs.Error("Branch does not belong to user's company", map[string]interface{}{
			"branch_id":  branchID,
			"company_id": claims.CompanyID,
		})
		return errPackage.NewDomainError("BranchUseCase", "DeactivateBranch", "Branch does not belong to user's company")
	}

	// 4. Desactivar la sucursal
	err = uc.companyService.DeactivateBranch(ctx, branchID)
	if err != nil {
		logs.Error("Failed to deactivate branch", map[string]interface{}{
			"error":     err.Error(),
			"branch_id": branchID,
		})
		return err
	}

	return nil
}

func (uc *BranchUseCase) ReactivateBranch(ctx context.Context, branchID string) error {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return errPackage.NewDomainErrorWithCause("BranchUseCase", "AssignZoneToBranch", "Failed to get claims from context", nil)
	}

	// 2. Obtener la sucursal para verificar propiedad
	branch, err := uc.companyService.GetBranchByID(ctx, branchID)
	if err != nil {
		logs.Error("Failed to get branch", map[string]interface{}{
			"error":     err.Error(),
			"branch_id": branchID,
		})
		return err
	}

	// 3. Verificar que la sucursal pertenezca a la empresa del usuario
	if branch.CompanyID != claims.CompanyID && claims.Role != "ADMIN" {
		logs.Error("Branch does not belong to user's company", map[string]interface{}{
			"branch_id":  branchID,
			"company_id": claims.CompanyID,
		})
		return errPackage.NewDomainError("BranchUseCase", "AssignZoneToBranch", "Branch does not belong to user's company")
	}

	err = uc.companyService.ReactivateBranch(ctx, branchID)
	if err != nil {
		logs.Error("Failed to reactivate branch", map[string]interface{}{
			"error":     err.Error(),
			"branch_id": branchID,
		})
		return err
	}

	return nil
}

// AssignZoneToBranch asigna una zona a una sucursal
func (uc *BranchUseCase) AssignZoneToBranch(ctx context.Context, branchID string, zoneID string) error {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return errPackage.NewDomainErrorWithCause("BranchUseCase", "AssignZoneToBranch", "Failed to get claims from context", nil)
	}

	// 2. Obtener la sucursal para verificar propiedad
	branch, err := uc.companyService.GetBranchByID(ctx, branchID)
	if err != nil {
		logs.Error("Failed to get branch", map[string]interface{}{
			"error":     err.Error(),
			"branch_id": branchID,
		})
		return err
	}

	// 3. Verificar que la sucursal pertenezca a la empresa del usuario
	if branch.CompanyID != claims.CompanyID && claims.Role != "ADMIN" {
		logs.Error("Branch does not belong to user's company", map[string]interface{}{
			"branch_id":  branchID,
			"company_id": claims.CompanyID,
		})
		return errPackage.NewDomainError("BranchUseCase", "AssignZoneToBranch", "Branch does not belong to user's company")
	}

	// 4. Asignar la zona a la sucursal
	err = uc.companyService.AssignZoneToBranch(ctx, branchID, zoneID)
	if err != nil {
		logs.Error("Failed to assign zone to branch", map[string]interface{}{
			"error":     err.Error(),
			"branch_id": branchID,
			"zone_id":   zoneID,
		})
		return err
	}

	return nil
}

// GetAvailableZonesForBranch obtiene las zonas disponibles para una sucursal
func (uc *BranchUseCase) GetAvailableZonesForBranch(ctx context.Context, branchID string) ([]entities.Zone, error) {
	// 1. Si se proporciona un branchID, verificar que pertenezca a la empresa del usuario
	if branchID != "" {
		// Obtener los claims del contexto
		claims, ok := ctx.Value("claims").(*auth.AuthClaims)
		if !ok {
			logs.Error("Failed to get claims from context", map[string]interface{}{
				"error": "Failed to get claims from context",
			})
			return nil, errPackage.NewDomainErrorWithCause("BranchUseCase", "GetAvailableZonesForBranch", "Failed to get claims from context", nil)
		}

		// Verificar propiedad de la sucursal
		branch, err := uc.companyService.GetBranchByID(ctx, branchID)
		if err != nil {
			logs.Error("Failed to get branch", map[string]interface{}{
				"error":     err.Error(),
				"branch_id": branchID,
			})
			return nil, err
		}

		if branch.CompanyID != claims.CompanyID && claims.Role != "ADMIN" {
			logs.Error("Branch does not belong to user's company", map[string]interface{}{
				"branch_id":  branchID,
				"company_id": claims.CompanyID,
			})
			return nil, errPackage.NewDomainError("BranchUseCase", "GetAvailableZonesForBranch", "Branch does not belong to user's company")
		}
	}

	// 2. Obtener las zonas disponibles
	zones, err := uc.companyService.GetAvailableZonesForBranch(ctx, branchID)
	if err != nil {
		logs.Error("Failed to get available zones", map[string]interface{}{
			"error":     err.Error(),
			"branch_id": branchID,
		})
		return nil, err
	}

	return zones, nil
}

// GetBranchMetrics obtiene las métricas de una sucursal
func (uc *BranchUseCase) GetBranchMetrics(ctx context.Context, branchID string) (*entities.BranchMetrics, error) {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		return nil, errPackage.NewDomainErrorWithCause("BranchUseCase", "GetBranchMetrics", "Failed to get claims from context", nil)
	}

	// 2. Verificar propiedad de la sucursal
	branch, err := uc.companyService.GetBranchByID(ctx, branchID)
	if err != nil {
		logs.Error("Failed to get branch", map[string]interface{}{
			"error":     err.Error(),
			"branch_id": branchID,
		})
		return nil, err
	}

	if branch.CompanyID != claims.CompanyID {
		logs.Error("Branch does not belong to user's company", map[string]interface{}{
			"branch_id":  branchID,
			"company_id": claims.CompanyID,
		})
		return nil, errPackage.NewDomainError("BranchUseCase", "GetBranchMetrics", "Branch does not belong to user's company")
	}

	// 3. Obtener las métricas
	metrics, err := uc.companyService.GetBranchMetrics(ctx, branchID)
	if err != nil {
		logs.Error("Failed to get branch metrics", map[string]interface{}{
			"error":     err.Error(),
			"branch_id": branchID,
		})
		return nil, err
	}

	return metrics, nil
}

// Función auxiliar para parsear los parámetros de consulta
func (uc *BranchUseCase) parseBranchQueryParams(r *http.Request) *entities.BranchQueryParams {
	params := &entities.BranchQueryParams{}

	// Filtros específicos para sucursales
	params.Name = r.URL.Query().Get("name")
	params.Code = r.URL.Query().Get("code")
	params.ContactName = r.URL.Query().Get("contact_name")
	params.ContactEmail = r.URL.Query().Get("contact_email")
	params.ZoneID = r.URL.Query().Get("zone_id")

	// Estado activo/inactivo
	isActive := r.URL.Query().Get("is_active")
	if isActive != "" {
		active := isActive == "true" || isActive == "1"
		params.IsActive = &active
	}

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

	return params
}

// Función auxiliar para filtrar y paginar las sucursales (simulando lo que normalmente se haría en BD)
func (uc *BranchUseCase) filterAndPaginateBranches(branches []entities.Branch, params *entities.BranchQueryParams) ([]entities.Branch, int64) {
	// Filtrar las sucursales
	var filteredBranches []entities.Branch
	for _, branch := range branches {
		// Aplicar filtros
		if params.Name != "" && branch.Name != params.Name {
			continue
		}
		if params.Code != "" && branch.Code != params.Code {
			continue
		}
		if params.ContactName != "" && branch.ContactName != params.ContactName {
			continue
		}
		if params.ContactEmail != "" && branch.ContactEmail != params.ContactEmail {
			continue
		}
		if params.ZoneID != "" && branch.ZoneID != params.ZoneID {
			continue
		}
		if params.IsActive != nil && branch.IsActive != *params.IsActive {
			continue
		}

		filteredBranches = append(filteredBranches, branch)
	}

	// Calcular total
	total := int64(len(filteredBranches))

	// Aplicar paginación
	start := (params.Page - 1) * params.PageSize
	end := start + params.PageSize
	if start >= len(filteredBranches) {
		return []entities.Branch{}, total
	}
	if end > len(filteredBranches) {
		end = len(filteredBranches)
	}

	return filteredBranches[start:end], total
}
