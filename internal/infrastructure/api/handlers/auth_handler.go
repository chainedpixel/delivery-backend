package handlers

import (
	"application/ports"
	"infrastructure/api/dto"
	"infrastructure/api/responser"
	errPackage "infrastructure/error"
	"net"
	"net/http"
	"shared/logs"
	"strings"
)

type AuthHandler struct {
	authUseCase ports.AuthenticatorUseCase
	respWriter  *responser.ResponseWriter
}

func NewAuthHandler(authUseCase ports.AuthenticatorUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		respWriter:  responser.NewResponseWriter(),
	}
}

// Login godoc
// @Summary      This endpoint is used to authenticate a users and return a JWT token to be used in subsequent requests
// @Description  Authenticate users and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login credentials"
// @Success      200  {object}  dto.LoginResponse
// @Failure      400  {object}  responser.APIErrorResponse
// @Failure      401  {object}  responser.APIErrorResponse
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener credenciales
	req, err := dto.NewLoginRequest(r.Body)
	if err != nil {
		h.respWriter.HandleError(w, errPackage.NewGeneralServiceError("AuthHandler", "Login", err))
		return
	}

	// 2. Autenticar
	token, err := h.authUseCase.Authenticate(r.Context(), req.ParseToCredentialsModel(getClientIP(r)))
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	logs.Debug("Client IP information", map[string]interface{}{
		"x_forwarded_for": r.Header.Get("X-Forwarded-For"),
		"x_real_ip":       r.Header.Get("X-Real-IP"),
		"remote_addr":     r.RemoteAddr,
		"final_ip":        getClientIP(r),
	})

	// 3. Responder
	h.respWriter.Success(w, http.StatusOK, dto.LoginResponse{
		Token: token,
	})
}

// Logout godoc
// @Summary      This endpoint is used to logout a users and invalidate the JWT token
// @Description  Logout users and invalidate JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  responser.APIErrorResponse
// @Failure      401  {object}  responser.APIErrorResponse
// @Router       /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener token
	token := r.Context().Value("userToken").(string)

	// 2. Desautenticar
	err := h.authUseCase.SignOut(r.Context(), token)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Responder
	h.respWriter.Success(w, http.StatusOK, map[string]interface{}{
		"message": "User logged out successfully",
	})
}

// getClientIP obtiene la dirección IP del cliente
// Se intenta obtener la dirección IP desde los headers X-Forwarded-For y X-Real-IP
// Si no se encuentra, se obtiene la dirección IP desde RemoteAddr
// Si no se puede obtener la dirección IP desde RemoteAddr, se retorna la dirección IP de RemoteAddr
func getClientIP(r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")

	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
