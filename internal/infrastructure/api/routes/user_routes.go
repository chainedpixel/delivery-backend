package routes

import (
	"github.com/gorilla/mux"
	"infrastructure/api/handlers"
	"net/http"
)

func RegisterUserRoutes(router *mux.Router, userHandler *handlers.UserHandler) {
	router.HandleFunc("/users/profile", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)

	router.HandleFunc("/users/profile", userHandler.GetUserProfile).Methods(http.MethodGet)
}
