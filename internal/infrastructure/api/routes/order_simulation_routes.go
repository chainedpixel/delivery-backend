package routes

import (
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterSimulationRoutes(router *mux.Router, simulationHandler *handlers.OrderSimulationHandler) {
	router.HandleFunc("/orders/{order_id}/simulate", simulationHandler.StartSimulation).Methods(http.MethodPost)
	router.HandleFunc("/orders/{order_id}/assign-driver", simulationHandler.AssignRandomDriver).Methods(http.MethodPost)
	router.HandleFunc("/orders/{order_id}/simulation-status", simulationHandler.GetSimulationStatus).Methods(http.MethodGet)
	router.HandleFunc("/orders/{order_id}/simulate-movement", simulationHandler.SimulateDriverMovement).Methods(http.MethodPost)
}
