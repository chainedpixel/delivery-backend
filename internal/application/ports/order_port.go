package ports

import (
	"context"
	"domain/delivery/models/entities"
	"infrastructure/api/dto"
	"net/http"
)

type OrdererUseCase interface {
	CreateOrder(ctx context.Context, authUserID string, reqOrder *dto.OrderCreateRequest) error
	GetOrdersByCompany(ctx context.Context, userID string, request *http.Request) ([]entities.Order, *entities.OrderQueryParams, int64, error)
	GetOrderByID(ctx context.Context, orderID string) (*entities.Order, error)
	ChangeStatus(ctx context.Context, id, status string) error
}
