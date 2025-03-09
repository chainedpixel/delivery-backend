package repositories

import (
	"context"
	"domain/delivery/models/entities"
	"domain/delivery/ports"
	"gorm.io/gorm"
)

type CompanyAddressRepository struct {
	db *gorm.DB
}

func NewCompanyAddressRepository(db *gorm.DB) ports.CompanyAddreser {
	return &CompanyAddressRepository{
		db: db,
	}
}

func (r *CompanyAddressRepository) GetCompanyAddresses(ctx context.Context) ([]entities.CompanyAddress, error) {
	var companyAddresses []entities.CompanyAddress
	err := r.db.WithContext(ctx).Find(&companyAddresses).Error
	if err != nil {
		return nil, err
	}

	return companyAddresses, nil
}

func (r *CompanyAddressRepository) GetCompanyAddressByID(ctx context.Context, id string) (*entities.CompanyAddress, error) {
	var companyAddress entities.CompanyAddress
	err := r.db.WithContext(ctx).First(&companyAddress, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &companyAddress, nil
}
