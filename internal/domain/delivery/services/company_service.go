package services

import (
	"context"
	"domain/delivery/interfaces"
	"domain/delivery/models/entities"
	"domain/delivery/ports"
	error2 "domain/error"
	"shared/logs"
)

type CompanyService struct {
	repo ports.CompanyAddreser
}

func NewCompanyService(repo ports.CompanyAddreser) interfaces.Companyrer {
	return &CompanyService{repo: repo}
}

func (c CompanyService) GetAddresses(ctx context.Context) ([]entities.CompanyAddress, error) {
	addresses, err := c.repo.GetCompanyAddresses(ctx)
	if err != nil {
		logs.Error("Failed to get company addresses", map[string]interface{}{
			"error": err,
		})
		return nil, error2.NewDomainErrorWithCause("CompanyService", "GetAddresses", "Error getting company addresses", err)
	}

	return addresses, nil
}

func (c CompanyService) GetAddressByID(ctx context.Context, id string) (*entities.CompanyAddress, error) {
	address, err := c.repo.GetCompanyAddressByID(ctx, id)
	if err != nil {
		logs.Error("Failed to get company address by ID", map[string]interface{}{
			"error": err,
			"id":    id,
		})
		return nil, error2.NewDomainErrorWithCause("CompanyService", "GetAddressByID", "Error getting company address by ID", err)
	}

	return address, nil
}
