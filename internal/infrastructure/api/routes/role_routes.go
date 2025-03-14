package routes

import (
	"github.com/gorilla/mux"
	"infrastructure/api/handlers"
	"net/http"
)

func RegisterRoleRoutes(router *mux.Router, roleHandler *handlers.RoleHandler) {
	router.HandleFunc("/roles", roleHandler.GetRoles).Methods(http.MethodGet)
	router.HandleFunc("/roles/{role}", roleHandler.GetRole).Methods(http.MethodGet)
}
