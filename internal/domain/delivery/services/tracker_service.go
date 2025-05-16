package services

import (
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/interfaces"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/websocket"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/ports"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
)

// TrackerService implementa la interfaz OrderTracker para rastrear pedidos en tiempo real
type TrackerService struct {
	trackerRepo ports.TrackerRepository
}

// NewTrackerService crea una nueva instancia del servicio de rastreo
func NewTrackerService(trackerRepo ports.TrackerRepository) interfaces.OrderTracker {
	return &TrackerService{
		trackerRepo: trackerRepo,
	}
}

// SendOrderUpdate envía una actualización del estado de un pedido
func (s *TrackerService) SendOrderUpdate(orderID string, data *websocket.OrderUpdateData) error {
	logs.Info("Sending order update through tracker service", map[string]interface{}{
		"order_id": orderID,
		"status":   data.Status,
	})
	return s.trackerRepo.SendOrderUpdate(orderID, data)
}

// SendLocationUpdate envía una actualización de la ubicación del repartidor
func (s *TrackerService) SendLocationUpdate(orderID string, data *websocket.LocationUpdateData) error {
	logs.Info("Sending location update through tracker service", map[string]interface{}{
		"order_id":  orderID,
		"latitude":  data.Latitude,
		"longitude": data.Longitude,
	})
	return s.trackerRepo.SendLocationUpdate(orderID, data)
}
