package repositories

import (
	wsModels "github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/websocket"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/ports"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/websocket"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
)

// TrackerRepository implementa el puerto TrackerRepository utilizando el hub de WebSocket
type TrackerRepository struct {
	hub *websocket.Hub
}

// NewTrackerRepository crea una nueva instancia del repositorio de rastreo
func NewTrackerRepository(hub *websocket.Hub) ports.TrackerRepository {
	return &TrackerRepository{
		hub: hub,
	}
}

// SendOrderUpdate envía una actualización de estado de pedido a los clientes suscritos
func (r *TrackerRepository) SendOrderUpdate(orderID string, data *wsModels.OrderUpdateData) error {
	logs.Info("Sending order update through tracker repository", map[string]interface{}{
		"order_id": orderID,
		"status":   data.Status,
	})

	r.hub.SendOrderUpdate(orderID, data)
	return nil
}

// SendLocationUpdate envía una actualización de ubicación a los clientes suscritos
func (r *TrackerRepository) SendLocationUpdate(orderID string, data *wsModels.LocationUpdateData) error {
	logs.Info("Sending location update through tracker repository", map[string]interface{}{
		"order_id":  orderID,
		"latitude":  data.Latitude,
		"longitude": data.Longitude,
	})

	r.hub.SendLocationUpdate(orderID, data)
	return nil
}
