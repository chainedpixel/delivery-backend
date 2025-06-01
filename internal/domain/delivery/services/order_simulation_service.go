package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/constants"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/interfaces"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/ports"
	errPackage "github.com/MarlonG1/delivery-backend/internal/domain/error"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
)

type OrderSimulationService struct {
	orderService   interfaces.Orderer
	companyRepo    ports.CompanyRepository
	userRepo       ports.UserRepository
	trackerService interfaces.OrderTracker
}

func NewOrderSimulationService(
	orderService interfaces.Orderer,
	companyRepo ports.CompanyRepository,
	userRepo ports.UserRepository,
	trackerService interfaces.OrderTracker,
) *OrderSimulationService {
	return &OrderSimulationService{
		orderService:   orderService,
		companyRepo:    companyRepo,
		userRepo:       userRepo,
		trackerService: trackerService,
	}
}

// SimulateOrderFlow simula todo el flujo de un pedido automáticamente
func (s *OrderSimulationService) SimulateOrderFlow(ctx context.Context, orderID string) error {
	logs.Info("Starting order simulation", map[string]interface{}{
		"orderID": orderID,
	})

	// 1. Verificar que el pedido existe
	order, err := s.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		return errPackage.NewDomainErrorWithCause("OrderSimulationService", "SimulateOrderFlow", "failed to get order", err)
	}

	// 2. Verificar que el pedido está en estado PENDING
	if order.Status != constants.OrderStatusPending {
		return errPackage.NewDomainError("OrderSimulationService", "SimulateOrderFlow",
			fmt.Sprintf("order must be in PENDING status, current: %s", order.Status))
	}

	// 3. Iniciar la simulación en una goroutine
	go s.executeSimulationFlow(context.Background(), orderID, order.CompanyID)

	logs.Info("Order simulation started", map[string]interface{}{
		"orderID": orderID,
	})

	return nil
}

// executeSimulationFlow ejecuta la simulación completa
func (s *OrderSimulationService) executeSimulationFlow(ctx context.Context, orderID, companyID string) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("Panic in order simulation", map[string]interface{}{
				"orderID": orderID,
				"error":   r,
			})
		}
	}()

	logs.Info("Starting simulation flow execution", map[string]interface{}{
		"orderID":   orderID,
		"companyID": companyID,
	})

	// Paso 1: Asignar driver automáticamente (opcional para simulación)
	if err := s.assignRandomDriverSafely(ctx, orderID, companyID); err != nil {
		logs.Warn("Could not assign real driver, continuing simulation without driver", map[string]interface{}{
			"orderID": orderID,
			"error":   err.Error(),
		})
		// Continuamos la simulación sin conductor real
	}

	// Delay antes del primer cambio de estado
	logs.Info("Simulation step 1: Waiting before ACCEPTED", map[string]interface{}{
		"orderID": orderID,
	})
	time.Sleep(3 * time.Second)

	// Paso 2: Cambiar a ACCEPTED
	logs.Info("Simulation step 2: Changing to ACCEPTED", map[string]interface{}{
		"orderID": orderID,
	})
	if err := s.changeOrderStatus(ctx, orderID, constants.OrderStatusAccepted); err != nil {
		logs.Error("Failed to change status to ACCEPTED", map[string]interface{}{
			"orderID": orderID,
			"error":   err.Error(),
		})
		return
	}

	// Delay antes de recoger
	waitTime := time.Duration(5+rand.Intn(10)) * time.Second
	logs.Info("Simulation step 3: Waiting before PICKED_UP", map[string]interface{}{
		"orderID":  orderID,
		"waitTime": waitTime.String(),
	})
	time.Sleep(waitTime)

	// Paso 3: Cambiar a PICKED_UP
	logs.Info("Simulation step 4: Changing to PICKED_UP", map[string]interface{}{
		"orderID": orderID,
	})
	if err := s.changeOrderStatus(ctx, orderID, constants.OrderStatusPickedUp); err != nil {
		logs.Error("Failed to change status to PICKED_UP", map[string]interface{}{
			"orderID": orderID,
			"error":   err.Error(),
		})
		return
	}

	// Delay antes de partir
	waitTime = time.Duration(3+rand.Intn(7)) * time.Second
	logs.Info("Simulation step 5: Waiting before IN_TRANSIT", map[string]interface{}{
		"orderID":  orderID,
		"waitTime": waitTime.String(),
	})
	time.Sleep(waitTime)

	// Paso 4: Cambiar a IN_TRANSIT
	logs.Info("Simulation step 6: Changing to IN_TRANSIT", map[string]interface{}{
		"orderID": orderID,
	})
	if err := s.changeOrderStatus(ctx, orderID, constants.OrderStatusInTransit); err != nil {
		logs.Error("Failed to change status to IN_TRANSIT", map[string]interface{}{
			"orderID": orderID,
			"error":   err.Error(),
		})
		return
	}

	// Simular movimiento del driver
	logs.Info("Simulation step 7: Starting driver movement simulation", map[string]interface{}{
		"orderID": orderID,
	})
	go s.simulateDriverMovement(ctx, orderID)

	// Delay del viaje
	waitTime = time.Duration(15+rand.Intn(30)) * time.Second
	logs.Info("Simulation step 8: Waiting for delivery", map[string]interface{}{
		"orderID":  orderID,
		"waitTime": waitTime.String(),
	})
	time.Sleep(waitTime)

	// Paso 5: Cambiar a DELIVERED
	logs.Info("Simulation step 9: Changing to DELIVERED", map[string]interface{}{
		"orderID": orderID,
	})
	if err := s.changeOrderStatus(ctx, orderID, constants.OrderStatusDelivered); err != nil {
		logs.Error("Failed to change status to DELIVERED", map[string]interface{}{
			"orderID": orderID,
			"error":   err.Error(),
		})
		return
	}

	logs.Info("Order simulation completed successfully", map[string]interface{}{
		"orderID": orderID,
	})
}

