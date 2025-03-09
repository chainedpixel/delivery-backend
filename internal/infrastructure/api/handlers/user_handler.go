package handlers

import (
	"application/ports"
	"domain/delivery/models/auth"
	"infrastructure/api/responser"
	errPackage "infrastructure/error"
	"net/http"
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
	// 1. Obtener el ID de los claims del token
	claims := r.Context().Value("claims").(*auth.AuthClaims)

	// 2. Obtener el perfil del usuario
	user, err := h.useCase.GetProfileInfo(r.Context(), claims.UserID)
	if err != nil {
		h.respWriter.HandleError(w, errPackage.NewGeneralServiceError("UserHandler", "GetUserProfile", err))
		return
	}

	h.respWriter.Success(w, http.StatusOK, user)
}
