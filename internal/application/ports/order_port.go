package ports

import (
	"context"
	"infrastructure/api/dto"
)

type OrdererUseCase interface {
	CreateOrder(ctx context.Context, reqOrder *dto.OrderCreateRequest) error
}
