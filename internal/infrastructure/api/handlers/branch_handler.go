package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MarlonG1/delivery-backend/internal/application/ports"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/auth"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/dto"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/responser"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
	"github.com/MarlonG1/delivery-backend/pkg/shared/mappers/request_mapper"
	"github.com/MarlonG1/delivery-backend/pkg/shared/mappers/response_mapper"
	"github.com/gorilla/mux"
)

type BranchHandler struct {
	useCase    ports.BranchUseCase
	respWriter *responser.ResponseWriter
}

func NewBranchHandler(useCase ports.BranchUseCase) *BranchHandler {
	return &BranchHandler{
		useCase:    useCase,
		respWriter: responser.NewResponseWriter(),
	}
}

// GetBranches godoc
// @Summary      Obtiene todas las sucursales de la empresa del usuario autenticado
// @Description  Obtiene todas las sucursales con filtros y paginación
// @Tags         branches
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Número de página"
// @Param        page_size query int false "Tamaño de página"
// @Param        name query string false "Nombre de la sucursal"
// @Param        code query string false "Código de la sucursal"
// @Param        contact_name query string false "Nombre del contacto"
// @Param        contact_email query string false "Email del contacto"
// @Param        zone_id query string false "ID de la zona"
// @Param        is_active query string false "Estado (activo/inactivo)"
// @Param        sort_by query string false "Campo por el cual ordenar"
// @Param        sort_direction query string false "Dirección de ordenamiento (asc/desc)"
// @Success      200  {object}  dto.PaginatedResponse
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/branches [get]
func (h *BranchHandler) GetBranches(w http.ResponseWriter, r *http.Request) {
	// Obtener las sucursales
	branches, params, total, err := h.useCase.GetBranches(r.Context(), r)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Convertir a DTO
	response := response_mapper.MapBranchesToResponse(branches, params, total)
	h.respWriter.Success(w, http.StatusOK, response)
}

// GetBranchByID godoc
// @Summary      Obtiene una sucursal por su ID
// @Description  Obtiene los detalles de una sucursal específica
// @Tags         branches
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        branch_id path string true "ID de la sucursal"
// @Success      200  {object}  dto.BranchResponse
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/branches/{branch_id} [get]
func (h *BranchHandler) GetBranchByID(w http.ResponseWriter, r *http.Request) {
	// Extraer ID de la sucursal
	vars := mux.Vars(r)
	branchID := vars["branch_id"]

	// Obtener la sucursal
	branch, err := h.useCase.GetBranchByID(r.Context(), branchID)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Convertir a DTO
	response := response_mapper.BranchToResponseDTO(branch, true)
	h.respWriter.Success(w, http.StatusOK, response)
}

// CreateBranch godoc
// @Summary      Crea una nueva sucursal
// @Description  Crea una nueva sucursal para la empresa del usuario autenticado
// @Tags         branches
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        branch body dto.BranchCreateRequest true "Información de la sucursal"
// @Success      201  {string}  string "Sucursal creada exitosamente"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/branches [post]
func (h *BranchHandler) CreateBranch(w http.ResponseWriter, r *http.Request) {
	var req dto.BranchCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Obtener companyID del usuario autenticado
	claims, ok := r.Context().Value("claims").(*auth.AuthClaims)
	if !ok {
		logs.Error("Failed to get claims from context", nil)
		h.respWriter.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Mapear a entidad
	branch, err := request_mapper.BranchRequestToBranch(claims.CompanyID, &req)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Crear la sucursal
	err = h.useCase.CreateBranch(r.Context(), claims.CompanyID, branch)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusCreated, "Sucursal creada exitosamente")
}

// UpdateBranch godoc
// @Summary      Actualiza una sucursal existente
// @Description  Actualiza la información de una sucursal
// @Tags         branches
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        branch_id path string true "ID de la sucursal"
// @Param        branch body dto.BranchUpdateRequest true "Información actualizada de la sucursal"
// @Success      200  {string}  string "Sucursal actualizada exitosamente"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/branches/{branch_id} [put]
func (h *BranchHandler) UpdateBranch(w http.ResponseWriter, r *http.Request) {
	// Extraer ID de la sucursal
	vars := mux.Vars(r)
	branchID := vars["branch_id"]

	var req dto.BranchUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Mapear a entidad
	branch, err := request_mapper.BranchUpdateRequestToBranch(branchID, &req)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Actualizar la sucursal
	err = h.useCase.UpdateBranch(r.Context(), branchID, branch)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "Sucursal actualizada exitosamente")
}

