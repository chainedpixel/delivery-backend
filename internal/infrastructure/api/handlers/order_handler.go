package handlers

import (
	"application/ports"
	"encoding/json"
	"infrastructure/api/dto"
	"infrastructure/api/responser"
	"net/http"
)

type OrderHandler struct {
	useCase    ports.OrdererUseCase
	respWriter *responser.ResponseWriter
}

func NewOrderHandler(useCase ports.OrdererUseCase) *OrderHandler {
	return &OrderHandler{
		useCase:    useCase,
		respWriter: responser.NewResponseWriter(),
	}
}

// CreateOrder godoc
// @Summary      This endpoint is used to create a new order
// @Description  Create a new order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order body dto.OrderCreateRequest true "Order information"
// @Success      201  {object}  dto.OrderResponse
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar solicitud
	var requestDTO dto.OrderCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&requestDTO); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 2. Verificar si la solicitud es válida
	if err := requestDTO.Validate(); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Llamar al caso de uso
	err := h.useCase.CreateOrder(r.Context(), &requestDTO)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 4. Responder
	h.respWriter.Success(w, http.StatusCreated, "Order created successfully")
}
