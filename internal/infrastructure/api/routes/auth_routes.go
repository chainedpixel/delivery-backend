package routes

import (
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterPublicAuthRoutes(router *mux.Router, authHandler *handlers.AuthHandler) {
	router.HandleFunc("/auth/login", authHandler.Login).Methods(http.MethodPost)
}

func RegisterProtectedAuthRoutes(router *mux.Router, authHandler *handlers.AuthHandler) {
	router.HandleFunc("/auth/logout", authHandler.Logout).Methods(http.MethodGet)
}
