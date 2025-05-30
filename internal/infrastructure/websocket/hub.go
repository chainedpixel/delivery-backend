package websocket

import (
	"context"
	"encoding/json"
	ws "github.com/gorilla/websocket"
	"sync"
	"time"

	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/websocket"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
)

// Client representa una conexión WebSocket con un cliente
type Client struct {
	hub      *Hub
	conn     *ws.Conn
	userID   string
	orderIDs map[string]bool  // Pedidos a los que está suscrito
	send     chan interface{} // Canal para enviar mensajes al cliente
	mu       sync.Mutex       // Mutex para proteger el mapa de orderIDs
	ctx      context.Context
	cancel   context.CancelFunc
}

// Hub maneja todas las conexiones WebSocket activas
type Hub struct {
	// Mapa de todos los clientes conectados
	clients map[*Client]bool

	// Mapa de pedidos y los clientes que los observan
	orders map[string]map[*Client]bool

	// Canal para registrar nuevos clientes
	register chan *Client

	// Canal para eliminar clientes
	unregister chan *Client

	// Canal para enviar actualizaciones de pedidos
	orderUpdates chan *OrderUpdate

	// Canal para enviar actualizaciones de ubicación
	locationUpdates chan *LocationUpdate

	// Mutex para proteger los mapas
	mu sync.Mutex
}

// OrderUpdate representa una actualización del estado de un pedido
type OrderUpdate struct {
	OrderID string
	Data    *websocket.OrderUpdateData
}

// LocationUpdate representa una actualización de ubicación
type LocationUpdate struct {
	OrderID string
	Data    *websocket.LocationUpdateData
}

// NewHub crea una nueva instancia del Hub
func NewHub() *Hub {
	return &Hub{
		clients:         make(map[*Client]bool),
		orders:          make(map[string]map[*Client]bool),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		orderUpdates:    make(chan *OrderUpdate),
		locationUpdates: make(chan *LocationUpdate),
	}
}

// Run inicia el ciclo principal del Hub
func (h *Hub) Run() {
	logs.Info("Starting WebSocket hub", nil)
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case update := <-h.orderUpdates:
			h.broadcastOrderUpdate(update)

		case update := <-h.locationUpdates:
			h.broadcastLocationUpdate(update)
		}
	}
}

// RegisterClient registra un nuevo cliente
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	logs.Info("New client registered", map[string]interface{}{
		"user_id": client.userID,
	})
}

// UnregisterClient elimina un cliente y sus suscripciones
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)

		// Remover el cliente de todos los pedidos a los que estaba suscrito
		for orderID := range client.orderIDs {
			if _, ok := h.orders[orderID]; ok {
				delete(h.orders[orderID], client)

				// Si no quedan clientes observando este pedido, limpiar el mapa
				if len(h.orders[orderID]) == 0 {
					delete(h.orders, orderID)
				}
			}
		}

		logs.Info("Client unregistered", map[string]interface{}{
			"user_id": client.userID,
		})
	}
}

// SubscribeToOrder suscribe un cliente a las actualizaciones de un pedido
func (h *Hub) SubscribeToOrder(client *Client, orderID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Registrar la suscripción en el cliente
	client.mu.Lock()
	client.orderIDs[orderID] = true
	client.mu.Unlock()

	// Registrar el cliente en el mapa de pedidos
	if _, ok := h.orders[orderID]; !ok {
		h.orders[orderID] = make(map[*Client]bool)
	}
	h.orders[orderID][client] = true

	logs.Info("Client subscribed to order", map[string]interface{}{
		"user_id":  client.userID,
		"order_id": orderID,
	})
}

// UnsubscribeFromOrder cancela la suscripción de un cliente a un pedido
func (h *Hub) UnsubscribeFromOrder(client *Client, orderID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Eliminar la suscripción del cliente
	client.mu.Lock()
	delete(client.orderIDs, orderID)
	client.mu.Unlock()

	// Eliminar el cliente del mapa de pedidos
	if clients, ok := h.orders[orderID]; ok {
		delete(clients, client)

		// Si no quedan clientes, eliminar el pedido del mapa
		if len(clients) == 0 {
			delete(h.orders, orderID)
		}
	}

	logs.Info("Client unsubscribed from order", map[string]interface{}{
		"user_id":  client.userID,
		"order_id": orderID,
	})
}

// BroadcastOrderUpdate envía una actualización de pedido a todos los clientes suscritos
func (h *Hub) broadcastOrderUpdate(update *OrderUpdate) {
	h.mu.Lock()
	clients, ok := h.orders[update.OrderID]
	h.mu.Unlock()

	if !ok {
		return // No hay clientes suscritos a este pedido
	}

	message := websocket.Message{
		Type:      websocket.ServerOrderUpdate,
		OrderID:   update.OrderID,
		Timestamp: time.Now(),
		Data:      update.Data,
	}

	// Convertir a JSON una sola vez para reutilizar
	msgJSON, err := json.Marshal(message)
	if err != nil {
		logs.Error("Failed to marshal order update message", map[string]interface{}{
			"error":    err.Error(),
			"order_id": update.OrderID,
		})
		return
	}

	// Enviar a todos los clientes suscritos
	for client := range clients {
		select {
		case client.send <- msgJSON:
		default:
			// Si el canal está lleno, desconectar al cliente
			h.unregister <- client
		}
	}

	logs.Info("Order update broadcasted", map[string]interface{}{
		"order_id": update.OrderID,
		"status":   update.Data.Status,
		"clients":  len(clients),
	})
}

