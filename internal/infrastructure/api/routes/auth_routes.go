package routes

import (
	"github.com/gorilla/mux"
	"infrastructure/api/handlers"
	"net/http"
)

func RegisterPublicAuthRoutes(router *mux.Router, authHandler *handlers.AuthHandler) {
	router.HandleFunc("/auth/login", authHandler.Login).Methods(http.MethodPost)
}

func RegisterProtectedAuthRoutes(router *mux.Router, authHandler *handlers.AuthHandler) {
	router.HandleFunc("/auth/logout", authHandler.Logout).Methods(http.MethodGet)
}
