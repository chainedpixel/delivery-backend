package services

import (
	"context"
	"domain/delivery/interfaces"
	"domain/delivery/models/entities"
	"domain/delivery/ports"
	errPackage "domain/error"
	"errors"
	"gorm.io/gorm"
	"shared/logs"
)

type CompanyService struct {
	repo ports.CompanyAddreser
}

func NewCompanyService(repo ports.CompanyAddreser) interfaces.Companyrer {
	return &CompanyService{repo: repo}
}

func (c *CompanyService) GetAddresses(ctx context.Context) ([]entities.CompanyAddress, error) {
	addresses, err := c.repo.GetCompanyAddresses(ctx)
	if err != nil {
		logs.Error("Failed to get company addresses", map[string]interface{}{
			"error": err,
		})
		return nil, errPackage.NewDomainErrorWithCause("CompanyService", "GetAddresses", "Error getting company addresses", err)
	}

	return addresses, nil
}

func (c *CompanyService) GetAddressByID(ctx context.Context, id string) (*entities.CompanyAddress, error) {
	address, err := c.repo.GetCompanyAddressByID(ctx, id)
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
		logs.Error("Failed to get company ID by user ID", map[string]interface{}{
			"error":   err,
			"user_id": userID,
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", errPackage.NewDomainErrorWithCause("CompanyService", "GetCompanyIDByUserID", "Error getting company ID by user ID", errPackage.ErrUserNotFoundOrUnauthorized)
		}

		return "", "", errPackage.NewDomainErrorWithCause("CompanyService", "GetCompanyIDByUserID", "Error getting company ID by user ID", err)
	}

	return companyID, branchID, nil
}
