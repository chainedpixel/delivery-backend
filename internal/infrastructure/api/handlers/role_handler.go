package handlers

import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/MarlonG1/delivery-backend/internal/application/ports"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/responser"

	_ "github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
)

type RoleHandler struct {
	roleUseCase ports.RolerUseCase
	respWriter  *responser.ResponseWriter
}

func NewRoleHandler(roleUseCase ports.RolerUseCase) *RoleHandler {
	return &RoleHandler{
		roleUseCase: roleUseCase,
		respWriter:  responser.NewResponseWriter(),
	}
}

// GetRoles godoc
// @Summary Get all roles
// @Description Get all roles
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.Role
// @Failure 400 {object} responser.APIErrorResponse
// @Router /api/v1/roles [get]
func (h *RoleHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener todos los roles
	roles, err := h.roleUseCase.GetRoles(r.Context())
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 2. Responder
	h.respWriter.Success(w, http.StatusOK, roles)
}

// GetRole godoc
// @Summary Get role by ID or name
// @Description Get role by ID or name
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role path string true "Role ID or name"
// @Success 200 {object} entities.Role
// @Failure 400 {object} responser.APIErrorResponse
// @Router /api/v1/roles/{role} [get]
func (h *RoleHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID o nombre del rol
	vars := mux.Vars(r)
	roleString := vars["role"]

	// 2. Obtener el rol por su ID o nombre
	role, err := h.roleUseCase.GetRoleByIDOrName(r.Context(), roleString)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, role)
}
