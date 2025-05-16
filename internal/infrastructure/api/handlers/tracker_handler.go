package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MarlonG1/delivery-backend/internal/application/ports"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/responser"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
	"github.com/gorilla/mux"

	_ "github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
)

// TrackerHandler maneja las peticiones relacionadas con el rastreo de pedidos
type TrackerHandler struct {
	useCase    ports.TrackerUseCase
	respWriter *responser.ResponseWriter
}

// NewTrackerHandler crea una nueva instancia del manejador de rastreo
func NewTrackerHandler(useCase ports.TrackerUseCase) *TrackerHandler {
	return &TrackerHandler{
		useCase:    useCase,
		respWriter: responser.NewResponseWriter(),
	}
}

// HandleWebSocket godoc
// @Summary      Establece una conexión WebSocket para rastrear pedidos en tiempo real
// @Description  Crea un canal WebSocket para recibir actualizaciones de pedidos
// @Tags         tracking
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      101  {string}  string "Conexión WebSocket establecida"
// @Failure      401  {object}  responser.APIErrorResponse
// @Router       /api/v1/tracking/ws [get]
func (h *TrackerHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	h.useCase.HandleWebSocket(w, r)
}

// UpdateDriverLocation godoc
// @Summary      Actualiza la ubicación del repartidor
// @Description  Actualiza la ubicación del repartidor para un pedido específico
// @Tags         tracking
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_id path string true "ID del pedido"
// @Param        location body LocationUpdateRequest true "Datos de ubicación"
// @Success      200  {string}  string "Ubicación actualizada correctamente"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/tracking/location/{order_id} [post]
func (h *TrackerHandler) UpdateDriverLocation(w http.ResponseWriter, r *http.Request) {

	// Extraer el ID del pedido
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	// Decodificar solicitud
	var req LocationUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Validar los datos
	if !req.isValid() {
		h.respWriter.Error(w, http.StatusBadRequest, "Datos de ubicación inválidos", nil)
		return
	}

	// Actualizar la ubicación
	err := h.useCase.UpdateDriverLocation(r.Context(), orderID, req.Latitude, req.Longitude)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	logs.Info("Driver location updated successfully", map[string]interface{}{
		"order_id":  orderID,
		"latitude":  req.Latitude,
		"longitude": req.Longitude,
	})
	h.respWriter.Success(w, http.StatusOK, "Ubicación actualizada correctamente")
}

// LocationUpdateRequest representa la solicitud para actualizar ubicación
type LocationUpdateRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// isValid verifica si los datos de ubicación son válidos
func (r *LocationUpdateRequest) isValid() bool {
	// Validación básica de coordenadas
	return r.Latitude >= -90 && r.Latitude <= 90 && r.Longitude >= -180 && r.Longitude <= 180
}
