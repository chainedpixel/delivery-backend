package routes

import (
	"github.com/gorilla/mux"
	"infrastructure/api/handlers"
)

func RegisterAuthRoutes(router *mux.Router, authHandler *handlers.AuthHandler) {
	router.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
}
