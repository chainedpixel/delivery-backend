package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/auth"
	"gorm.io/gorm"
	"time"

	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/interfaces"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/ports"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/value_objects"
	errPackage "github.com/MarlonG1/delivery-backend/internal/domain/error"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
)

type CompanyService struct {
	repo           ports.CompanyRepository
	metricsService interfaces.MetricsService
}

func NewCompanyService(repo ports.CompanyRepository, metricsService interfaces.MetricsService) interfaces.Companyrer {
	return &CompanyService{
		repo:           repo,
		metricsService: metricsService,
	}
}

// Métodos existentes
func (c *CompanyService) GetAddresses(ctx context.Context, companyID string) ([]entities.CompanyAddress, error) {
	addresses, err := c.repo.GetCompanyAddresses(ctx, companyID)
	if err != nil {
		logs.Error("Failed to get company addresses", map[string]interface{}{
			"error": err,
		})
		return nil, errPackage.NewDomainErrorWithCause("CompanyService", "GetAddresses", "Error getting company addresses", err)
	}

	return addresses, nil
}

func (c *CompanyService) GetAddressByID(ctx context.Context, id, companyID string) (*entities.CompanyAddress, error) {
	address, err := c.repo.GetCompanyAddressByID(ctx, id, companyID)
	if err != nil {
		logs.Error("Failed to get company address by ID", map[string]interface{}{
			"error": err,
			"id":    id,
		})
		return nil, errPackage.NewDomainErrorWithCause("CompanyService", "GetAddressByID", "Error getting company address by ID", err)
	}

	return address, nil
}

func (c *CompanyService) GetCompanyAndBranchForUser(ctx context.Context, userID string) (string, string, error) {
	companyID, branchID, err := c.repo.GetCompanyAndBranchForUser(ctx, userID)
	if err != nil {
		logs.Error("Failed to get company and branch for user", map[string]interface{}{
			"error":   err,
			"user_id": userID,
		})
		return "", "", errPackage.NewDomainErrorWithCause("CompanyService", "GetCompanyAndBranchForUser", "Error getting company and branch for user", err)
	}

	return companyID, branchID, nil
}

