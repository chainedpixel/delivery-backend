package ports

import (
	"context"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"net/http"
)

// Definición de la interfaz de casos de uso para Companies
type CompanyUseCase interface {
	// CRUD básico
	GetCompanyByID(ctx context.Context) (*entities.Company, error)
	CreateCompany(ctx context.Context, company *entities.Company) error
	UpdateCompany(ctx context.Context, company *entities.Company) error
	DeactivateCompany(ctx context.Context, companyID string) error
	ReactivateCompany(ctx context.Context, companyID string) error

	// Operaciones relacionadas con direcciones
	GetCompanyAddresses(ctx context.Context) ([]entities.CompanyAddress, error)
	AddCompanyAddress(ctx context.Context, address *entities.CompanyAddress) error
	UpdateCompanyAddress(ctx context.Context, addressID string, address *entities.CompanyAddress) error
	DeleteCompanyAddress(ctx context.Context, addressID string) error

	// Operaciones relacionadas con métricas
	GetCompanyMetrics(ctx context.Context) (*entities.CompanyMetrics, error)
}

// Definición de la interfaz de casos de uso para Branches
type BranchUseCase interface {
	// CRUD básico
	GetBranches(ctx context.Context, request *http.Request) ([]entities.Branch, *entities.BranchQueryParams, int64, error)
	GetBranchByID(ctx context.Context, branchID string) (*entities.Branch, error)
	GetBranchesByCompany(ctx context.Context, companyID string) ([]entities.Branch, error)
	CreateBranch(ctx context.Context, companyID string, branch *entities.Branch) error
	UpdateBranch(ctx context.Context, branchID string, branch *entities.Branch) error
	DeactivateBranch(ctx context.Context, branchID string) error
	ReactivateBranch(ctx context.Context, branchID string) error

	// Operaciones relacionadas con zonas
	AssignZoneToBranch(ctx context.Context, branchID string, zoneID string) error
	GetAvailableZonesForBranch(ctx context.Context, branchID string) ([]entities.Zone, error)

	// Operaciones relacionadas con métricas
	GetBranchMetrics(ctx context.Context, branchID string) (*entities.BranchMetrics, error)
}
