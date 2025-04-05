package ports

import (
	"context"

	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
)

type CompanyRepository interface {
	// Métodos para direcciones
	GetCompanyAddresses(ctx context.Context, companyID string) ([]entities.CompanyAddress, error)
	GetCompanyAddressByID(ctx context.Context, id, companyID string) (*entities.CompanyAddress, error)
	GetCompanyAndBranchForUser(ctx context.Context, userID string) (string, string, error)

	// Métodos para CRUD de Company
	GetCompanyByID(ctx context.Context, id string) (*entities.Company, error)
	CreateCompany(ctx context.Context, company *entities.Company) error
	UpdateCompany(ctx context.Context, company *entities.Company) error
	DeactivateCompany(ctx context.Context, id string) error
	ReactivateCompany(ctx context.Context, id string) error
	AddCompanyAddress(ctx context.Context, address *entities.CompanyAddress) error
	UpdateCompanyAddress(ctx context.Context, address *entities.CompanyAddress) error
	DeleteCompanyAddress(ctx context.Context, addressID string) error
	GetCompanies(ctx context.Context, params *entities.CompanyQueryParams) ([]entities.Company, int64, error)

	// Métodos para verificaciones
	ExistsTaxID(ctx context.Context, taxID string, excludeID string) (bool, error)
	ExistsCompanyName(ctx context.Context, name string, excludeID string) (bool, error)

	// Métodos para gestión de sucursales (branches)
	GetCompanyBranches(ctx context.Context, companyID string) ([]entities.Branch, error)
	GetBranchByID(ctx context.Context, branchID string) (*entities.Branch, error)
	CreateBranch(ctx context.Context, branch *entities.Branch) error
	UpdateBranch(ctx context.Context, branch *entities.Branch) error
	DeactivateBranch(ctx context.Context, branchID string) error
	ReactivateBranch(ctx context.Context, branchID string) error

	// Métodos para zonas
	GetZoneByID(ctx context.Context, zoneID string) (*entities.Zone, error)
	GetAllActiveZones(ctx context.Context) ([]entities.Zone, error)
	GetBranchesByZone(ctx context.Context, zoneID string) ([]entities.Branch, error)
}