// Nuevos métodos para CRUD de Company
func (c *CompanyService) GetCompanyByID(ctx context.Context, id string) (*entities.Company, error) {
	company, err := c.repo.GetCompanyByID(ctx, id)
	if err != nil {
		logs.Error("Failed to get company by ID", map[string]interface{}{
			"error": err,
			"id":    id,
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errPackage.NewDomainErrorWithCause("CompanyService", "GetCompanyByID", "Company not found", errPackage.ErrCompanyNotFound)
		}

		return nil, errPackage.NewDomainErrorWithCause("CompanyService", "GetCompanyByID", "Error getting company by ID", err)
	}

	// Verificar si la empresa está activa
	if !company.IsActive {
		logs.Warn("Company is inactive", map[string]interface{}{
			"company_id": id,
		})
		return company, nil // Devolvemos la empresa aunque esté inactiva, pero con un warning en logs
	}

	return company, nil
}

// CreateCompany crea una nueva empresa con sus respectivas validaciones de dominio
func (c *CompanyService) CreateCompany(ctx context.Context, company *entities.Company) error {
	// 1. Validar la empresa
	if err := c.ValidateCompany(ctx, company); err != nil {
		return err
	}

	// 2. Crear la empresa
	err := c.repo.CreateCompany(ctx, company)
	if err != nil {
		logs.Error("Failed to create company", map[string]interface{}{
			"error":   err,
			"company": company.Name,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "CreateCompany", "Error creating company", err)
	}

	logs.Info("Company created successfully", map[string]interface{}{
		"company_id":   company.ID,
		"company_name": company.Name,
	})

	return nil
}

// UpdateCompany actualiza una empresa existente con validaciones de dominio
func (c *CompanyService) UpdateCompany(ctx context.Context, company *entities.Company) error {
	// 1. Verificar que la empresa existe
	existingCompany, err := c.GetCompanyByID(ctx, company.ID)
	if err != nil {
		return err
	}

	// 2. Validar la empresa actualizada
	if err := c.ValidateCompanyUpdate(ctx, company, existingCompany); err != nil {
		return err
	}

	// 3. Actualizar la empresa
	err = c.repo.UpdateCompany(ctx, company)
	if err != nil {
		logs.Error("Failed to update company", map[string]interface{}{
			"error":      err,
			"company_id": company.ID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "UpdateCompany", "Error updating company", err)
	}

	logs.Info("Company updated successfully", map[string]interface{}{
		"company_id": company.ID,
	})

	return nil
}

// DeactivateCompany desactiva una empresa
func (c *CompanyService) DeactivateCompany(ctx context.Context, id string) error {
	// 1. Verificar que la empresa existe
	company, err := c.GetCompanyByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. Verificar que la empresa no está ya desactivada
	if !company.IsActive {
		logs.Warn("Company is already inactive", map[string]interface{}{
			"company_id": id,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "DeactivateCompany", "Company is already inactive", errPackage.ErrCompanyInactive)
	}

	// 3. Desactivar la empresa
	err = c.repo.DeactivateCompany(ctx, id)
	if err != nil {
		logs.Error("Failed to deactivate company", map[string]interface{}{
			"error":      err,
			"company_id": id,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "DeactivateCompany", "Error deactivating company", err)
	}

	logs.Info("Company deactivated successfully", map[string]interface{}{
		"company_id": id,
	})

	return nil
}

// ReactivateCompany reactiva una empresa
func (c *CompanyService) ReactivateCompany(ctx context.Context, id string) error {
	// 1. Verificar que la empresa existe
	company, err := c.repo.GetCompanyByID(ctx, id)
	if err != nil {
		logs.Error("Failed to get company by ID", map[string]interface{}{
			"error": err,
			"id":    id,
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errPackage.NewDomainErrorWithCause("CompanyService", "ReactivateCompany", "Company not found", errPackage.ErrCompanyNotFound)
		}

		return errPackage.NewDomainErrorWithCause("CompanyService", "ReactivateCompany", "Error getting company by ID", err)
	}

	// 2. Verificar que la empresa está desactivada
	if company.IsActive {
		logs.Warn("Company is already active", map[string]interface{}{
			"company_id": id,
		})
		return nil // No es un error, simplemente ya está activa
	}

	// 3. Reactivar la empresa
	err = c.repo.ReactivateCompany(ctx, id)
	if err != nil {
		logs.Error("Failed to reactivate company", map[string]interface{}{
			"error":      err,
			"company_id": id,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "ReactivateCompany", "Error reactivating company", err)
	}

	logs.Info("Company reactivated successfully", map[string]interface{}{
		"company_id": id,
	})

	return nil
}

// ValidateCompany valida una empresa según las reglas de negocio
func (c *CompanyService) ValidateCompany(ctx context.Context, company *entities.Company) error {
	// 1. Validar campos obligatorios
	if company.Name == "" || company.LegalName == "" || company.TaxID == "" ||
		company.ContactEmail == "" || company.ContactPhone == "" {
		logs.Error("Invalid company data", map[string]interface{}{
			"company": company.Name,
		})
		return errPackage.NewDomainError("CompanyService", "ValidateCompany", errPackage.ErrInvalidCompanyData.Error())
	}

	// 2. Validar valores específicos
	// Verificar formato de email
	emailVO := value_objects.NewEmail(company.ContactEmail)
	if !emailVO.IsValid() {
		logs.Error("Invalid email format", map[string]interface{}{
			"email": company.ContactEmail,
		})
		return errPackage.NewDomainError("CompanyService", "ValidateCompany", "Invalid email format")
	}

	// 3. Validar unicidad
	// Verificar que no exista otra empresa con el mismo TaxID
	exists, err := c.repo.ExistsTaxID(ctx, company.TaxID, "")
	if err != nil {
		logs.Error("Failed to check if tax ID exists", map[string]interface{}{
			"error":  err,
			"tax_id": company.TaxID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateCompany", "Error checking tax ID uniqueness", err)
	}

	if exists {
		logs.Error("Tax ID already exists", map[string]interface{}{
			"tax_id": company.TaxID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateCompany", "Tax ID already exists", errPackage.ErrDuplicateTaxID)
	}

	// Verificar que no exista otra empresa con el mismo nombre
	exists, err = c.repo.ExistsCompanyName(ctx, company.Name, "")
	if err != nil {
		logs.Error("Failed to check if company name exists", map[string]interface{}{
			"error": err,
			"name":  company.Name,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateCompany", "Error checking company name uniqueness", err)
	}

	if exists {
		logs.Error("Company name already exists", map[string]interface{}{
			"name": company.Name,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateCompany", "Company name already exists", errPackage.ErrDuplicateCompanyName)
	}

	// 4. Validar fechas de contrato
	if company.ContractEndDate != nil && company.ContractEndDate.Before(company.ContractStartDate) {
		logs.Error("Invalid contract dates", map[string]interface{}{
			"start_date": company.ContractStartDate,
			"end_date":   company.ContractEndDate,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateCompany", "Contract end date cannot be before start date", errPackage.ErrInvalidContractDates)
	}

	// 5. Validar detalles de contrato si existen
	if company.ContractDetails != "" {
		var contractDetails map[string]interface{}
		if err := json.Unmarshal([]byte(company.ContractDetails), &contractDetails); err != nil {
			logs.Error("Invalid contract details format", map[string]interface{}{
				"error":            err,
				"contract_details": company.ContractDetails,
			})
			return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateCompany", "Invalid contract details format", err)
		}
	}

	return nil
}

// ValidateCompanyUpdate valida una actualización de empresa según las reglas de negocio
func (c *CompanyService) ValidateCompanyUpdate(ctx context.Context, company *entities.Company, existingCompany *entities.Company) error {
	// 1. Validar campos que pueden cambiar
	// Verificar formato de email si se actualiza
	if company.ContactEmail != "" && company.ContactEmail != existingCompany.ContactEmail {
		emailVO := value_objects.NewEmail(company.ContactEmail)
		if !emailVO.IsValid() {
			logs.Error("Invalid email format", map[string]interface{}{
				"email": company.ContactEmail,
			})
			return errPackage.NewDomainError("CompanyService", "ValidateCompanyUpdate", "Invalid email format")
		}
	}

	// 2. Validar unicidad si se cambia el nombre
	if company.Name != "" && company.Name != existingCompany.Name {
		exists, err := c.repo.ExistsCompanyName(ctx, company.Name, company.ID)
		if err != nil {
			logs.Error("Failed to check if company name exists", map[string]interface{}{
				"error": err,
				"name":  company.Name,
			})
			return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateCompanyUpdate", "Error checking company name uniqueness", err)
		}

		if exists {
			logs.Error("Company name already exists", map[string]interface{}{
				"name": company.Name,
			})
			return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateCompanyUpdate", "Company name already exists", errPackage.ErrDuplicateCompanyName)
		}
	}

	// 3. Validar fechas de contrato si se actualiza
	if company.ContractEndDate != nil && existingCompany.ContractStartDate.After(time.Time{}) {
		if company.ContractEndDate.Before(existingCompany.ContractStartDate) {
			logs.Error("Invalid contract dates", map[string]interface{}{
				"start_date": existingCompany.ContractStartDate,
				"end_date":   company.ContractEndDate,
			})
			return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateCompanyUpdate", "Contract end date cannot be before start date", errPackage.ErrInvalidContractDates)
		}
	}

	// 4. Validar detalles de contrato si se actualizan
	if company.ContractDetails != "" {
		var contractDetails map[string]interface{}
		if err := json.Unmarshal([]byte(company.ContractDetails), &contractDetails); err != nil {
			logs.Error("Invalid contract details format", map[string]interface{}{
				"error":            err,
				"contract_details": company.ContractDetails,
			})
			return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateCompanyUpdate", "Invalid contract details format", err)
		}
	}

	return nil
}

// GetCompanyBranches obtiene las sucursales de una empresa
func (c *CompanyService) GetCompanyBranches(ctx context.Context, companyID string) ([]entities.Branch, error) {
	// 1. Verificar que la empresa existe
	_, err := c.GetCompanyByID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	// 2. Obtener las sucursales
	branches, err := c.repo.GetCompanyBranches(ctx, companyID)
	if err != nil {
		logs.Error("Failed to get company branches", map[string]interface{}{
			"error":      err,
			"company_id": companyID,
		})
		return nil, errPackage.NewDomainErrorWithCause("CompanyService", "GetCompanyBranches", "Error getting company branches", err)
	}

	return branches, nil
}

// AddBranchToCompany añade una sucursal a una empresa
func (c *CompanyService) AddBranchToCompany(ctx context.Context, companyID string, branch *entities.Branch) error {
	// 1. Verificar que la empresa existe y está activa
	company, err := c.GetCompanyByID(ctx, companyID)
	if err != nil {
		return err
	}

	if !company.IsActive {
		logs.Error("Cannot add branch to inactive company", map[string]interface{}{
			"company_id": companyID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "AddBranchToCompany", "Cannot add branch to inactive company", errPackage.ErrCompanyInactive)
	}

	// 2. Validar la sucursal
	if err := c.ValidateBranch(ctx, branch); err != nil {
		return err
	}

	// 3. Asegurarse de que la sucursal pertenece a la empresa correcta
	branch.CompanyID = companyID

	// 4. Crear la sucursal
	err = c.repo.CreateBranch(ctx, branch)
	if err != nil {
		logs.Error("Failed to create branch", map[string]interface{}{
			"error":       err,
			"branch_name": branch.Name,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "AddBranchToCompany", "Error creating branch", err)
	}

	logs.Info("Branch added to company successfully", map[string]interface{}{
		"company_id": companyID,
		"branch_id":  branch.ID,
	})

	return nil
}

// UpdateBranch actualiza una sucursal
func (c *CompanyService) UpdateBranch(ctx context.Context, branch *entities.Branch) error {
	// 1. Verificar que la sucursal existe
	existingBranch, err := c.repo.GetBranchByID(ctx, branch.ID)
	if err != nil {
		logs.Error("Failed to get branch by ID", map[string]interface{}{
			"error":     err,
			"branch_id": branch.ID,
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errPackage.NewDomainErrorWithCause("CompanyService", "UpdateBranch", "Branch not found", errPackage.ErrBranchNotFound)
		}

		return errPackage.NewDomainErrorWithCause("CompanyService", "UpdateBranch", "Error getting branch by ID", err)
	}

	// 2. Validar la sucursal actualizada
	if err := c.ValidateBranchUpdate(ctx, branch, existingBranch); err != nil {
		return err
	}

	// 3. Asegurarse de que no se cambia la compañía
	if branch.CompanyID != "" && branch.CompanyID != existingBranch.CompanyID {
		logs.Error("Cannot change branch company", map[string]interface{}{
			"branch_id":          branch.ID,
			"current_company_id": existingBranch.CompanyID,
			"new_company_id":     branch.CompanyID,
		})
		return errPackage.NewDomainError("CompanyService", "UpdateBranch", errPackage.ErrCannotChangeCompany.Error())
	}

	// 4. Actualizar la sucursal
	err = c.repo.UpdateBranch(ctx, branch)
	if err != nil {
		logs.Error("Failed to update branch", map[string]interface{}{
			"error":     err,
			"branch_id": branch.ID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "UpdateBranch", "Error updating branch", err)
	}

	logs.Info("Branch updated successfully", map[string]interface{}{
		"branch_id": branch.ID,
	})

	return nil
}

// DeactivateBranch desactiva una sucursal
func (c *CompanyService) DeactivateBranch(ctx context.Context, branchID string) error {
	// 1. Verificar que la sucursal existe
	branch, err := c.repo.GetBranchByID(ctx, branchID)
	if err != nil {
		logs.Error("Failed to get branch by ID", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errPackage.NewDomainErrorWithCause("CompanyService", "DeactivateBranch", "Branch not found", errPackage.ErrBranchNotFound)
		}

		return errPackage.NewDomainErrorWithCause("CompanyService", "DeactivateBranch", "Error getting branch by ID", err)
	}

	// 2. Verificar que la sucursal no está ya desactivada
	if !branch.IsActive {
		logs.Warn("Branch is already inactive", map[string]interface{}{
			"branch_id": branchID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "DeactivateBranch", "Branch is already inactive", errPackage.ErrBranchInactive)
	}

	// 3. Desactivar la sucursal
	err = c.repo.DeactivateBranch(ctx, branchID)
	if err != nil {
		logs.Error("Failed to deactivate branch", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "DeactivateBranch", "Error deactivating branch", err)
	}

	logs.Info("Branch deactivated successfully", map[string]interface{}{
		"branch_id": branchID,
	})

	return nil
}

// ReactivateBranch reactiva una sucursal
func (c *CompanyService) ReactivateBranch(ctx context.Context, branchID string) error {
	// 1. Verificar que la sucursal existe
	branch, err := c.repo.GetBranchByID(ctx, branchID)
	if err != nil {
		logs.Error("Failed to get branch by ID", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errPackage.NewDomainErrorWithCause("CompanyService", "ReactivateBranch", "Branch not found", errPackage.ErrBranchNotFound)
		}

		return errPackage.NewDomainErrorWithCause("CompanyService", "ReactivateBranch", "Error getting branch by ID", err)
	}

	// 2. Verificar que la sucursal no está ya desactivada
	if !branch.IsActive {
		logs.Warn("Branch is already inactive", map[string]interface{}{
			"branch_id": branchID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "ReactivateBranch", "Branch is already inactive", errPackage.ErrBranchInactive)
	}

	// 3. Desactivar la sucursal
	err = c.repo.ReactivateBranch(ctx, branchID)
	if err != nil {
		logs.Error("Failed to deactivate branch", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "ReactivateBranch", "Error deactivating branch", err)
	}

	logs.Info("Branch reactivated successfully", map[string]interface{}{
		"branch_id": branchID,
	})

	return nil
}

// ValidateBranch valida una sucursal según las reglas de negocio
func (c *CompanyService) ValidateBranch(ctx context.Context, branch *entities.Branch) error {
	// 1. Validar campos obligatorios
	if branch.Name == "" || branch.Code == "" || branch.ContactName == "" ||
		branch.ContactPhone == "" || branch.ContactEmail == "" || branch.ZoneID == "" {
		logs.Error("Invalid branch data", map[string]interface{}{
			"branch": branch.Name,
		})
		return errPackage.NewDomainError("CompanyService", "ValidateBranch", errPackage.ErrInvalidBranchData.Error())
	}

	// 2. Validar valores específicos
	// Verificar formato de email
	emailVO := value_objects.NewEmail(branch.ContactEmail)
	if !emailVO.IsValid() {
		logs.Error("Invalid email format", map[string]interface{}{
			"email": branch.ContactEmail,
		})
		return errPackage.NewDomainError("CompanyService", "ValidateBranch", "Invalid email format")
	}

	// 3. Validar que la zona existe y está activa
	zone, err := c.repo.GetZoneByID(ctx, branch.ZoneID)
	if err != nil {
		logs.Error("Failed to get zone by ID", map[string]interface{}{
			"error":   err,
			"zone_id": branch.ZoneID,
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateBranch", "Zone not found", errPackage.ErrZoneNotFound)
		}

		return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateBranch", "Error getting zone by ID", err)
	}

	if !zone.IsActive {
		logs.Error("Zone is inactive", map[string]interface{}{
			"zone_id": branch.ZoneID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateBranch", "Zone is inactive", errPackage.ErrZoneInactive)
	}

	// 4. Validar horarios de operación si existen
	if branch.OperatingHours != "" {
		var operatingHours struct {
			Weekdays struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"weekdays"`
			Weekends struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"weekends"`
		}

		if err := json.Unmarshal([]byte(branch.OperatingHours), &operatingHours); err != nil {
			logs.Error("Invalid operating hours format", map[string]interface{}{
				"error":           err,
				"operating_hours": branch.OperatingHours,
			})
			return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateBranch", "Invalid operating hours format", err)
		}

		// Verificar que los horarios tengan el formato correcto (HH:MM)
		operatingHoursVO, err := value_objects.NewOperatingHoursFromJSON(branch.OperatingHours)
		if err != nil {
			logs.Error("Invalid operating hours format", map[string]interface{}{
				"error": err,
			})
			return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateBranch", "Invalid operating hours format", err)
		}
		if operatingHoursVO == nil || !operatingHoursVO.IsValid() {
			logs.Error("Invalid operating hours", map[string]interface{}{
				"operating_hours": branch.OperatingHours,
			})
			return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateBranch", "Invalid operating hours", errPackage.ErrInvalidOperatingHours)
		}
	}

	return nil
}

// ValidateBranchUpdate valida una actualización de sucursal según las reglas de negocio
func (c *CompanyService) ValidateBranchUpdate(ctx context.Context, branch *entities.Branch, existingBranch *entities.Branch) error {
	// 1. Validar campos que pueden cambiar
	// Verificar formato de email si se actualiza
	if branch.ContactEmail != "" && branch.ContactEmail != existingBranch.ContactEmail {
		emailVO := value_objects.NewEmail(branch.ContactEmail)
		if !emailVO.IsValid() {
			logs.Error("Invalid email format", map[string]interface{}{
				"email": branch.ContactEmail,
			})
			return errPackage.NewDomainError("CompanyService", "ValidateBranchUpdate", "Invalid email format")
		}
	}

	// 2. Validar que la zona existe y está activa si se cambia
	if branch.ZoneID != "" && branch.ZoneID != existingBranch.ZoneID {
		// Validar zona
		err := c.ValidateBranchZoneAssignment(ctx, branch.ID, branch.ZoneID)
		if err != nil {
			return err
		}
	}

	// 3. Validar horarios de operación si se actualizan
	if branch.OperatingHours != "" && branch.OperatingHours != existingBranch.OperatingHours {
		operatingHoursVO, err := value_objects.NewOperatingHoursFromJSON(branch.OperatingHours)
		if err != nil {
			logs.Error("Invalid operating hours format", map[string]interface{}{
				"error": err,
			})
			return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateBranchUpdate", "Invalid operating hours format", err)
		}
		if operatingHoursVO == nil || !operatingHoursVO.IsValid() {
			logs.Error("Invalid operating hours", map[string]interface{}{
				"operating_hours": branch.OperatingHours,
			})
			return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateBranchUpdate", "Invalid operating hours", errPackage.ErrInvalidOperatingHours)
		}
	}

	return nil
}

// AssignZoneToBranch asigna una zona a una sucursal
func (c *CompanyService) AssignZoneToBranch(ctx context.Context, branchID string, zoneID string) error {
	// 1. Verificar que la sucursal existe
	branch, err := c.repo.GetBranchByID(ctx, branchID)
	if err != nil {
		logs.Error("Failed to get branch by ID", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errPackage.NewDomainErrorWithCause("CompanyService", "AssignZoneToBranch", "Branch not found", errPackage.ErrBranchNotFound)
		}

		return errPackage.NewDomainErrorWithCause("CompanyService", "AssignZoneToBranch", "Error getting branch by ID", err)
	}

	// 2. Validar la asignación de zona
	if err := c.ValidateBranchZoneAssignment(ctx, branchID, zoneID); err != nil {
		return err
	}

	// 3. Actualizar la zona de la sucursal
	branch.ZoneID = zoneID
	branch.UpdatedAt = time.Now()

	err = c.repo.UpdateBranch(ctx, branch)
	if err != nil {
		logs.Error("Failed to update branch zone", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
			"zone_id":   zoneID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "AssignZoneToBranch", "Error updating branch zone", err)
	}

	logs.Info("Zone assigned to branch successfully", map[string]interface{}{
		"branch_id": branchID,
		"zone_id":   zoneID,
	})

	return nil
}

// ValidateBranchZoneAssignment valida la asignación de una zona a una sucursal
func (c *CompanyService) ValidateBranchZoneAssignment(ctx context.Context, branchID string, zoneID string) error {
	// 1. Verificar que la zona existe y está activa
	zone, err := c.repo.GetZoneByID(ctx, zoneID)
	if err != nil {
		logs.Error("Failed to get zone by ID", map[string]interface{}{
			"error":   err,
			"zone_id": zoneID,
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateBranchZoneAssignment", "Zone not found", errPackage.ErrZoneNotFound)
		}

		return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateBranchZoneAssignment", "Error getting zone by ID", err)
	}

	if !zone.IsActive {
		logs.Error("Zone is inactive", map[string]interface{}{
			"zone_id": zoneID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateBranchZoneAssignment", "Zone is inactive", errPackage.ErrZoneInactive)
	}

	// 2. Verificar que la sucursal existe (si se está actualizando una sucursal existente)
	if branchID != "" {
		branch, err := c.repo.GetBranchByID(ctx, branchID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				logs.Error("Failed to get branch by ID", map[string]interface{}{
					"error":     err,
					"branch_id": branchID,
				})
				return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateBranchZoneAssignment", "Error getting branch by ID", err)
			}
			// Si la sucursal no existe, es una sucursal nueva, así que continuamos
		} else {
			// 3. Verificar reglas específicas de negocio para la asignación de zonas
			// Por ejemplo, podríamos limitar el número de sucursales por zona por empresa

			// Obtener todas las sucursales en la zona para la misma empresa
			branchesInZone, err := c.repo.GetBranchesByZone(ctx, zoneID)
			if err != nil {
				logs.Error("Failed to get branches by zone", map[string]interface{}{
					"error":   err,
					"zone_id": zoneID,
				})
				return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateBranchZoneAssignment", "Error getting branches by zone", err)
			}

			// Contar sucursales de la misma empresa en la zona
			companyBranchesInZone := 0
			for _, b := range branchesInZone {
				if b.CompanyID == branch.CompanyID && b.ID != branchID {
					companyBranchesInZone++
				}
			}

			// Ejemplo: limitar a 3 sucursales por zona por empresa
			if companyBranchesInZone >= 3 {
				logs.Error("Company already has maximum number of branches in this zone", map[string]interface{}{
					"company_id":       branch.CompanyID,
					"zone_id":          zoneID,
					"current_branches": companyBranchesInZone,
				})
				return errPackage.NewDomainErrorWithCause("CompanyService", "ValidateBranchZoneAssignment", "Company already has maximum number of branches in this zone", errPackage.ErrTooManyBranchesInZone)
			}
		}
	}

	return nil
}

// GetAvailableZonesForBranch obtiene las zonas disponibles para una sucursal
func (c *CompanyService) GetAvailableZonesForBranch(ctx context.Context, branchID string) ([]entities.Zone, error) {
	// 1. Verificar que la sucursal existe (si se proporciona un ID)
	var branch *entities.Branch
	var err error

	if branchID != "" {
		branch, err = c.repo.GetBranchByID(ctx, branchID)
		if err != nil {
			logs.Error("Failed to get branch by ID", map[string]interface{}{
				"error":     err,
				"branch_id": branchID,
			})

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errPackage.NewDomainErrorWithCause("CompanyService", "GetAvailableZonesForBranch", "Branch not found", errPackage.ErrBranchNotFound)
			}

			return nil, errPackage.NewDomainErrorWithCause("CompanyService", "GetAvailableZonesForBranch", "Error getting branch by ID", err)
		}
	}

	// 2. Obtener todas las zonas activas
	allZones, err := c.repo.GetAllActiveZones(ctx)
	if err != nil {
		logs.Error("Failed to get active zones", map[string]interface{}{
			"error": err,
		})
		return nil, errPackage.NewDomainErrorWithCause("CompanyService", "GetAvailableZonesForBranch", "Error getting active zones", err)
	}

	// 3. Si no hay un ID de sucursal, devolver todas las zonas activas
	if branchID == "" || branch == nil {
		return allZones, nil
	}

	// 4. Filtrar zonas según las reglas de negocio para la sucursal y su empresa
	var availableZones []entities.Zone
	for _, zone := range allZones {
		// Verificar si la sucursal puede ser asignada a esta zona
		err := c.ValidateBranchZoneAssignment(ctx, branchID, zone.ID)
		if err == nil {
			// Si no hay error, la zona está disponible
			availableZones = append(availableZones, zone)
		}
	}

	return availableZones, nil
}

func (c *CompanyService) GetCompanyMetrics(ctx context.Context, companyID string) (*entities.CompanyMetrics, error) {

	// 1. Obtener las métricas de la empresa
	metrics, err := c.metricsService.GetCompanyMetrics(ctx, companyID)
	if err != nil {
		logs.Error("Failed to calculate company metrics", map[string]interface{}{
			"error":      err,
			"company_id": companyID,
		})
		return nil, errPackage.NewDomainErrorWithCause("CompanyService", "GetCompanyMetrics", "Error calculating company metrics", err)
	}

	return metrics, nil
}

func (c *CompanyService) GetBranchMetrics(ctx context.Context, branchID string) (*entities.BranchMetrics, error) {
	// 1. Obtener los claims del contexto
	claims, ok := ctx.Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Claims not found in context",
		})
		return nil, errPackage.NewDomainError("CompanyService", "GetBranchMetrics", "Claims not found in context")
	}

	branch, err := c.repo.GetBranchByID(ctx, branchID)
	if err != nil {
		logs.Error("Failed to get branch by ID", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
		})
		return nil, errPackage.NewDomainErrorWithCause("CompanyService", "GetAvailableZonesForBranch", "Error getting branch by ID", err)
	}
	if branch.CompanyID != claims.CompanyID {
		return nil, errPackage.NewDomainErrorWithCause("CompanyService", "GetBranchMetrics", "Branch not found", errPackage.ErrBranchNotFound)
	}

	// 2. Obtener las métricas de la sucursal
	metrics, err := c.metricsService.GetBranchMetrics(ctx, branchID, claims.CompanyID)
	if err != nil {
		logs.Error("Failed to calculate branch metrics", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
		})
		return nil, errPackage.NewDomainErrorWithCause("CompanyService", "GetBranchMetrics", "Error calculating branch metrics", err)
	}

	return metrics, nil
}

func (c *CompanyService) AddCompanyAddress(ctx context.Context, companyID string, address *entities.CompanyAddress) error {

	// 1. Verificar que la empresa existe
	_, err := c.repo.GetCompanyByID(ctx, companyID)
	if err != nil {
		logs.Error("Failed to get company by ID", map[string]interface{}{
			"error":      err,
			"company_id": companyID,
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errPackage.NewDomainErrorWithCause("CompanyService", "AddCompanyAddress", "Company not found", errPackage.ErrCompanyNotFound)
		}

		return errPackage.NewDomainErrorWithCause("CompanyService", "AddCompanyAddress", "Error getting company by ID", err)
	}

	// 2. Validar los datos proporcionados para la dirección
	if address == nil || address.AddressLine1 == "" || address.City == "" {
		logs.Error("Invalid company address data", map[string]interface{}{
			"company_id": companyID,
			"address":    address,
		})
		return errPackage.NewDomainError("CompanyService", "AddCompanyAddress", "Invalid address data")
	}

	// 3. Guardar la dirección en el repositorio
	address.CompanyID = companyID
	err = c.repo.AddCompanyAddress(ctx, address)
	if err != nil {
		logs.Error("Failed to create company address", map[string]interface{}{
			"error":      err,
			"company_id": companyID,
			"address":    address,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "AddCompanyAddress", "Error creating company address", err)
	}

	return nil
}

func (c *CompanyService) UpdateCompanyAddress(ctx context.Context, companyID string, address *entities.CompanyAddress) error {

	// 1. Verificar que la dirección proporcionada es válida
	if address == nil || address.ID == "" || address.AddressLine1 == "" || address.City == "" {
		logs.Error("Invalid company address data", map[string]interface{}{
			"address": address,
		})
		return errPackage.NewDomainError("CompanyService", "UpdateCompanyAddress", "Invalid address data")
	}

	// 2. Verificar que la dirección existe en la base de datos
	existingAddress, err := c.repo.GetCompanyAddressByID(ctx, address.ID, companyID)
	if err != nil {
		logs.Error("Failed to get company address by ID", map[string]interface{}{
			"error":   err,
			"address": address,
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errPackage.NewDomainErrorWithCause("CompanyService", "UpdateCompanyAddress", "Address not found", errPackage.ErrAddressNotFound)
		}

		return errPackage.NewDomainErrorWithCause("CompanyService", "UpdateCompanyAddress", "Error getting address by ID", err)
	}

	// 3. Verificar que la dirección está asociada a una empresa válida
	if existingAddress.CompanyID == "" {
		logs.Error("Address is not associated with a company", map[string]interface{}{
			"address_id": address.ID,
		})
		return errPackage.NewDomainError("CompanyService", "UpdateCompanyAddress", "Address is not associated with a company")
	}

	// 4. Actualizar la dirección en el repositorio
	err = c.repo.UpdateCompanyAddress(ctx, address)
	if err != nil {
		logs.Error("Failed to update company address", map[string]interface{}{
			"error":      err,
			"address_id": address.ID,
			"address":    address,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "UpdateCompanyAddress", "Error updating company address", err)
	}

	return nil
}

func (c *CompanyService) DeleteCompanyAddress(ctx context.Context, addressID, companyID string) error {
	// 1.  Verificar que la dirección existe en la base de datos
	_, err := c.repo.GetCompanyAddressByID(ctx, addressID, companyID)
	if err != nil {
		logs.Error("Failed to get company address by ID", map[string]interface{}{
			"error":   err,
			"address": addressID,
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errPackage.NewDomainErrorWithCause("CompanyService", "UpdateCompanyAddress", "Address not found", errPackage.ErrAddressNotFound)
		}

		return errPackage.NewDomainErrorWithCause("CompanyService", "UpdateCompanyAddress", "Error getting address by ID", err)
	}

	// 2. Eliminar la direccion
	err = c.repo.DeleteCompanyAddress(ctx, addressID)
	if err != nil {
		logs.Error("Failed to delete company address", map[string]interface{}{
			"error":   err,
			"address": addressID,
		})
		return errPackage.NewDomainErrorWithCause("CompanyService", "UpdateCompanyAddress", "Error deleting company address", err)
	}

	return nil
}

func (c *CompanyService) GetBranchByID(ctx context.Context, branchID string) (*entities.Branch, error) {
	branch, err := c.repo.GetBranchByID(ctx, branchID)
	if err != nil {
		logs.Error("Failed to get branch by ID", map[string]interface{}{
			"error":     err,
			"branch_id": branchID,
		})
	}

	return branch, err
}

func (c *CompanyService) GetCompanies(ctx context.Context, params *entities.CompanyQueryParams) ([]entities.Company, int64, error) {
	companies, total, err := c.repo.GetCompanies(ctx, params)
	if err != nil {
		logs.Error("Failed to get companies", map[string]interface{}{
			"error": err,
		})
		return nil, 0, errPackage.NewDomainErrorWithCause("CompanyService", "GetCompanies", "Error getting companies", err)
	}

	return companies, total, nil
}
