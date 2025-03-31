package interfaces

import (
	"context"

	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
)

type Companyrer interface {
	// Métodos existentes
	GetAddresses(ctx context.Context, companyID string) ([]entities.CompanyAddress, error)
	GetAddressByID(ctx context.Context, id, companyID string) (*entities.CompanyAddress, error)
	GetCompanyAndBranchForUser(ctx context.Context, userID string) (string, string, error)
	GetCompanyByID(ctx context.Context, id string) (*entities.Company, error)

	// CRUD de Company
	GetBranchByID(ctx context.Context, branchID string) (*entities.Branch, error)
	CreateCompany(ctx context.Context, company *entities.Company) error
	UpdateCompany(ctx context.Context, company *entities.Company) error
	DeactivateCompany(ctx context.Context, id string) error
	ReactivateCompany(ctx context.Context, id string) error
	ValidateCompany(ctx context.Context, company *entities.Company) error
	ValidateCompanyUpdate(ctx context.Context, company *entities.Company, existingCompany *entities.Company) error
	AddCompanyAddress(ctx context.Context, companyID string, address *entities.CompanyAddress) error
	UpdateCompanyAddress(ctx context.Context, companyID string, address *entities.CompanyAddress) error
	DeleteCompanyAddress(ctx context.Context, addressID, companyID string) error

	// CRUD de Branch
	GetCompanyBranches(ctx context.Context, companyID string) ([]entities.Branch, error)
	AddBranchToCompany(ctx context.Context, companyID string, branch *entities.Branch) error
	UpdateBranch(ctx context.Context, branch *entities.Branch) error
	DeactivateBranch(ctx context.Context, branchID string) error
	ReactivateBranch(ctx context.Context, branchID string) error
	ValidateBranch(ctx context.Context, branch *entities.Branch) error
	ValidateBranchUpdate(ctx context.Context, branch *entities.Branch, existingBranch *entities.Branch) error

	// Gestión de zonas
	AssignZoneToBranch(ctx context.Context, branchID string, zoneID string) error
	ValidateBranchZoneAssignment(ctx context.Context, branchID string, zoneID string) error
	GetAvailableZonesForBranch(ctx context.Context, branchID string) ([]entities.Zone, error)

	// Métricas
	GetCompanyMetrics(ctx context.Context, companyID string) (*entities.CompanyMetrics, error)
	GetBranchMetrics(ctx context.Context, branchID string) (*entities.BranchMetrics, error)
}
