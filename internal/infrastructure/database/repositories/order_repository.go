package repositories

import (
	"context"
	"domain/delivery/models/orders"
	"domain/delivery/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
	errPackage "infrastructure/error"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) ports.OrderRepository {
	return &orderRepository{
		db: db,
	}
}

// CreateOrder crea un nuevo pedido
func (r *orderRepository) CreateOrder(ctx context.Context, order *orders.Order) error {
	if order == nil {
		return errPackage.ErrNilOrder
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if order != nil {
			if err := tx.Create(order).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// GetOrderByID obtiene un pedido por ID
func (r *orderRepository) GetOrderByID(ctx context.Context, id string) (*orders.Order, error) {
	var order orders.Order
	err := r.applyOrderPreloads(r.db.WithContext(ctx)).
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &order, err
}

// GetOrderByQR obtiene un pedido por QR
func (r *orderRepository) GetOrderByQR(ctx context.Context, qr *orders.QRCode) (*orders.Order, error) {
	var order orders.Order
	err := r.applyOrderPreloads(r.db.WithContext(ctx)).
		First(&order, "id = ?", qr.OrderID).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOrderByTrackingNumber obtiene un pedido por número de seguimiento
func (r *orderRepository) GetOrderByTrackingNumber(ctx context.Context, trackingNumber string) (*orders.Order, error) {
	var order orders.Order
	err := r.applyOrderPreloads(r.db.WithContext(ctx)).
		First(&order, "tracking_number = ?", trackingNumber).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOrdersByUserID obtiene los pedidos de un usuario
func (r *orderRepository) GetOrdersByUserID(ctx context.Context, userID string) ([]orders.Order, error) {
	var orders []orders.Order
	err := r.applyOrderPreloads(r.db.WithContext(ctx)).
		Find(&orders, "client_id = ?", userID).Error

	return orders, err
}

// GetOrders obtiene todos los pedidos
func (r *orderRepository) GetOrders(ctx context.Context) ([]orders.Order, error) {
	var orders []orders.Order
	err := r.applyOrderPreloads(r.db.WithContext(ctx)).
		Find(&orders).Error

	return orders, err
}

// UpdateOrder actualiza un pedido
func (r *orderRepository) UpdateOrder(ctx context.Context, orderID string, order *orders.Order) error {
	if order == nil {
		return errPackage.ErrNilOrder
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if order != nil {
			if err := tx.Save(order).Where("id = ?", orderID).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// DeleteOrder elimina un pedido
func (r *orderRepository) DeleteOrder(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&orders.Order{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

// ChangeStatus cambia el estado de un pedido
func (r *orderRepository) ChangeStatus(ctx context.Context, id string, status string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&orders.Order{}).Where("id = ?", id).Update("status", status).Error; err != nil {
			return err
		}

		// Guardar historial de estado
		statusHistory := orders.StatusHistory{
			ID:      uuid.NewString(),
			OrderID: id,
			Status:  status,
		}

		return tx.Create(&statusHistory).Error
	})

	return err
}

func (r *orderRepository) AssignDriverToOrder(ctx context.Context, orderID, driverID string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&orders.Order{}).Where("id = ?", orderID).Update("driver_id", driverID).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *orderRepository) CreateQRData(ctx context.Context, qr *orders.QRCode) error {
	if qr == nil {
		return errPackage.ErrNilQR
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if qr != nil {
			if err := tx.Create(qr).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (r *orderRepository) applyOrderPreloads(query *gorm.DB) *gorm.DB {
	return query.
		Preload("Company").
		Preload("Branch").
		Preload("Client").
		Preload("Driver").
		Preload("Detail").
		Preload("PackageDetail").
		Preload("DeliveryAddress").
		Preload("PickupAddress").
		Preload("Tracking").
		Preload("QRCode").
		Preload("StatusHistory").
		Preload("WarehouseTrackings").
		Preload("WarehouseInventory")
}
