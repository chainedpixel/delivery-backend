package order

import (
	"context"
	"github.com/MarlonG1/delivery-backend/internal/application/ports"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/interfaces"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/auth"
	wsModels "github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/websocket"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/websocket"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
	ws "github.com/gorilla/websocket"
	"net/http"
)

// TrackerUseCase implementa el caso de uso para el rastreador de pedidos
type TrackerUseCase struct {
	trackerService interfaces.OrderTracker
	orderService   interfaces.Orderer
	hub            *websocket.Hub
	upgrader       ws.Upgrader
}

// NewTrackerUseCase crea una nueva instancia del caso de uso
func NewTrackerUseCase(trackerService interfaces.OrderTracker, orderService interfaces.Orderer, hub *websocket.Hub) ports.TrackerUseCase {
	return &TrackerUseCase{
		trackerService: trackerService,
		orderService:   orderService,
		hub:            hub,
		upgrader: ws.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// En producción, validar el origen con más detalle
				return true
			},
		},
	}
}

// HandleWebSocket maneja una nueva conexión WebSocket
func (uc *TrackerUseCase) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 1. Verificar si el usuario está autenticado
	claims, ok := r.Context().Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", map[string]interface{}{
			"error": "Failed to get claims from context",
		})
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Actualizar la conexión HTTP a WebSocket
	conn, err := uc.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logs.Error("Failed to upgrade to websocket", map[string]interface{}{
			"error":   err.Error(),
			"user_id": claims.UserID,
		})
		return
	}

	// 3. Crear un nuevo cliente
	client := websocket.NewClient(uc.hub, conn, claims.UserID)

	// 4. Iniciar el cliente
	client.Start()

	logs.Info("New WebSocket connection established", map[string]interface{}{
		"user_id": claims.UserID,
	})
}

// SendOrderUpdate envía una actualización del estado de un pedido
func (uc *TrackerUseCase) SendOrderUpdate(ctx context.Context, orderID string, data *wsModels.OrderUpdateData) error {
	return uc.trackerService.SendOrderUpdate(orderID, data)
}

// SendLocationUpdate envía una actualización de ubicación
func (uc *TrackerUseCase) SendLocationUpdate(ctx context.Context, orderID string, data *wsModels.LocationUpdateData) error {
	return uc.trackerService.SendLocationUpdate(orderID, data)
}

// UpdateDriverLocation actualiza la ubicación del repartidor para un pedido
func (uc *TrackerUseCase) UpdateDriverLocation(ctx context.Context, orderID string, latitude, longitude float64) error {
	return uc.orderService.UpdateDriverLocation(ctx, orderID, latitude, longitude)
}
