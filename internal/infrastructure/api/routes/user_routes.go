package routes

import (
	"github.com/gorilla/mux"
	"infrastructure/api/handlers"
	"net/http"
)

func RegisterUserRoutes(router *mux.Router, userHandler *handlers.UserHandler) {
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)
	router.HandleFunc("/users/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)
	router.HandleFunc("/users/profile", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)
	router.HandleFunc("/users/recover/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)
	router.HandleFunc("/users/roles/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)
	router.HandleFunc("/users/sessions/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)

	router.HandleFunc("/users/roles/{user_id}", userHandler.GetUserRoles).Methods(http.MethodGet)
	router.HandleFunc("/users/roles/{user_id}", userHandler.AssignRoleToUser).Methods(http.MethodPost)
	router.HandleFunc("/users/roles/{user_id}", userHandler.UnassignRole).Methods(http.MethodDelete)

	router.HandleFunc("/users/recover/{user_id}", userHandler.RecoverUser).Methods(http.MethodGet)
	router.HandleFunc("/users/sessions/{user_id}", userHandler.CleanAllSessions).Methods(http.MethodDelete)

	router.HandleFunc("/users", userHandler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/users", userHandler.GetAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/profile", userHandler.GetUserProfile).Methods(http.MethodGet)

	router.HandleFunc("/users/{user_id}", userHandler.GetUserByID).Methods(http.MethodGet)
	router.HandleFunc("/users/{user_id}", userHandler.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/users/{user_id}", userHandler.ActivateOrDeactivateUser).Methods(http.MethodPatch)
	router.HandleFunc("/users/{user_id}", userHandler.DeleteUser).Methods(http.MethodDelete)
}
