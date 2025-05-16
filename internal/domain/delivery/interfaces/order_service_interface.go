package interfaces

import (
	"context"

	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
)

type Orderer interface {
	CreateOrder(ctx context.Context, order *entities.Order) error
	ChangeStatus(ctx context.Context, id, status string) error
	GetOrderByID(ctx context.Context, orderID string) (*entities.Order, error)
	GetOrders(ctx context.Context) ([]entities.Order, error)
	GetOrdersByCompany(ctx context.Context, companyID string, params *entities.OrderQueryParams) ([]entities.Order, int64, error)
	UpdateOrder(ctx context.Context, orderID string, order *entities.Order) error
	GetOrderByTrackingNumber(ctx context.Context, trackingNumber string) (*entities.Order, error)
	GetOrdersByClientID(ctx context.Context, clientID string) ([]entities.Order, error)
	AssignDriverToOrder(ctx context.Context, orderID, driverID string) error
	SoftDeleteOrder(ctx context.Context, id string) error
	OrderIsDeleted(ctx context.Context, orderID string) bool
	RestoreOrder(ctx context.Context, id string) error
	IsAvailableForDelete(ctx context.Context, orderID string) error
	UpdateDriverLocation(ctx context.Context, orderID string, latitude, longitude float64) error
}
