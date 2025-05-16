package ports

import (
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/websocket"
)

// TrackerRepository define las operaciones necesarias para el rastreo de pedidos
type TrackerRepository interface {
	// SendOrderUpdate envía actualizaciones del estado de un pedido a los clientes suscritos
	SendOrderUpdate(orderID string, data *websocket.OrderUpdateData) error

	// SendLocationUpdate envía actualizaciones de ubicación de un repartidor a los clientes suscritos
	SendLocationUpdate(orderID string, data *websocket.LocationUpdateData) error
}
