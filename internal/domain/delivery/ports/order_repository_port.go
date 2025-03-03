package ports

import (
	"context"
	"domain/delivery/models/orders"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *orders.Order) error
	CreateQRData(ctx context.Context, qr *orders.QRCode) error
	GetOrderByID(ctx context.Context, id string) (*orders.Order, error)
	GetOrderByQR(ctx context.Context, qr *orders.QRCode) (*orders.Order, error)
	GetOrderByTrackingNumber(ctx context.Context, trackingNumber string) (*orders.Order, error)
	GetOrdersByUserID(ctx context.Context, userID string) ([]orders.Order, error)
	GetOrders(ctx context.Context) ([]orders.Order, error)
	UpdateOrder(ctx context.Context, orderID string, order *orders.Order) error
	DeleteOrder(ctx context.Context, id string) error
	ChangeStatus(ctx context.Context, id string, status string) error
	AssignDriverToOrder(ctx context.Context, orderID, driverID string) error
}