// assignRandomDriverSafely intenta asignar un driver real, pero no falla si no encuentra ninguno
func (s *OrderSimulationService) assignRandomDriverSafely(ctx context.Context, orderID, companyID string) error {
	// Obtener drivers disponibles de la empresa
	drivers, err := s.getAvailableDrivers(ctx, companyID)
	if err != nil {
		logs.Warn("Failed to get available drivers", map[string]interface{}{
			"orderID":   orderID,
			"companyID": companyID,
			"error":     err.Error(),
		})
		return s.simulateFakeDriverAssignment(ctx, orderID)
	}

	if len(drivers) == 0 {
		logs.Info("No real drivers available, simulating fake driver assignment", map[string]interface{}{
			"orderID": orderID,
		})
		return s.simulateFakeDriverAssignment(ctx, orderID)
	}

	// Seleccionar driver aleatorio
	selectedDriver := drivers[rand.Intn(len(drivers))]

	// Asignar el driver al pedido
	err = s.orderService.AssignDriverToOrder(ctx, orderID, selectedDriver.ID)
	if err != nil {
		logs.Warn("Failed to assign real driver, simulating fake assignment", map[string]interface{}{
			"orderID":  orderID,
			"driverID": selectedDriver.ID,
			"error":    err.Error(),
		})
		return s.simulateFakeDriverAssignment(ctx, orderID)
	}

	logs.Info("Real driver assigned to order", map[string]interface{}{
		"orderID":    orderID,
		"driverID":   selectedDriver.ID,
		"driverName": selectedDriver.FullName,
	})

	return nil
}

// simulateFakeDriverAssignment simula la asignación de un conductor ficticio
func (s *OrderSimulationService) simulateFakeDriverAssignment(ctx context.Context, orderID string) error {
	logs.Info("Simulating fake driver assignment for demonstration", map[string]interface{}{
		"orderID": orderID,
	})

	// En lugar de asignar un driver real, solo enviamos una notificación de simulación
	// El frontend puede mostrar información de driver simulado
	return nil
}

// getAvailableDrivers obtiene drivers disponibles de la empresa
func (s *OrderSimulationService) getAvailableDrivers(ctx context.Context, companyID string) ([]entities.User, error) {
	// Obtener usuarios de la empresa con rol de driver
	params := &entities.UserQueryParams{
		PaginationQueryParams: entities.PaginationQueryParams{
			Page:     1,
			PageSize: 100,
		},
		Status: true,
	}

	users, _, err := s.userRepo.GetAllUsersFromCompany(ctx, companyID, params)
	if err != nil {
		return nil, err
	}

	// Filtrar solo drivers
	var drivers []entities.User
	for _, user := range users {
		for _, userRole := range user.Roles {
			if userRole.Role != nil && userRole.Role.Name == constants.Driver {
				drivers = append(drivers, user)
				break
			}
		}
	}

	return drivers, nil
}

