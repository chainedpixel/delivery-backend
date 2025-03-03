package services

import (
	"context"
	"domain/delivery/constants"
	"domain/delivery/interfaces"
	"domain/delivery/models/orders"
	"domain/delivery/ports"
	"domain/delivery/value_objects"
	errPackage "domain/error"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"shared/logs"
	"time"
)

type OrderService struct {
	repo ports.OrderRepository
}

func NewOrderService(repo ports.OrderRepository) interfaces.Orderer {
	return &OrderService{
		repo: repo,
	}
}

func (o OrderService) GetOrderByTrackingNumber(ctx context.Context, trackingNumber string) (*orders.Order, error) {
	order, err := o.repo.GetOrderByTrackingNumber(ctx, trackingNumber)
	if err != nil {
		logs.Error("Failed to get order by tracking number", map[string]interface{}{
			"trackingNumber": trackingNumber,
			"error":          err.Error(),
		})
		return nil, errPackage.NewDomainErrorWithCause("OrderService", "GetOrderByTrackingNumber", "failed to get order by tracking number", err)
	}

	return order, nil
}

func (o OrderService) GetOrdersByClientID(ctx context.Context, clientID string) ([]orders.Order, error) {
	getOrders, err := o.repo.GetOrdersByUserID(ctx, clientID)
	if err != nil {
		logs.Error("Failed to get orders by client id", map[string]interface{}{
			"clientID": clientID,
			"error":    err.Error(),
		})
		return nil, errPackage.NewDomainErrorWithCause("OrderService", "GetOrdersByClientID", "failed to get orders by client id", err)
	}

	return getOrders, nil
}

func (o OrderService) AssignDriverToOrder(ctx context.Context, orderID, driverID string) error {
	err := o.repo.AssignDriverToOrder(ctx, orderID, driverID)
	if err != nil {
		logs.Error("Failed to assign driver to order", map[string]interface{}{
			"orderID":  orderID,
			"driverID": driverID,
			"error":    err.Error(),
		})
		return errPackage.NewDomainErrorWithCause("OrderService", "AssignDriverToOrder", "failed to assign driver to order", err)
	}

	return nil
}

func (o OrderService) CreateOrder(ctx context.Context, order *orders.Order) error {
	if order == nil {
		logs.Error("Order is nil")
		return errPackage.NewDomainError("OrderService", "CreateOrder", "order is nil")
	}

	// 1. Generar estado historico inicial
	statusHistory := &orders.StatusHistory{
		ID:      uuid.NewString(),
		OrderID: order.ID,
		Status:  constants.OrderStatusPending,
	}
	order.StatusHistory = append(order.StatusHistory, *statusHistory)

	// 2. Generar tracking number
	order.TrackingNumber = generateTrackingNumber()

	//3. Crear QR
	err := o.repo.CreateQRData(ctx, generateQRCode(*order))
	if err != nil {
		logs.Error("Failed to create qr code", map[string]interface{}{
			"orderID":        order.ID,
			"trackingNumber": order.TrackingNumber,
			"error":          err.Error(),
		})
		return errPackage.NewDomainErrorWithCause("OrderService", "CreateOrder", "failed to create qr code", err)
	}

	//4. Crear pedido
	err = o.repo.CreateOrder(ctx, order)
	if err != nil {
		logs.Error("Failed to create order", map[string]interface{}{
			"orderID":        order.ID,
			"trackingNumber": order.TrackingNumber,
			"error":          err.Error(),
		})
		return errPackage.NewDomainErrorWithCause("OrderService", "CreateOrder", "failed to create order", err)
	}

	return nil
}

func (o OrderService) ChangeStatus(ctx context.Context, id, status string) error {
	// 1. Validar que el estado sea valido
	if !value_objects.NewOrderStatus(status).IsValid() {
		logs.Error("Invalid order status", map[string]interface{}{
			"status": status,
		})
		return errPackage.NewDomainError("OrderService", "ChangeStatus", "invalid order status")
	}

	// 2. Obtener pedido para obtener estado actual
	order, err := o.repo.GetOrderByID(ctx, id)
	if err != nil {
		logs.Error("Failed to get order by id", map[string]interface{}{
			"orderID": id,
			"error":   err.Error(),
		})
		return errPackage.NewDomainErrorWithCause("OrderService", "ChangeStatus", "failed to get order by id", err)
	}

	// 3. Validar que la transicion de estados sea valida
	if !value_objects.NewOrderStatus(order.Status).CanTransitionTo(value_objects.NewOrderStatus(status)) {
		logs.Error("Invalid transition", map[string]interface{}{
			"from": order.Status,
			"to":   status,
		})
		return errPackage.NewDomainError("OrderService", "ChangeStatus", fmt.Sprintf("invalid transition from %s to %s", order.Status, status))
	}

	// 3. Cambiar estado
	err = o.repo.ChangeStatus(ctx, id, status)
	if err != nil {
		logs.Error("Failed to change status", map[string]interface{}{
			"orderID": id,
			"status":  status,
			"error":   err.Error(),
		})
		return errPackage.NewDomainErrorWithCause("OrderService", "ChangeStatus", "failed to change status", err)
	}

	return nil
}

func (o OrderService) GetOrderByID(ctx context.Context, orderID string) (*orders.Order, error) {
	order, err := o.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		logs.Error("Failed to get order by id", map[string]interface{}{
			"orderID": orderID,
			"error":   err.Error(),
		})
		return nil, errPackage.NewDomainErrorWithCause("OrderService", "GetOrder", "failed to get order by id", err)
	}

	return order, nil
}

func (o OrderService) GetOrders(ctx context.Context) ([]orders.Order, error) {
	dbOrders, err := o.repo.GetOrders(ctx)
	if err != nil {
		logs.Error("Failed to get orders", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errPackage.NewDomainErrorWithCause("OrderService", "GetOrders", "failed to get orders", err)
	}

	return dbOrders, nil
}

func (o OrderService) UpdateOrder(ctx context.Context, orderID string, order *orders.Order) error {
	err := o.repo.UpdateOrder(ctx, orderID, order)
	if err != nil {
		logs.Error("Failed to update order", map[string]interface{}{
			"orderID": orderID,
			"error":   err.Error(),
		})
		return errPackage.NewDomainErrorWithCause("OrderService", "UpdateOrder", "failed to update order", err)
	}

	return nil
}

func generateQRCode(order orders.Order) *orders.QRCode {
	return &orders.QRCode{
		OrderID: order.ID,
		QRData:  order.TrackingNumber,
	}
}

func generateTrackingNumber() string {
	// Formato: [prefijo]-[timestamp]-[aleatorio]
	prefix := "DEL"
	timestamp := time.Now().Format("060102")
	random := fmt.Sprintf("%04d", rand.Intn(10000))

	return fmt.Sprintf("%s-%s-%s", prefix, timestamp, random)
}
