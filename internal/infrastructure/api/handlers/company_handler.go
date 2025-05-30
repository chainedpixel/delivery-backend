package handlers

import (
	"encoding/json"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
	"net/http"

	"github.com/MarlonG1/delivery-backend/internal/application/ports"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/dto"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/responser"
	"github.com/MarlonG1/delivery-backend/pkg/shared/mappers/request_mapper"
	"github.com/MarlonG1/delivery-backend/pkg/shared/mappers/response_mapper"
	"github.com/gorilla/mux"

	_ "github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
)

type CompanyHandler struct {
	useCase    ports.CompanyUseCase
	respWriter *responser.ResponseWriter
}

func NewCompanyHandler(useCase ports.CompanyUseCase) *CompanyHandler {
	return &CompanyHandler{
		useCase:    useCase,
		respWriter: responser.NewResponseWriter(),
	}
}

// GetCompanyProfile godoc
// @Summary      Obtener la información de principal de la compañia
// @Description   Obtener la información de principal de la compañia
// @Tags         companies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.CompanyResponse
// @Failure      401  {object}  responser.APIErrorResponse
// @Router       /api/v1/companies/profile [get]
func (h *CompanyHandler) GetCompanyProfile(w http.ResponseWriter, r *http.Request) {
	// Obtener la información de la compañía
	company, err := h.useCase.GetCompanyByID(r.Context())
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Obtener las métricas de la compañía
	metrics, err := h.useCase.GetCompanyMetrics(r.Context())
	if err != nil {
		// Solo log del error de métricas, no detener la respuesta
		logs.Warn("Failed to get company metrics", map[string]interface{}{
			"error":      err.Error(),
			"company_id": company.ID,
		})
		// Continuamos sin métricas
	}

	// Convertir a DTO usando el nuevo mapper que incluye métricas
	response := response_mapper.CompanyToResponseWithMetricsDTO(company, metrics, true)

	h.respWriter.Success(w, http.StatusOK, response)
}

// CreateCompany godoc
// @Summary      Crear una nueva compañia
// @Description  Crear una nueva compañia
// @Tags         companies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        company body dto.CompanyCreateRequest true "Company information"
// @Success      201  {string}  string "Company created successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/companies [post]
func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var req dto.CompanyCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Mapear a entidad
	company, err := request_mapper.CompanyRequestToCompany(&req)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Crear la empresa
	err = h.useCase.CreateCompany(r.Context(), company)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusCreated, "Company created successfully")
}

// UpdateCompany godoc
// @Summary      Actualizar la información de la compañia
// @Description  Actualizar la información de la compañia
// @Tags         companies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        company body dto.CompanyUpdateRequest true "Company information"
// @Success      200  {string}  string "Company updated successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/companies/profile [put]
func (h *CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	var req dto.CompanyUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Mapear a entidad
	company, err := request_mapper.CompanyUpdateRequestToCompany(&req)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Actualizar la empresa
	err = h.useCase.UpdateCompany(r.Context(), company)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "Company updated successfully")
}

// GetCompanyAddresses godoc
// @Summary      Obtener todas las direcciones de la compañia autenticada
// @Description  Obtener todas las direcciones de la compañia autenticada
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  dto.CompanyAddressResponse
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/companies/addresses [get]
func (h *CompanyHandler) GetCompanyAddresses(w http.ResponseWriter, r *http.Request) {
	addresses, err := h.useCase.GetCompanyAddresses(r.Context())
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Convertir a DTOs
	addressesResponse := make([]dto.CompanyAddressResponse, len(addresses))
	for i, address := range addresses {
		addressesResponse[i] = response_mapper.CompanyAddressToResponseDTO(address)
	}

	h.respWriter.Success(w, http.StatusOK, addressesResponse)
}