// BroadcastLocationUpdate envía una actualización de ubicación a todos los clientes suscritos
func (h *Hub) broadcastLocationUpdate(update *LocationUpdate) {
	if h == nil || h.orders == nil {
		logs.Warn("Hub or orders map is nil", nil)
		return
	}

	h.mu.Lock()
	clients, ok := h.orders[update.OrderID]
	h.mu.Unlock()

	if !ok || len(clients) == 0 {
		logs.Info("No clients subscribed to this order", map[string]interface{}{
			"order_id": update.OrderID,
		})
		return
	}

	message := websocket.Message{
		Type:      websocket.ServerLocation,
		OrderID:   update.OrderID,
		Timestamp: time.Now(),
		Data:      update.Data,
	}

	// Convertir a JSON una sola vez para reutilizar
	msgJSON, err := json.Marshal(message)
	if err != nil {
		logs.Error("Failed to marshal location update message", map[string]interface{}{
			"error":    err.Error(),
			"order_id": update.OrderID,
		})
		return
	}

	// Enviar a todos los clientes suscritos
	for client := range clients {
		select {
		case client.send <- msgJSON:
		default:
			// Si el canal está lleno, desconectar al cliente
			h.unregister <- client
		}
	}

	logs.Info("Location update broadcasted", map[string]interface{}{
		"order_id": update.OrderID,
		"clients":  len(clients),
	})
}

// SendOrderUpdate envía una actualización de pedido
func (h *Hub) SendOrderUpdate(orderID string, data *websocket.OrderUpdateData) {
	h.orderUpdates <- &OrderUpdate{
		OrderID: orderID,
		Data:    data,
	}
}

// SendLocationUpdate envía una actualización de ubicación
func (h *Hub) SendLocationUpdate(orderID string, data *websocket.LocationUpdateData) {
	logs.Debug("Details function", map[string]interface{}{
		"order_id":    orderID,
		"data is nil": data == nil,
		"latitude":    data.Latitude,
		"longitude":   data.Longitude,
	})

	h.locationUpdates <- &LocationUpdate{
		OrderID: orderID,
		Data:    data,
	}
}

// NewClient crea un nuevo cliente WebSocket
func NewClient(hub *Hub, conn *ws.Conn, userID string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		hub:      hub,
		conn:     conn,
		userID:   userID,
		orderIDs: make(map[string]bool),
		send:     make(chan interface{}, 256),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start inicia las goroutines de lectura y escritura del cliente
func (c *Client) Start() {
	// Registrar cliente en el hub
	c.hub.register <- c

	// Iniciar goroutines
	go c.readPump()
	go c.writePump()
}

// ReadPump maneja la recepción de mensajes del cliente
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		c.cancel()
	}()

	// Configurar la conexión WebSocket
	c.conn.SetReadLimit(4096) // Límite de tamaño de mensaje
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Ciclo de lectura
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				logs.Error("Unexpected close error", map[string]interface{}{
					"error":   err.Error(),
					"user_id": c.userID,
				})
			}
			break
		}

		// Procesar el mensaje
		var message websocket.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			logs.Error("Failed to unmarshal client message", map[string]interface{}{
				"error":   err.Error(),
				"user_id": c.userID,
			})
			continue
		}

		// Manejar según el tipo de mensaje
		switch message.Type {
		case websocket.ClientSubscribe:
			if message.OrderID != "" {
				c.hub.SubscribeToOrder(c, message.OrderID)
			}
		case websocket.ClientUnsubscribe:
			if message.OrderID != "" {
				c.hub.UnsubscribeFromOrder(c, message.OrderID)
			}
		}
	}
}

// WritePump maneja el envío de mensajes al cliente
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Canal cerrado, cerrar conexión
				c.conn.WriteMessage(ws.CloseMessage, []byte{})
				return
			}

			// Enviar el mensaje
			var err error
			switch msg := message.(type) {
			case []byte:
				// Mensaje ya convertido a JSON
				err = c.conn.WriteMessage(ws.TextMessage, msg)
			default:
				// Mensaje objeto, convertir a JSON
				err = c.conn.WriteJSON(message)
			}

			if err != nil {
				logs.Error("Failed to send message to client", map[string]interface{}{
					"error":   err.Error(),
					"user_id": c.userID,
				})
				return
			}

		case <-ticker.C:
			// Enviar ping periódicamente para mantener la conexión activa
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}

		case <-c.ctx.Done():
			// Contexto cancelado, cerrar goroutine
			return
		}
	}
}
