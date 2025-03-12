package handlers

import (
	"application/ports"
	"domain/delivery/models/auth"
	"encoding/json"
	"github.com/gorilla/mux"
	"infrastructure/api/dto"
	"infrastructure/api/responser"
	"net/http"
	"shared/logs"
	"shared/mappers/response_mapper"
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
// @Success      201  {object}  string "Order created successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener los claims del contexto
	claims, ok := r.Context().Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to retrieve claims from context", nil)
		h.respWriter.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// 2. Decodificar solicitud
	var requestDTO dto.OrderCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&requestDTO); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Verificar si la solicitud es válida
	if err := requestDTO.Validate(); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 4. Llamar al caso de uso
	err := h.useCase.CreateOrder(r.Context(), claims.UserID, &requestDTO)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 5. Responder
	h.respWriter.Success(w, http.StatusCreated, "Order created successfully")
}

// UpdateOrder godoc
// @Summary      This endpoint is used to update an order by ID
// @Description  Update order by ID
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_id path string true "Order ID"
// @Param        order body dto.OrderUpdateRequest true "Order information"
// @Success      200  {object}  string "Order updated successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/orders/{order_id} [put]
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del pedido
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	// 2. Decodificar solicitud
	var requestDTO dto.OrderUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&requestDTO); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Verificar si la solicitud es válida
	if err := requestDTO.Validate(); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 4. Llamar al caso de uso
	err := h.useCase.UpdateOrder(r.Context(), orderID, &requestDTO)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 5. Responder
	h.respWriter.Success(w, http.StatusOK, "Order updated successfully")
}

// GetOrderByID godoc
// @Summary      This endpoint is used to get an order by ID
// @Description  Get order by ID
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_id path string true "Order ID"
// @Success      200  {object}  dto.OrderResponse
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/orders/{order_id} [get]
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del pedido
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	// 2. Obtener pedido
	order, err := h.useCase.GetOrderByID(r.Context(), orderID)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Mapear a DTO
	response := response_mapper.OrderToResponseDTO(order)

	// 4. Responder
	h.respWriter.Success(w, http.StatusOK, response)
}

// ChangeOrderStatus godoc
// @Summary      This endpoint is used to change the status of an order
// @Description  Change order status
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_id path string true "Order ID"
// @Param        status query string true "New status"
// @Success      200  {object}  dto.OrderResponse
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/orders/{order_id} [patch]
func (h *OrderHandler) ChangeOrderStatus(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del pedido
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	// 2. Extraer el nuevo estado del pedido
	status := r.URL.Query().Get("status")

	// 3. Cambiar estado
	err := h.useCase.ChangeStatus(r.Context(), orderID, status)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 4. Responder
	h.respWriter.Success(w, http.StatusOK, "Order status changed successfully")
}

// GetOrdersByCompany godoc
// @Summary      This endpoint is used to get orders by company
// @Description  Get orders by company
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number"
// @Param page_size query int false "Page size"
// @param sort_by query string false "Sort by"
// @param tracking_number query string false "Tracking number"
// @Param location query string false "Location like address, city, state, postal code, etc."
// @Param sort_direction query string false "Sort order"
// @Param status query string false "Order status"
// @Param start_date query string false "Start date"
// @Param end_date query string false "End date"
// @Param include_deleted query bool false "Include deleted orders"
// @Success      200  {object}  dto.PaginatedResponse
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/orders [get]
func (h *OrderHandler) GetOrdersByCompany(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer los claims del contexto
	claims, ok := r.Context().Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to retrieve claims from context", nil)
		h.respWriter.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// 2. Obtener pedidos y parametros de consulta
	orders, params, total, err := h.useCase.GetOrdersByCompany(r.Context(), claims.UserID, r)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Mapear a DTOs
	response := response_mapper.MapOrdersToResponse(orders, params, total)

	// 4. Responder
	h.respWriter.Success(w, http.StatusOK, response)
}

// DeleteOrder godoc
// @Summary      This endpoint is used to delete an order
// @Description  Delete order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_id path string true "Order ID"
// @Success      200  {object}  string "Order deleted successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/orders/{order_id} [delete]
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del pedido
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	// 2. Eliminar pedido
	err := h.useCase.DeleteOrder(r.Context(), orderID)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Responder
	h.respWriter.Success(w, http.StatusOK, "Order deleted successfully")
}

// RestoreOrder godoc
// @Summary      This endpoint is used to restore a deleted order
// @Description  Restore order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_id path string true "Order ID"
// @Success      200  {object}  string "Order restored successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/orders/recovery/{order_id} [get]
func (h *OrderHandler) RestoreOrder(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del pedido
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	// 2. Restaurar pedido
	err := h.useCase.RestoreOrder(r.Context(), orderID)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Responder
	h.respWriter.Success(w, http.StatusOK, "Order restored successfully")
}
