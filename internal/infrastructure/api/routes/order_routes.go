package routes

import (
	"github.com/gorilla/mux"
	"infrastructure/api/handlers"
	"net/http"
)

func RegisterOrderRoutes(router *mux.Router, orderHandler *handlers.OrderHandler) {
	router.HandleFunc("/orders", orderHandler.CreateOrder).Methods(http.MethodPost)
	router.HandleFunc("/orders", orderHandler.GetOrdersByCompany).Methods(http.MethodGet)
	router.HandleFunc("/orders/{order_id}", orderHandler.GetOrderByID).Methods(http.MethodGet)
	router.HandleFunc("/orders/{order_id}", orderHandler.DeleteOrder).Methods(http.MethodDelete)
	router.HandleFunc("/orders/{order_id}", orderHandler.ChangeOrderStatus).Methods(http.MethodPatch)
	router.HandleFunc("/orders/{order_id}", orderHandler.UpdateOrder).Methods(http.MethodPut)
	router.HandleFunc("/orders/recovery/{order_id}", orderHandler.RestoreOrder).Methods(http.MethodGet)
}
