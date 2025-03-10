package responser

import (
	domainErr "domain/error"
	"encoding/json"
	"errors"
	"fmt"
	errPackage "infrastructure/error"
	"net/http"
	"shared/logs"
	"strings"
)

type errorType string

const (
	errorValidation errorType = "VALIDATION" // error de validación de datos (Nivel de creacion en ValueObjects)
	errorBusiness   errorType = "BUSINESS"   // error de negocio (Nivel de creacion en procesos y concordancia de datos)
	errorSystem     errorType = "SYSTEM"     // error de sistema (Nivel de creacion en infraestructura y servicios)
)

type ResponseWriter struct{}

// NewResponseWriter crea una nueva instancia de ResponseWriter
func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{}
}

// Success envía una respuesta exitosa con el código de estado y los datos proporcionados.
func (w *ResponseWriter) Success(rw http.ResponseWriter, status int, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)

	json.NewEncoder(rw).Encode(&APIResponse{
		Success: true,
		Data:    data,
	})
}

// Error envía una respuesta de error con el código de estado y el mensaje proporcionado.
func (w *ResponseWriter) Error(rw http.ResponseWriter, status int, message string, details []string) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	json.NewEncoder(rw).Encode(APIErrorResponse{
		Error: &APIError{
			Message: message,
			Details: details,
			Code:    deriveErrorCode(status),
		},
	})
}

// HandleError maneja los diferentes tipos de errores y envía una respuesta de error con el código de estado y el mensaje correspondiente.
func (w *ResponseWriter) HandleError(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "application/json")

	logs.Error("Error processing request", map[string]interface{}{
		"error_type": getErrorType(err),
		"error":      err.Error(),
	})

	// TODO: Dejo el switch para que se pueda extender a futuro con más tipos de errores.
	switch errorType := getErrorType(err); errorType {
	case errorValidation:
		w.handleValidationError(rw, err)
	case errorBusiness:
		w.handleBusinessError(rw, err)
	default:
		w.handleSystemError(rw, err)
	}
}

// handleValidationError maneja los errores de validación y envía una respuesta de error con el código de estado y el mensaje correspondiente.
func (w *ResponseWriter) handleValidationError(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(rw).Encode(APIResponse{
		Success: false,
		Error: &APIError{
			Message: err.Error(),
			Code:    "VALIDATION_ERROR",
		},
	})
}

// handleBusinessError maneja los errores de negocio y envía una respuesta de error con el código de estado y el mensaje correspondiente.
func (w *ResponseWriter) handleBusinessError(rw http.ResponseWriter, err error) {
	var svcErr *errPackage.ServiceError
	if errors.As(err, &svcErr) {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(APIResponse{
			Success: false,
			Error: &APIError{
				Message: svcErr.Error(),
				Code:    fmt.Sprintf("BUSINESS_%s_ERROR", strings.ToUpper(svcErr.Type)),
			},
		})
		return
	}

	var domainErr *domainErr.DomainError
	if errors.As(err, &domainErr) {
		code := http.StatusBadRequest
		if domainErr.IsNotFoundError() {
			code = http.StatusNotFound
		}
		rw.WriteHeader(code)
		json.NewEncoder(rw).Encode(APIResponse{
			Success: false,
			Error: &APIError{
				Message: domainErr.Error(),
				Code:    deriveErrorCode(code),
			},
		})
		return
	}
}

// handleSystemError maneja los errores de sistema y envía una respuesta de error con el código de estado y el mensaje correspondiente.
func (w *ResponseWriter) handleSystemError(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(rw).Encode(APIResponse{
		Success: false,
		Error: &APIError{
			Message: "An unexpected error occurred",
			Code:    "SYSTEM_ERROR",
		},
	})
}

// deriveErrorCode deriva el código de error de acuerdo al estado y mensaje proporcionado.
func deriveErrorCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusMethodNotAllowed:
		return "METHOD_NOT_ALLOWED"
	case http.StatusInternalServerError:
		return "INTERNAL_SERVER_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}

// getErrorType obtiene el tipo de error de acuerdo al tipo de error proporcionado.
func getErrorType(err error) errorType {
	switch err.(type) {
	case *errPackage.ServiceError,
		*domainErr.DomainError:
		return errorBusiness
	default:
		return errorSystem
	}
}
