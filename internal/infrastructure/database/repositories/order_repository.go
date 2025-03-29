package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/constants"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/ports"
	errPackage "github.com/MarlonG1/delivery-backend/internal/infrastructure/error"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) ports.OrdererRepository {
	return &orderRepository{
		db: db,
	}
}

// CreateOrder crea un nuevo pedido
func (r *orderRepository) CreateOrder(ctx context.Context, order *entities.Order) error {
	if order == nil {
		return errPackage.ErrNilOrder
	}

	jsonMess, _ := json.Marshal(order)
	logs.Info("OrdererRepository.CreateOrder", map[string]interface{}{
		"order": string(jsonMess),
	})

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Crear la orden y sus entidades relacionadas normalmente
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 2. Actualizar los campos espaciales directamente usando gorm.Expr
		// Para dirección de entrega
		if order.DeliveryAddress != nil {
			if err := tx.Model(&entities.DeliveryAddress{}).
				Where("order_id = ?", order.ID).
				Update("location", gorm.Expr("ST_PointFromText(?)", "POINT(-90.5091 14.6234)")).
				Error; err != nil {
				return err
			}
		}

		// Para dirección de recogida
		if order.PickupAddress != nil {
			if err := tx.Model(&entities.PickupAddress{}).
				Where("order_id = ?", order.ID).
				Update("location", gorm.Expr("ST_PointFromText(?)", "POINT(-90.5191 14.6334)")).
				Error; err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (r *orderRepository) GetOrdersByCompany(ctx context.Context, companyID string, params *entities.OrderQueryParams) ([]entities.Order, int64, error) {
	var orders []entities.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Order{}).Where("company_id = ?", companyID)
	if params != nil {
		if params.Status != "" {
			query = query.Where("status = ?", params.Status)
		}

		if params.Location != "" {
			query = query.Joins("LEFT JOIN delivery_addresses ON orders.id = delivery_addresses.order_id")
			searchPattern := "%" + params.Location + "%"
			locationSearch := r.db.Where(
				"delivery_addresses.address_line1 LIKE ? OR delivery_addresses.address_line2 LIKE ? OR "+
					"delivery_addresses.city LIKE ? OR delivery_addresses.state LIKE ? OR "+
					"delivery_addresses.postal_code LIKE ? OR delivery_addresses.recipient_name LIKE ? OR "+
					"delivery_addresses.address_notes LIKE ?",
				searchPattern, searchPattern, searchPattern, searchPattern,
				searchPattern, searchPattern, searchPattern,
			)

			query = query.Where(locationSearch)
		}

		if params.IncludeDeleted {
			query = query.Unscoped()
		} else {
			query = query.Where("deleted_at IS NULL")
		}

		if params.TrackingNumber != "" {
			query = query.Where("tracking_number = ?", params.TrackingNumber)
		}

		if params.StartDate != nil && params.EndDate != nil {
			query = query.Where("created_at BETWEEN ? AND ?", params.StartDate, params.EndDate)
		} else if params.StartDate != nil {
			query = query.Where("created_at >= ?", params.StartDate)
		} else if params.EndDate != nil {
			query = query.Where("created_at <= ?", params.EndDate)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params != nil {
		if params.Page > 0 && params.PageSize > 0 {
			offset := (params.Page - 1) * params.PageSize
			query = query.Offset(offset).Limit(params.PageSize)
		}

		if params.SortBy != "" {
			direction := "ASC"
			if params.SortDirection == "desc" {
				direction = "DESC"
			}
			query = query.Order(params.SortBy + " " + direction)
		} else {
			query = query.Order("created_at DESC")
		}
	} else {
		query = query.Order("created_at DESC")
	}

	err := r.applyOrderPreloads(query).Find(&orders).Error
	return orders, total, err
}

// GetOrderByID obtiene un pedido por ID
func (r *orderRepository) GetOrderByID(ctx context.Context, id string) (*entities.Order, error) {
	var order entities.Order
	err := r.applyOrderPreloads(r.db.WithContext(ctx)).
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	order.DeliveryAddress.Latitude, order.DeliveryAddress.Longitude, err = r.GetLocationCoordinates(ctx, order.ID, "delivery")
	if err != nil {
		return nil, err
	}
	order.PickupAddress.Latitude, order.PickupAddress.Longitude, err = r.GetLocationCoordinates(ctx, order.ID, "pickup")
	if err != nil {
		return nil, err
	}

	return &order, err
}

// GetOrderByQR obtiene un pedido por QR
func (r *orderRepository) GetOrderByQR(ctx context.Context, qr *entities.QRCode) (*entities.Order, error) {
	var order entities.Order
	err := r.applyOrderPreloads(r.db.WithContext(ctx)).
		First(&order, "id = ?", qr.OrderID).Error
	if err != nil {
		return nil, err
	}

	order.DeliveryAddress.Latitude, order.DeliveryAddress.Longitude, err = r.GetLocationCoordinates(ctx, order.ID, "delivery")
	if err != nil {
		return nil, err
	}
	order.PickupAddress.Latitude, order.PickupAddress.Longitude, err = r.GetLocationCoordinates(ctx, order.ID, "pickup")
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOrderByTrackingNumber obtiene un pedido por número de seguimiento
func (r *orderRepository) GetOrderByTrackingNumber(ctx context.Context, trackingNumber string) (*entities.Order, error) {
	var order entities.Order
	err := r.applyOrderPreloads(r.db.WithContext(ctx)).
		First(&order, "tracking_number = ?", trackingNumber).Error
	if err != nil {
		return nil, err
	}

	order.DeliveryAddress.Latitude, order.DeliveryAddress.Longitude, err = r.GetLocationCoordinates(ctx, order.ID, "delivery")
	if err != nil {
		return nil, err
	}
	order.PickupAddress.Latitude, order.PickupAddress.Longitude, err = r.GetLocationCoordinates(ctx, order.ID, "pickup")
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOrdersByUserID obtiene los pedidos de un usuario
func (r *orderRepository) GetOrdersByUserID(ctx context.Context, userID string) ([]entities.Order, error) {
	var orders []entities.Order
	err := r.applyOrderPreloads(r.db.WithContext(ctx)).
		Find(&orders, "client_id = ?", userID).Error

	return orders, err
}

// GetOrders obtiene todos los pedidos
func (r *orderRepository) GetOrders(ctx context.Context) ([]entities.Order, error) {
	var orders []entities.Order
	err := r.applyOrderPreloads(r.db.WithContext(ctx)).
		Find(&orders).Error

	return orders, err
}

// UpdateOrder actualiza un pedido
func (r *orderRepository) UpdateOrder(ctx context.Context, orderID string, order *entities.Order) error {
	if order == nil {
		return errPackage.ErrNilOrder
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Actualizar la tabla principal orders
		if err := tx.Model(&entities.Order{}).Where("id = ?", orderID).Updates(map[string]interface{}{
			"updated_at": order.UpdatedAt,
		}).Error; err != nil {
			return err
		}

		// 2. Actualizar detalles del pedido si existen
		if order.Detail != nil {
			if err := tx.Model(&entities.Details{}).Where("order_id = ?", orderID).Updates(order.Detail).Error; err != nil {
				return err
			}
		}

		// 3. Actualizar detalles del paquete si existen
		if order.PackageDetail != nil {
			if err := tx.Model(&entities.PackageDetail{}).Where("order_id = ?", orderID).Updates(order.PackageDetail).Error; err != nil {
				return err
			}
		}

		// 4. Actualizar dirección de entrega si existe
		if order.DeliveryAddress != nil {
			if err := tx.Model(&entities.DeliveryAddress{}).Where("order_id = ?", orderID).Updates(order.DeliveryAddress).Error; err != nil {
				return err
			}
		}

		// 5. Actualizar dirección de recogida si existe
		if order.PickupAddress != nil {
			if err := tx.Model(&entities.PickupAddress{}).Where("order_id = ?", orderID).Updates(order.PickupAddress).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteOrder elimina un pedido
func (r *orderRepository) DeleteOrder(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&entities.Order{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

// ChangeStatus cambia el estado de un pedido
func (r *orderRepository) ChangeStatus(ctx context.Context, id string, status string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entities.Order{}).Where("id = ?", id).Update("status", status).Error; err != nil {
			return err
		}

		// Guardar historial de estado
		statusHistory := entities.StatusHistory{
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
		if err := tx.Model(&entities.Order{}).Where("id = ?", orderID).Update("driver_id", driverID).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *orderRepository) CreateQRData(ctx context.Context, qr *entities.QRCode) error {
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

// GetLocationCoordinates obtiene las coordenadas de una dirección
func (r *orderRepository) GetLocationCoordinates(ctx context.Context, orderID string, addressType string) (float64, float64, error) {
	var tableName string
	if addressType == "delivery" {
		tableName = "delivery_addresses"
	} else {
		tableName = "pickup_addresses"
	}

	var result struct {
		Lat float64
		Lng float64
	}

	err := r.db.WithContext(ctx).Raw(
		fmt.Sprintf("SELECT ST_Y(location) as lat, ST_X(location) as lng FROM %s WHERE order_id = ?", tableName),
		orderID,
	).Scan(&result).Error

	return result.Lat, result.Lng, err
}

// SoftDeleteOrder realiza una eliminación lógica del pedido
func (r *orderRepository) SoftDeleteOrder(ctx context.Context, id string) error {
	now := time.Now()

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Marcar el pedido como eliminado
		if err := tx.Model(&entities.Order{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"deleted_at": now,
				"status":     constants.OrderStatusDeleted,
			}).Error; err != nil {
			return err
		}

		// 2. Crear un registro en el historial de estados
		statusHistory := entities.StatusHistory{
			ID:          uuid.NewString(),
			OrderID:     id,
			Status:      constants.OrderStatusDeleted,
			Description: "Pedido eliminado por el usuario",
			CreatedAt:   now,
		}

		if err := tx.Create(&statusHistory).Error; err != nil {
			return err
		}

		return nil
	})
}

// RestoreOrder restaura un pedido previamente eliminado lógicamente
func (r *orderRepository) RestoreOrder(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Restaurar el pedido
		if err := tx.Model(&entities.Order{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"deleted_at": nil,
				"status":     constants.OrderStatusRestored,
			}).Error; err != nil {
			return err
		}

		// 2. Crear un registro en el historial de estados
		statusHistory := entities.StatusHistory{
			ID:          uuid.NewString(),
			OrderID:     id,
			Status:      constants.OrderStatusRestored,
			Description: "Pedido restaurado",
			CreatedAt:   time.Now(),
		}

		return tx.Create(&statusHistory).Error
	})
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
