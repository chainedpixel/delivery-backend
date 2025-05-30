package interfaces

import (
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/websocket"
)

// OrderTracker define las funciones necesarias para rastrear pedidos en tiempo real
type OrderTracker interface {
	// SendOrderUpdate envía una actualización del estado de un pedido a todos los clientes suscritos
	SendOrderUpdate(orderID string, data *websocket.OrderUpdateData) error

	// SendLocationUpdate envía una actualización de ubicación del repartidor a todos los clientes suscritos
	SendLocationUpdate(orderID string, data *websocket.LocationUpdateData) error
}
