package websocket

import (
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"time"
)

// MessageType define el tipo de mensaje que se enviará por websocket
type MessageType string

const (
	// Tipos de mensajes del cliente al servidor
	ClientSubscribe   MessageType = "SUBSCRIBE"   // Cliente solicita suscribirse a actualizaciones de un pedido
	ClientUnsubscribe MessageType = "UNSUBSCRIBE" // Cliente solicita cancelar la suscripción

	// Tipos de mensajes del servidor al cliente
	ServerOrderUpdate MessageType = "ORDER_UPDATE" // Actualización del estado del pedido
	ServerLocation    MessageType = "LOCATION"     // Actualización de la ubicación del repartidor
	ServerError       MessageType = "ERROR"        // Mensaje de error
)

// Message representa un mensaje genérico de WebSocket
type Message struct {
	Type      MessageType `json:"type"`               // Tipo de mensaje
	OrderID   string      `json:"order_id,omitempty"` // ID del pedido (cuando aplica)
	Timestamp time.Time   `json:"timestamp"`          // Hora del mensaje
	Data      interface{} `json:"data,omitempty"`     // Datos del mensaje (depende del tipo)
}

// OrderUpdateData contiene los datos para un mensaje de actualización de pedido
type OrderUpdateData struct {
	Status      string     `json:"status"`                // Estado actual del pedido
	Description string     `json:"description,omitempty"` // Descripción opcional
	UpdatedAt   time.Time  `json:"updated_at"`            // Hora de la actualización
	Order       *OrderInfo `json:"order,omitempty"`       // Información resumida del pedido
}

// LocationUpdateData contiene los datos de actualización de ubicación
type LocationUpdateData struct {
	Latitude  float64   `json:"latitude"`  // Latitud del repartidor
	Longitude float64   `json:"longitude"` // Longitud del repartidor
	UpdatedAt time.Time `json:"updated_at"`
	Address   string    `json:"address,omitempty"` // Dirección aproximada (opcional)
}

// ErrorData contiene información sobre un error
type ErrorData struct {
	Code    string `json:"code"`              // Código de error
	Message string `json:"message"`           // Mensaje de error
	Details string `json:"details,omitempty"` // Detalles adicionales (opcional)
}

// OrderInfo contiene información resumida del pedido para actualizaciones
type OrderInfo struct {
	ID             string  `json:"id"`
	TrackingNumber string  `json:"tracking_number"`
	Status         string  `json:"status"`
	DriverName     string  `json:"driver_name,omitempty"`
	EstimatedTime  int     `json:"estimated_time,omitempty"` // Tiempo estimado en minutos
	CompanyName    string  `json:"company_name"`
	Progress       float64 `json:"progress"` // Porcentaje de progreso (0-100)
}

// OrderInfoFromEntity convierte una entidad Order a un OrderInfo
func OrderInfoFromEntity(order *entities.Order) *OrderInfo {
	info := &OrderInfo{
		ID:             order.ID,
		TrackingNumber: order.TrackingNumber,
		Status:         order.Status,
		CompanyName:    order.Company.Name,
	}

	// Agregar datos del repartidor si existe
	if order.Driver != nil && order.Driver.User != nil {
		info.DriverName = order.Driver.User.FullName
	}

	// Calcular progreso según el estado
	switch order.Status {
	case "PENDING":
		info.Progress = 10
	case "ACCEPTED":
		info.Progress = 25
	case "PICKED_UP":
		info.Progress = 50
	case "IN_TRANSIT":
		info.Progress = 75
	case "DELIVERED":
		info.Progress = 100
	default:
		info.Progress = 0
	}

	// Calcular tiempo estimado basado en info del detalle si existe
	if order.Detail != nil && !order.Detail.DeliveryDeadline.IsZero() {
		// Tiempo estimado en minutos desde ahora hasta la fecha de entrega
		estimatedTime := int(time.Until(order.Detail.DeliveryDeadline).Minutes())
		if estimatedTime > 0 {
			info.EstimatedTime = estimatedTime
		}
	}

	return info
}
