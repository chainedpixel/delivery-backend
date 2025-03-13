package handlers

import (
	"application/ports"
	"encoding/json"
	"github.com/gorilla/mux"
	"infrastructure/api/dto"
	"infrastructure/api/responser"
	errPackage "infrastructure/error"
	"net/http"
	"shared/mappers/request_mapper"
	"shared/mappers/response_mapper"
)

type UserHandler struct {
	useCase    ports.UserUseCase
	respWriter *responser.ResponseWriter
}

func NewUserHandler(useCase ports.UserUseCase) *UserHandler {
	return &UserHandler{
		useCase:    useCase,
		respWriter: responser.NewResponseWriter(),
	}
}

// GetUserProfile godoc
// @Summary      This endpoint is used to get the authenticated users profile information using the JWT token
// @Description  Get authenticated users profile information
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  entities.User
// @Failure      401  {object}  responser.APIErrorResponse
// @Router       /api/v1/users/profile [get]
func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener el perfil del usuario
	user, err := h.useCase.GetProfileInfo(r.Context())
	if err != nil {
		h.respWriter.HandleError(w, errPackage.NewGeneralServiceError("UserHandler", "GetUserProfile", err))
		return
	}

	h.respWriter.Success(w, http.StatusOK, user)
}

// CreateUser godoc
// @Summary      This endpoint is used to create a new user
// @Description  Create a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user body dto.UserDTO true "User object that needs to be created"
// @Success      201  string  "User created successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar el usuario del body
	var req dto.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 2. Verificar si la solicitud es válida
	if err := req.Validate(); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Mapear el DTO a la entidad
	user, err := request_mapper.UserRequestToModel(&req)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 4. Ejecutar el caso de uso
	err = h.useCase.CreateUser(r.Context(), user)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusCreated, "User created successfully")
}

// UpdateUser godoc
// @Summary      This endpoint is used to update a user by ID
// @Description  Update a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id path string true "User ID"
// @Param        user body dto.UpdateUserDTO true "User object that needs to be updated"
// @Success      200  string  "User updated successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/users/{user_id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del usuario
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// 2. Decodificar solicitud
	var req dto.UpdateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Mapear el DTO a la entidad
	user, err := request_mapper.UpdateUserRequestToModel(&req)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 4. Ejecutar el caso de uso
	err = h.useCase.UpdateUser(r.Context(), userID, user)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "User updated successfully")
}

// GetAllUsers godoc
// @Summary      This endpoint is used to get all users
// @Description  Get all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        company_id query string false "Company ID"
// @Param        page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param        sort_by query string false "Sort by"
// @Param        sort_direction query string false "Order by"
// @Param status query string false "Status"
// @Param creation_date query string false "Creation date"
// @Param include_deleted query string false "Include deleted"
// @Param 	  name query string false "Name"
// @Param 	  email query string false "Email"
// @Param 	  phone query string false "Phone"
// @Success      200  {object}  dto.PaginatedResponse
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/users [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// 1. Ejectuar el caso de uso
	users, params, total, err := h.useCase.GetAllUsers(r.Context(), r)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 2. Mapear a DTO
	response := response_mapper.MapUsersToResponse(users, params, total)

	h.respWriter.Success(w, http.StatusOK, response)
}

// GetUserByID godoc
// @Summary      This endpoint is used to get a user by ID
// @Description  Get a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id path string true "User ID"
// @Success      200  {object}  entities.User
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/users/{user_id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del usuario
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// 2. Ejecutar el caso de uso
	user, err := h.useCase.GetUserByID(r.Context(), userID)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, user)
}

// RecoverUser godoc
// @Summary      This endpoint is used to recover a user by ID
// @Description  Recover a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id path string true "User ID"
// @Success      200  string  "User recovered successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/users/recover/{user_id} [get]
func (h *UserHandler) RecoverUser(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del usuario
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// 2. Ejecutar el caso de uso
	err := h.useCase.RecoverUser(r.Context(), userID)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "User recovered successfully")
}

// DeleteUser godoc
// @Summary      This endpoint is used to delete a user by ID
// @Description  Delete a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id path string true "User ID"
// @Success      200  string  "User deleted successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/users/{user_id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del usuario
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// 2. Ejecutar el caso de uso
	err := h.useCase.DeleteUser(r.Context(), userID)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "User deleted successfully")
}

// ActivateOrDeactivateUser godoc
// @Summary      This endpoint is used to activate or deactivate a user by ID
// @Description  Activate or deactivate a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id path string true "User ID"
// @Param        active body dto.ActivateUserDTO true "Activate or deactivate user"
// @Success      200  string  "User activated or deactivated successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/users/{user_id} [patch]
func (h *UserHandler) ActivateOrDeactivateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del usuario
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// 2. Decodificar solicitud
	var req dto.ActivateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Ejecutar el caso de uso
	err := h.useCase.ActivateOrDeactivateUser(r.Context(), userID, req.Active)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "User activated or deactivated successfully")
}

// AssignRoleToUser godoc
// @Summary      This endpoint is used to assign a role to a user
// @Description  Assign a role to a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id path string true "User ID"
// @Param        role body dto.AssignRoleDTO true "Role object that needs to be assigned"
// @Success      200  string  "Role assigned to user successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/users/roles/{user_id} [post]
func (h *UserHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del usuario
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// 2. Decodificar solicitud
	var req dto.AssignRoleDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Verificar si la solicitud es válida
	if err := req.Validate(); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 4. Ejecutar el caso de uso
	err := h.useCase.AssignRoleToUser(r.Context(), userID, req.Role)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "Role assigned to user successfully")
}

// GetUserRoles godoc
// @Summary      This endpoint is used to get all roles of a user
// @Description  Get all roles of a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id path string true "User ID"
// @Success      200  {array}  entities.Role
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/users/roles/{user_id} [get]
func (h *UserHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del usuario
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// 2. Ejecutar el caso de uso
	roles, err := h.useCase.GetUserRoles(r.Context(), userID)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, roles)
}

// UnassignRole godoc
// @Summary      This endpoint is used to unassign a role from a user
// @Description  Unassign a role from a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id path string true "User ID"
// @Param        role body dto.AssignRoleDTO true "Role object that needs to be unassigned"
// @Success      200  string  "Role unassigned from user successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/users/roles/{user_id} [delete]
func (h *UserHandler) UnassignRole(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del usuario
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// 2. Decodificar solicitud
	var req dto.AssignRoleDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Verificar si la solicitud es válida
	if err := req.Validate(); err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 4. Ejecutar el caso de uso
	err := h.useCase.UnassignRole(r.Context(), userID, req.Role)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "Role unassigned from user successfully")
}

// CleanAllSessions godoc
// @Summary      This endpoint is used to clean all sessions of a user
// @Description  Clean all sessions of a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id path string true "User ID"
// @Success      200  string  "Sessions cleaned successfully"
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/users/sessions/{user_id} [delete]
func (h *UserHandler) CleanAllSessions(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del usuario
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// 2. Ejecutar el caso de uso
	err := h.useCase.CleanAllSessions(r.Context(), userID)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, "Sessions cleaned successfully")
}
