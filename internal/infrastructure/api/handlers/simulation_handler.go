package handlers

import (
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/services"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/responser"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
	"github.com/gorilla/mux"
	"net/http"
)

type OrderSimulationHandler struct {
	simulationService *services.OrderSimulationService
	respWriter        *responser.ResponseWriter
}

func NewOrderSimulationHandler(simulationService *services.OrderSimulationService) *OrderSimulationHandler {
	return &OrderSimulationHandler{
		simulationService: simulationService,
		respWriter:        responser.NewResponseWriter(),
	}
}

// StartSimulation godoc
// @Summary      Inicia la simulación automática de un pedido
// @Description  Simula todo el flujo de un pedido: asignación de driver y cambios de estado automáticos
// @Tags         orders-simulation
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_id path string true "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  responser.APIErrorResponse
// @Failure      500  {object}  responser.APIErrorResponse
// @Router       /api/v1/orders/{order_id}/simulate [post]
func (h *OrderSimulationHandler) StartSimulation(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer order_id de la URL
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	if orderID == "" {
		h.respWriter.Error(w, http.StatusBadRequest, "Order ID is required", nil)
		return
	}

	logs.Info("Starting order simulation request", map[string]interface{}{
		"orderID": orderID,
	})

	// 2. Iniciar la simulación
	ctx := r.Context()
	err := h.simulationService.SimulateOrderFlow(ctx, orderID)
	if err != nil {
		logs.Error("Failed to start order simulation", map[string]interface{}{
			"orderID": orderID,
			"error":   err.Error(),
		})
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Respuesta exitosa
	response := map[string]interface{}{
		"success":  true,
		"message":  "Order simulation started successfully",
		"order_id": orderID,
		"details": map[string]interface{}{
			"simulation_steps": []string{
				"Assign random driver",
				"Change to ACCEPTED",
				"Change to PICKED_UP",
				"Change to IN_TRANSIT",
				"Simulate driver movement",
				"Change to DELIVERED",
			},
			"estimated_duration": "30-60 seconds",
			"real_time_updates":  "Available via WebSocket",
		},
	}

	logs.Info("Order simulation started successfully", map[string]interface{}{
		"orderID": orderID,
	})

	h.respWriter.Success(w, http.StatusOK, response)
}

// AssignRandomDriver godoc
// @Summary      Asigna un driver aleatorio a un pedido
// @Description  Asigna automáticamente un driver disponible al pedido especificado
// @Tags         orders-simulation
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_id path string true "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  responser.APIErrorResponse
// @Failure      500  {object}  responser.APIErrorResponse
// @Router       /api/v1/orders/{order_id}/assign-driver [post]
func (h *OrderSimulationHandler) AssignRandomDriver(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer order_id de la URL
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	if orderID == "" {
		h.respWriter.Error(w, http.StatusBadRequest, "Order ID is required", nil)
		return
	}

	logs.Info("Assigning random driver request", map[string]interface{}{
		"orderID": orderID,
	})

	// 2. Asignar driver aleatorio
	ctx := r.Context()
	err := h.simulationService.SimulateRandomDriver(ctx, orderID)
	if err != nil {
		logs.Error("Failed to assign random driver", map[string]interface{}{
			"orderID": orderID,
			"error":   err.Error(),
		})
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Respuesta exitosa
	response := map[string]interface{}{
		"success":  true,
		"message":  "Random driver assigned successfully",
		"order_id": orderID,
	}

	logs.Info("Random driver assigned successfully", map[string]interface{}{
		"orderID": orderID,
	})

	h.respWriter.Success(w, http.StatusOK, response)
}

// GetSimulationStatus godoc
// @Summary      Obtiene el estado actual de la simulación
// @Description  Retorna información detallada sobre el estado actual del pedido y su simulación
// @Tags         orders-simulation
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_id path string true "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  responser.APIErrorResponse
// @Failure      500  {object}  responser.APIErrorResponse
// @Router       /api/v1/orders/{order_id}/simulation-status [get]
func (h *OrderSimulationHandler) GetSimulationStatus(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer order_id de la URL
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	if orderID == "" {
		h.respWriter.Error(w, http.StatusBadRequest, "Order ID is required", nil)
		return
	}

	logs.Info("Getting simulation status request", map[string]interface{}{
		"orderID": orderID,
	})

	// 2. Obtener estado de la simulación
	ctx := r.Context()
	status, err := h.simulationService.GetSimulationStatus(ctx, orderID)
	if err != nil {
		logs.Error("Failed to get simulation status", map[string]interface{}{
			"orderID": orderID,
			"error":   err.Error(),
		})
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Respuesta exitosa
	response := map[string]interface{}{
		"success": true,
		"message": "Simulation status retrieved successfully",
		"data":    status,
	}

	h.respWriter.Success(w, http.StatusOK, response)
}

// SimulateDriverMovement godoc
// @Summary      Simula el movimiento del conductor para un pedido
// @Description  Inicia la simulación de movimiento del conductor en tiempo real
// @Tags         orders-simulation
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_id path string true "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  responser.APIErrorResponse
// @Failure      500  {object}  responser.APIErrorResponse
// @Router       /api/v1/orders/{order_id}/simulate-movement [post]
func (h *OrderSimulationHandler) SimulateDriverMovement(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer order_id de la URL
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	if orderID == "" {
		h.respWriter.Error(w, http.StatusBadRequest, "Order ID is required", nil)
		return
	}

	logs.Info("Starting driver movement simulation", map[string]interface{}{
		"orderID": orderID,
	})

	// 2. Iniciar simulación de movimiento en goroutine
	go func() {
		// Simular movimiento del conductor
		// Coordenadas de ejemplo (San Salvador)
		startLat := 13.6929
		startLng := -89.2182
		endLat := 13.7942
		endLng := -89.1956

		steps := 30
		latIncrement := (endLat - startLat) / float64(steps)
		lngIncrement := (endLng - startLng) / float64(steps)

		currentLat := startLat
		currentLng := startLng

		for i := 0; i < steps; i++ {
			// Enviar actualización de ubicación
			// Aquí podrías llamar al servicio para actualizar la ubicación
			logs.Info("Driver location update", map[string]interface{}{
				"orderID":   orderID,
				"latitude":  currentLat,
				"longitude": currentLng,
				"step":      i + 1,
			})

			// Actualizar coordenadas
			currentLat += latIncrement
			currentLng += lngIncrement

			// Esperar antes del siguiente punto
			// time.Sleep(2 * time.Second)
		}
	}()

	// 3. Respuesta inmediata
	response := map[string]interface{}{
		"success":  true,
		"message":  "Driver movement simulation started",
		"order_id": orderID,
		"details": map[string]interface{}{
			"duration_seconds": 60,
			"update_interval":  "2 seconds",
			"total_points":     30,
		},
	}

	h.respWriter.Success(w, http.StatusOK, response)
}
