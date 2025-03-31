package repositories

import (
	"context"
	"time"

	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/ports"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) ports.CompanyRepository {
	return &CompanyRepository{
		db: db,
	}
}

func (r *CompanyRepository) GetCompanyAddresses(ctx context.Context, companyID string) ([]entities.CompanyAddress, error) {
	var companyAddresses []entities.CompanyAddress
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Find(&companyAddresses).Error
	if err != nil {
		return nil, err
	}

	return companyAddresses, nil
}

func (r *CompanyRepository) GetCompanyAddressByID(ctx context.Context, companyID, id string) (*entities.CompanyAddress, error) {
	var companyAddress entities.CompanyAddress
	err := r.db.WithContext(ctx).First(&companyAddress, "id = ? AND company_id = ?", id, companyID).Error
	if err != nil {
		return nil, err
	}

	return &companyAddress, nil
}

func (r *CompanyRepository) GetCompanyAndBranchForUser(ctx context.Context, userID string) (string, string, error) {
	var company entities.CompanyUser
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&company).Error
	if err != nil {
		return "", "", err
	}

	return company.CompanyID, company.BranchID, nil
}

// Nuevos métodos para CRUD de Company
func (r *CompanyRepository) GetCompanyByID(ctx context.Context, id string) (*entities.Company, error) {
	var company entities.Company
	err := r.db.WithContext(ctx).
		Preload("Address").
		Preload("Branches").
		First(&company, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &company, nil
}

func (r *CompanyRepository) CreateCompany(ctx context.Context, company *entities.Company) error {
	return r.db.WithContext(ctx).Create(company).Error
}

func (r *CompanyRepository) UpdateCompany(ctx context.Context, company *entities.Company) error {
	return r.db.WithContext(ctx).Save(company).Error
}

func (r *CompanyRepository) DeactivateCompany(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&entities.Company{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active":  false,
			"updated_at": time.Now(),
		}).Error
}

func (r *CompanyRepository) ReactivateCompany(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&entities.Company{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active":  true,
			"updated_at": time.Now(),
		}).Error
}

// Métodos para verificaciones
func (r *CompanyRepository) ExistsTaxID(ctx context.Context, taxID string, excludeID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Company{}).
		Where("tax_id = ?", taxID)

	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *CompanyRepository) ExistsCompanyName(ctx context.Context, name string, excludeID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Company{}).
		Where("name = ?", name)

	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Métodos para gestión de sucursales (branches)
func (r *CompanyRepository) GetCompanyBranches(ctx context.Context, companyID string) ([]entities.Branch, error) {
	var branches []entities.Branch
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Find(&branches).Error
	if err != nil {
		return nil, err
	}

	return branches, nil
}

func (r *CompanyRepository) GetBranchByID(ctx context.Context, branchID string) (*entities.Branch, error) {
	var branch entities.Branch
	err := r.db.WithContext(ctx).
		Preload("Company").
		Preload("Zone").
		First(&branch, "id = ?", branchID).Error
	if err != nil {
		return nil, err
	}

	return &branch, nil
}

func (r *CompanyRepository) CreateBranch(ctx context.Context, branch *entities.Branch) error {
	return r.db.WithContext(ctx).Create(branch).Error
}

func (r *CompanyRepository) UpdateBranch(ctx context.Context, branch *entities.Branch) error {
	return r.db.WithContext(ctx).Save(branch).Error
}

func (r *CompanyRepository) DeactivateBranch(ctx context.Context, branchID string) error {
	return r.db.WithContext(ctx).Model(&entities.Branch{}).
		Where("id = ?", branchID).
		Updates(map[string]interface{}{
			"is_active":  false,
			"updated_at": time.Now(),
		}).Error
}

func (r *CompanyRepository) ReactivateBranch(ctx context.Context, branchID string) error {
	return r.db.WithContext(ctx).Model(&entities.Branch{}).
		Where("id = ?", branchID).
		Updates(map[string]interface{}{
			"is_active":  true,
			"updated_at": time.Now(),
		}).Error
}

// Métodos para zonas
func (r *CompanyRepository) GetZoneByID(ctx context.Context, zoneID string) (*entities.Zone, error) {
	var zone entities.Zone
	err := r.db.WithContext(ctx).
		First(&zone, "id = ?", zoneID).Error
	if err != nil {
		return nil, err
	}

	return &zone, nil
}

func (r *CompanyRepository) GetAllActiveZones(ctx context.Context) ([]entities.Zone, error) {
	var zones []entities.Zone
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&zones).Error
	if err != nil {
		return nil, err
	}

	return zones, nil
}

func (r *CompanyRepository) GetBranchesByZone(ctx context.Context, zoneID string) ([]entities.Branch, error) {
	var branches []entities.Branch
	err := r.db.WithContext(ctx).
		Where("zone_id = ?", zoneID).
		Find(&branches).Error
	if err != nil {
		return nil, err
	}

	return branches, nil
}

func (r *CompanyRepository) AddCompanyAddress(ctx context.Context, address *entities.CompanyAddress) error {
	return r.db.WithContext(ctx).Create(address).Error
}

func (r *CompanyRepository) UpdateCompanyAddress(ctx context.Context, address *entities.CompanyAddress) error {
	return r.db.WithContext(ctx).Save(address).Error
}

func (r *CompanyRepository) DeleteCompanyAddress(ctx context.Context, addressID string) error {
	return r.db.WithContext(ctx).Delete(&entities.CompanyAddress{}, "id = ?", addressID).Error
}
