package routes

import (
	"github.com/gorilla/mux"
	"infrastructure/api/handlers"
	"net/http"
)

func RegisterUserRoutes(router *mux.Router, userHandler *handlers.UserHandler) {
	router.HandleFunc("/users/profile", userHandler.GetUserProfile).Methods(http.MethodGet)
}
