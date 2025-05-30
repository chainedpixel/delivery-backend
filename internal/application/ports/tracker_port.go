package ports

import (
	"context"
	"net/http"

	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/websocket"
)

// TrackerUseCase define la interfaz para el caso de uso del rastreador
type TrackerUseCase interface {
	// HandleWebSocket gestiona una nueva conexión WebSocket
	HandleWebSocket(w http.ResponseWriter, r *http.Request)

	// SendOrderUpdate envía actualizaciones del estado de un pedido
	SendOrderUpdate(ctx context.Context, orderID string, data *websocket.OrderUpdateData) error

	// SendLocationUpdate envía actualizaciones de ubicación
	SendLocationUpdate(ctx context.Context, orderID string, data *websocket.LocationUpdateData) error

	// UpdateDriverLocation actualiza la ubicación de un repartidor
	UpdateDriverLocation(ctx context.Context, orderID string, latitude, longitude float64) error
}
