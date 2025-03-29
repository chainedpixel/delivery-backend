package ports

import (
	"context"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
)

type CompanyAddreser interface {
	GetCompanyAddresses(ctx context.Context) ([]entities.CompanyAddress, error)
	GetCompanyAddressByID(ctx context.Context, id string) (*entities.CompanyAddress, error)
	GetCompanyAndBranchForUser(ctx context.Context, userID string) (string, string, error)
}