// DeactivateBranch godoc
// @Summary      Desactiva una sucursal
// @Description  Cambia el estado de una sucursal a inactivo
// @Tags         branches
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        branch_id path string true "ID de la sucursal"
// @Success      200  {string}  string "Sucursal desactivada exitosamente"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/branches/deactivate/{branch_id} [post]
func (h *BranchHandler) DeactivateBranch(w http.ResponseWriter, r *http.Request) {
	// Extraer ID de la sucursal
	vars := mux.Vars(r)
	branchID := vars["branch_id"]

	// Desactivar la sucursal
	err := h.useCase.DeactivateBranch(r.Context(), branchID)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "Sucursal desactivada exitosamente")
}

// ReactivateBranch godoc

// ReactivateBranch godoc
// @Summary	  Reactiva una sucursal
// @Description  Cambia el estado de una sucursal a activo
// @Tags		 branches
// @Accept	   json
// @Produce	  json
// @Security	 BearerAuth
// @Param		branch_id path string true "ID de la sucursal"
// @Success	  200  {string}  string "Sucursal reactivada exitosamente"
// @Failure	  400  {object}  responser.APIErrorResponse
// @Router	   /api/v1/branches/reactivate/{branch_id} [post]
func (h *BranchHandler) ReactivateBranch(w http.ResponseWriter, r *http.Request) {
	// Extraer ID de la sucursal
	vars := mux.Vars(r)
	branchID := vars["branch_id"]

	// Reactivar la sucursal
	err := h.useCase.ReactivateBranch(r.Context(), branchID)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "Sucursal reactivada exitosamente")
}

// AssignZoneToBranch godoc
// @Summary      Asigna una zona a una sucursal
// @Description  Asigna una zona geográfica específica a una sucursal
// @Tags         branches, zones
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        branch_id path string true "ID de la sucursal"
// @Param        request body dto.ZoneAssignmentRequest true "Información de la zona a asignar"
// @Success      200  {string}  string "Zona asignada exitosamente"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/branches/zones/{branch_id} [post]
func (h *BranchHandler) AssignZoneToBranch(w http.ResponseWriter, r *http.Request) {
	// Extraer ID de la sucursal
	vars := mux.Vars(r)
	branchID := vars["branch_id"]

	// Decodificar la solicitud
	var req dto.ZoneAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Asignar la zona
	err := h.useCase.AssignZoneToBranch(r.Context(), branchID, req.ZoneID)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "Zona asignada exitosamente")
}

// GetAvailableZonesForBranch godoc
// @Summary      Obtiene las zonas disponibles para una sucursal
// @Description  Retorna las zonas que pueden ser asignadas a una sucursal específica
// @Tags         branches, zones
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        branch_id path string true "ID de la sucursal"
// @Success      200  {array}  entities.Zone
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/branches/available-zones/{branch_id} [get]
func (h *BranchHandler) GetAvailableZonesForBranch(w http.ResponseWriter, r *http.Request) {
	// Extraer ID de la sucursal
	vars := mux.Vars(r)
	branchID := vars["branch_id"]

	// Obtener las zonas disponibles
	zones, err := h.useCase.GetAvailableZonesForBranch(r.Context(), branchID)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, zones)
}

// GetBranchMetrics godoc
// @Summary      Obtiene las métricas de una sucursal
// @Description  Retorna métricas como órdenes totales, completadas, canceladas, ingresos, etc.
// @Tags         branches, metrics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        branch_id path string true "ID de la sucursal"
// @Success      200  {object}  dto.BranchMetricsResponse
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/branches/metrics/{branch_id} [get]
func (h *BranchHandler) GetBranchMetrics(w http.ResponseWriter, r *http.Request) {
	// Extraer ID de la sucursal
	vars := mux.Vars(r)
	branchID := vars["branch_id"]

	// Obtener las métricas
	metrics, err := h.useCase.GetBranchMetrics(r.Context(), branchID)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Convertir a DTO
	response := response_mapper.BranchToMetricsDTO(metrics)
	h.respWriter.Success(w, http.StatusOK, response)
}