// changeOrderStatus cambia el estado del pedido
func (s *OrderSimulationService) changeOrderStatus(ctx context.Context, orderID, newStatus string) error {
	err := s.orderService.ChangeStatus(ctx, orderID, newStatus)
	if err != nil {
		return err
	}

	logs.Info("Order status changed", map[string]interface{}{
		"orderID":   orderID,
		"newStatus": newStatus,
	})

	return nil
}

// simulateDriverMovement simula el movimiento del conductor
func (s *OrderSimulationService) simulateDriverMovement(ctx context.Context, orderID string) {
	if s.trackerService == nil {
		logs.Warn("TrackerService is nil, cannot simulate driver movement", map[string]interface{}{
			"orderID": orderID,
		})
		return
	}

	logs.Info("Starting driver movement simulation", map[string]interface{}{
		"orderID": orderID,
	})

	// Coordenadas de ejemplo (San Salvador)
	startLat := 13.6929
	startLng := -89.2182
	endLat := 13.7942
	endLng := -89.1956

	// Simular movimiento por 20 pasos en 20 segundos
	steps := 20
	latIncrement := (endLat - startLat) / float64(steps)
	lngIncrement := (endLng - startLng) / float64(steps)

	currentLat := startLat
	currentLng := startLng

	for i := 0; i < steps; i++ {
		// Agregar algo de variación aleatoria para hacer más realista el movimiento
		randomLat := currentLat + (rand.Float64()*0.001 - 0.0005)
		randomLng := currentLng + (rand.Float64()*0.001 - 0.0005)

		// Enviar actualización de ubicación
		err := s.orderService.UpdateDriverLocation(ctx, orderID, randomLat, randomLng)
		if err != nil {
			logs.Error("Failed to update driver location", map[string]interface{}{
				"orderID": orderID,
				"step":    i + 1,
				"lat":     randomLat,
				"lng":     randomLng,
				"error":   err.Error(),
			})
		} else {
			logs.Debug("Driver location updated", map[string]interface{}{
				"orderID": orderID,
				"step":    i + 1,
				"lat":     randomLat,
				"lng":     randomLng,
			})
		}

		// Actualizar coordenadas para el siguiente paso
		currentLat += latIncrement
		currentLng += lngIncrement

		// Esperar 1 segundo antes del siguiente punto
		time.Sleep(1 * time.Second)
	}

	logs.Info("Driver movement simulation completed", map[string]interface{}{
		"orderID":    orderID,
		"totalSteps": steps,
		"finalLat":   currentLat,
		"finalLng":   currentLng,
	})
}

// SimulateRandomDriver simula solo la asignación de driver
func (s *OrderSimulationService) SimulateRandomDriver(ctx context.Context, orderID string) error {
	order, err := s.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	return s.assignRandomDriverSafely(ctx, orderID, order.CompanyID)
}

// GetSimulationStatus obtiene el estado actual de la simulación
func (s *OrderSimulationService) GetSimulationStatus(ctx context.Context, orderID string) (map[string]interface{}, error) {
	order, err := s.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	status := map[string]interface{}{
		"order_id":        orderID,
		"current_status":  order.Status,
		"has_driver":      order.DriverID != nil,
		"tracking_number": order.TrackingNumber,
		"created_at":      order.CreatedAt,
		"updated_at":      order.UpdatedAt,
	}

	if order.DriverID != nil {
		status["driver_id"] = *order.DriverID
	}

	// Determinar qué pasos faltan - actualizado con todos los estados
	steps := []string{
		constants.OrderStatusPending,
		constants.OrderStatusAccepted,
		constants.OrderStatusPickedUp,
		constants.OrderStatusInTransit,
		constants.OrderStatusDelivered,
	}

	currentIndex := -1
	for i, step := range steps {
		if step == order.Status {
			currentIndex = i
			break
		}
	}

	var remainingSteps []string
	if currentIndex >= 0 && currentIndex < len(steps)-1 {
		remainingSteps = steps[currentIndex+1:]
	}

	status["current_step"] = currentIndex + 1
	status["total_steps"] = len(steps)
	status["remaining_steps"] = remainingSteps
	status["completed"] = order.Status == constants.OrderStatusDelivered || order.Status == constants.OrderStatusCompleted

	return status, nil
}
