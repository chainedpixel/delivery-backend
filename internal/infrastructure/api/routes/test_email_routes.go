package routes

import (
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

// RegisterEmailTestRoutes registra las rutas para pruebas de email
func RegisterEmailTestRoutes(router *mux.Router, emailTestHandler *handlers.EmailTestHandler) {
	router.HandleFunc("/test/email", emailTestHandler.SendTestEmail).Methods(http.MethodPost)
	router.HandleFunc("/test/email/status", emailTestHandler.GetEmailStatus).Methods(http.MethodGet)
	router.HandleFunc("/test/email/welcome", emailTestHandler.SendWelcomeEmail).Methods(http.MethodPost)
	router.HandleFunc("/test/email/order", emailTestHandler.SendOrderTestEmail).Methods(http.MethodPost)
}
