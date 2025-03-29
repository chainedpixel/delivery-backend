package routes

import (
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterRoleRoutes(router *mux.Router, roleHandler *handlers.RoleHandler) {
	router.HandleFunc("/roles", roleHandler.GetRoles).Methods(http.MethodGet)
	router.HandleFunc("/roles/{role}", roleHandler.GetRole).Methods(http.MethodGet)
}
