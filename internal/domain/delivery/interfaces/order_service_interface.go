package interfaces

import (
	"context"
	"domain/delivery/models/orders"
)

type Orderer interface {
	CreateOrder(ctx context.Context, order *orders.Order) error
	ChangeStatus(ctx context.Context, id, status string) error
	GetOrderByID(ctx context.Context, orderID string) (*orders.Order, error)
	GetOrders(ctx context.Context) ([]orders.Order, error)
	UpdateOrder(ctx context.Context, orderID string, order *orders.Order) error
	GetOrderByTrackingNumber(ctx context.Context, trackingNumber string) (*orders.Order, error)
	GetOrdersByClientID(ctx context.Context, clientID string) ([]orders.Order, error)
	AssignDriverToOrder(ctx context.Context, orderID, driverID string) error
}
