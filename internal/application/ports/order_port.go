package ports

import (
	"context"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/dto"
	"net/http"
)

type OrdererUseCase interface {
	CreateOrder(ctx context.Context, authUserID string, reqOrder *dto.OrderCreateRequest) error
	UpdateOrder(ctx context.Context, orderID string, reqOrder *dto.OrderUpdateRequest) error
	GetOrdersByCompany(ctx context.Context, userID string, request *http.Request) ([]entities.Order, *entities.OrderQueryParams, int64, error)
	GetOrderByID(ctx context.Context, orderID string) (*entities.Order, error)
	ChangeStatus(ctx context.Context, id, status string) error
	DeleteOrder(ctx context.Context, id string) error
	RestoreOrder(ctx context.Context, id string) error
}
