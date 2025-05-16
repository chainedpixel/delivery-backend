package routes

import (
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/handlers"
	"github.com/gorilla/mux"
)

// RegisterTrackerRoutes registra las rutas relacionadas con el tracking de pedidos
func RegisterTrackerRoutes(router *mux.Router, handler *handlers.TrackerHandler) {
	router.HandleFunc("/tracking/ws", handler.HandleWebSocket).Methods("GET")
	router.HandleFunc("/tracking/location/{order_id}", handler.UpdateDriverLocation).Methods("POST")
}
