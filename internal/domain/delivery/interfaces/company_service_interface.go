package interfaces

import (
	"context"

	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
)

type Companyrer interface {
	GetAddresses(ctx context.Context) ([]entities.CompanyAddress, error)
	GetAddressByID(ctx context.Context, id string) (*entities.CompanyAddress, error)
	GetCompanyAndBranchForUser(ctx context.Context, userID string) (string, string, error)
}
