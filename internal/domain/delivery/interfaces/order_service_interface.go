package interfaces

import (
	"context"
	"domain/delivery/models/entities"
)

type Orderer interface {
	CreateOrder(ctx context.Context, order *entities.Order) error
	ChangeStatus(ctx context.Context, id, status string) error
	GetOrderByID(ctx context.Context, orderID string) (*entities.Order, error)
	GetOrders(ctx context.Context) ([]entities.Order, error)
	UpdateOrder(ctx context.Context, orderID string, order *entities.Order) error
	GetOrderByTrackingNumber(ctx context.Context, trackingNumber string) (*entities.Order, error)
	GetOrdersByClientID(ctx context.Context, clientID string) ([]entities.Order, error)
	AssignDriverToOrder(ctx context.Context, orderID, driverID string) error
}
