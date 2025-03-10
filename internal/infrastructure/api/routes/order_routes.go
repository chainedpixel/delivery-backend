package routes

import (
	"github.com/gorilla/mux"
	"infrastructure/api/handlers"
	"net/http"
)

func RegisterOrderRoutes(router *mux.Router, orderHandler *handlers.OrderHandler) {
	router.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)
	router.HandleFunc("/orders/{order_id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)

	router.HandleFunc("/orders", orderHandler.CreateOrder).Methods(http.MethodPost)
	router.HandleFunc("/orders", orderHandler.GetOrdersByCompany).Methods(http.MethodGet)
	router.HandleFunc("/orders/{order_id}", orderHandler.GetOrderByID).Methods(http.MethodGet)
	router.HandleFunc("/orders/{order_id}", orderHandler.ChangeOrderStatus).Methods(http.MethodPatch)
}
