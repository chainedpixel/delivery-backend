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
