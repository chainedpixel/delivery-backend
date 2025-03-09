package interfaces

import (
	"context"
	"domain/delivery/models/entities"
)

type Companyrer interface {
	GetAddresses(ctx context.Context) ([]entities.CompanyAddress, error)
	GetAddressByID(ctx context.Context, id string) (*entities.CompanyAddress, error)
}
