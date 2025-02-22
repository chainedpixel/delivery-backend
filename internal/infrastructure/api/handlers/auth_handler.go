package handlers

import (
	"application/ports"
	"infrastructure/api/dto"
	"infrastructure/api/responser"
	"net/http"
)

type AuthHandler struct {
	authUseCase ports.AuthUseCase
	respWriter  *responser.ResponseWriter
}

func NewAuthHandler(authUseCase ports.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		respWriter:  responser.NewResponseWriter(),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener credenciales
	req, err := dto.NewLoginRequest(r.Body)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 2. Autenticar
	token, err := h.authUseCase.Authenticate(r.Context(), req.ParseToCredentialsModel(r.RemoteAddr))
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Responder
	h.respWriter.Success(w, http.StatusOK, dto.LoginResponse{
		Token: token,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener token
	token := r.Context().Value("systemToken").(string)

	// 2. Desautenticar
	err := h.authUseCase.SignOut(r.Context(), token)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Responder
	h.respWriter.Success(w, http.StatusOK, nil)
}