// GetCompanies godoc
// @Summary      Obtener lista de compañías
// @Description  Obtener lista de compañías con filtros y paginación
// @Tags         companies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        name         query  string  false  "Nombre de la empresa"
// @Param        legal_name   query  string  false  "Nombre legal de la empresa"
// @Param        tax_id       query  string  false  "Identificador fiscal"
// @Param        is_active    query  boolean false  "Estado activo/inactivo"
// @Param        page         query  int     false  "Número de página"  default(1)
// @Param        page_size    query  int     false  "Tamaño de página"  default(10)
// @Param        sort_by      query  string  false  "Campo para ordenar"
// @Param        sort_direction query  string  false  "Dirección (asc o desc)"  default(desc)
// @Success      200  {object}  dto.PaginatedResponse
// @Failure      401  {object}  responser.APIErrorResponse
// @Router       /api/v1/companies [get]
func (h *CompanyHandler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	// Obtener las empresas con filtros y paginación
	companies, params, total, err := h.useCase.GetCompanies(r.Context(), r)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Mapear las empresas a la respuesta DTO
	response := response_mapper.MapCompaniesToSimpleList(companies, params, total)
	h.respWriter.Success(w, http.StatusOK, response)
}

// AddCompanyAddress godoc
// @Summary      Agregar una nueva dirección para la compañia autenticada
// @Description Agregar una nueva dirección para la compañia autenticada
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        address body dto.CompanyAddressDTO true "Address information"
// @Success      201  {string}  string "Address added successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/companies/addresses [post]
func (h *CompanyHandler) AddCompanyAddress(w http.ResponseWriter, r *http.Request) {
	var req dto.CompanyAddressDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Mapear a entidad - el companyID se asignará desde el caso de uso
	address := request_mapper.CompanyAddressDTOToEntity(&req)

	// Añadir la dirección
	err := h.useCase.AddCompanyAddress(r.Context(), address)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusCreated, "Address added successfully")
}

// UpdateCompanyAddress godoc
// @Summary      Actualizar una dirección para la compañia autenticada
// @Description  Actualizar una dirección para la compañia autenticada
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        address_id path string true "Address ID"
// @Param        address body dto.CompanyAddressDTO true "Address information"
// @Success      200  {string}  string "Address updated successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/companies/addresses/{address_id} [put]
func (h *CompanyHandler) UpdateCompanyAddress(w http.ResponseWriter, r *http.Request) {
	// Extraer ID de la dirección
	vars := mux.Vars(r)
	addressID := vars["address_id"]

	var req dto.CompanyAddressDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Mapear a entidad
	address := request_mapper.CompanyAddressDTOToEntity(&req)

	// Actualizar la dirección
	err := h.useCase.UpdateCompanyAddress(r.Context(), addressID, address)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "Address updated successfully")
}

// DeleteCompanyAddress godoc
// @Summary      Eliminar una dirección de la compañia autenticada
// @Description  Eliminar una dirección de la compañia autenticada
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        address_id path string true "Address ID"
// @Success      200  {string}  string "Address deleted successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/companies/addresses/{address_id} [delete]
func (h *CompanyHandler) DeleteCompanyAddress(w http.ResponseWriter, r *http.Request) {
	// Extraer ID de la dirección
	vars := mux.Vars(r)
	addressID := vars["address_id"]

	// Eliminar la dirección
	err := h.useCase.DeleteCompanyAddress(r.Context(), addressID)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "Address deleted successfully")
}

// GetCompanyMetrics godoc
// @Summary      Obtener las métricas de la compañia
// @Description  Obtener las métricas de la compañia
// @Tags         companies, metrics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.CompanyMetricsResponse
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/companies/metrics [get]
func (h *CompanyHandler) GetCompanyMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.useCase.GetCompanyMetrics(r.Context())
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// Convertir a DTO
	response := response_mapper.CompanyToMetricsDTO(metrics)
	h.respWriter.Success(w, http.StatusOK, response)
}

// DeactivateCompany godoc
// @Summary      Deactivate the authenticated company
// @Description  Deactivate the authenticated company
// @Tags         companies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {string}  string "Company deactivated successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/companies/deactivate [post]
func (h *CompanyHandler) DeactivateCompany(w http.ResponseWriter, r *http.Request) {
	// El companyID se obtiene del contexto en el caso de uso
	err := h.useCase.DeactivateCompany(r.Context(), "")
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "Company deactivated successfully")
}

// ReactivateCompany godoc
// @Summary      Reactivate the authenticated company
// @Description  Reactivate the authenticated company
// @Tags         companies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {string}  string "Company reactivated successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/companies/reactivate [post]
func (h *CompanyHandler) ReactivateCompany(w http.ResponseWriter, r *http.Request) {
	// El companyID se obtiene del contexto en el caso de uso
	err := h.useCase.ReactivateCompany(r.Context(), "")
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "Company reactivated successfully")
}
