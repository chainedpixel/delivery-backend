package dto

import (
	"time"
)

// BranchCreateRequest representa la solicitud para crear una nueva sucursal
type BranchCreateRequest struct {
	// Nombre de la sucursal
	Name string `json:"name" example:"Sucursal Norte" binding:"required"`

	// Código único de la sucursal
	Code string `json:"code" example:"SUC-NORTE-001" binding:"required"`

	// Nombre del contacto en la sucursal
	ContactName string `json:"contact_name" example:"Gerente Norte" binding:"required"`

	// Teléfono del contacto en la sucursal
	ContactPhone string `json:"contact_phone" example:"+573001112233" binding:"required"`

	// Email del contacto en la sucursal
	ContactEmail string `json:"contact_email" example:"norte@expressdelivery.com" binding:"required,email"`

	// ID de la zona a la que pertenece la sucursal
	ZoneID string `json:"zone_id" example:"f8c3e8d7-b6a5-4d3c-9f1e-0a2b4c6d8e0f" binding:"required"`

	// Horarios de operación de la sucursal
	OperatingHours OperatingHoursDTO `json:"operating_hours" binding:"required"`

	// Dirección asociada a la sucursal (opcional)
	Address *CompanyAddressDTO `json:"address,omitempty"`
}

// OperatingHoursDTO representa los horarios de operación
type OperatingHoursDTO struct {
	// Horario para días de semana
	Weekdays OperatingHoursTimeDTO `json:"weekdays" binding:"required"`

	// Horario para fines de semana
	Weekends OperatingHoursTimeDTO `json:"weekends" binding:"required"`
}

// OperatingHoursTimeDTO representa un horario de inicio y fin
type OperatingHoursTimeDTO struct {
	// Hora de inicio (formato HH:MM)
	Start string `json:"start" example:"08:00" binding:"required"`

	// Hora de fin (formato HH:MM)
	End string `json:"end" example:"20:00" binding:"required"`
}

// BranchUpdateRequest representa la solicitud para actualizar una sucursal existente
type BranchUpdateRequest struct {
	// Nombre de la sucursal
	Name string `json:"name,omitempty" example:"Sucursal Norte"`

	// Código único de la sucursal
	Code string `json:"code,omitempty" example:"SUC-NORTE-001"`

	// Nombre del contacto en la sucursal
	ContactName string `json:"contact_name,omitempty" example:"Gerente Norte"`

	// Teléfono del contacto en la sucursal
	ContactPhone string `json:"contact_phone,omitempty" example:"+573001112233"`

	// Email del contacto en la sucursal
	ContactEmail string `json:"contact_email,omitempty" example:"norte@expressdelivery.com" binding:"omitempty,email"`

	// ID de la zona a la que pertenece la sucursal
	ZoneID string `json:"zone_id,omitempty" example:"f8c3e8d7-b6a5-4d3c-9f1e-0a2b4c6d8e0f"`

	// Horarios de operación de la sucursal
	OperatingHours *OperatingHoursDTO `json:"operating_hours,omitempty"`
}

// BranchResponse representa la respuesta con información de una sucursal
type BranchResponse struct {
	// ID único de la sucursal
	ID string `json:"id" example:"b5f8c3d1-2e59-4c4b-a6e8-e5f3c0c3d1b5"`

	// ID de la empresa a la que pertenece
	CompanyID string `json:"company_id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`

	// Nombre de la empresa
	CompanyName string `json:"company_name,omitempty" example:"Express Delivery Co."`

	// Nombre de la sucursal
	Name string `json:"name" example:"Sucursal Norte"`

	// Código único de la sucursal
	Code string `json:"code" example:"SUC-NORTE-001"`

	// Nombre del contacto en la sucursal
	ContactName string `json:"contact_name" example:"Gerente Norte"`

	// Teléfono del contacto en la sucursal
	ContactPhone string `json:"contact_phone" example:"+573001112233"`

	// Email del contacto en la sucursal
	ContactEmail string `json:"contact_email" example:"norte@expressdelivery.com"`

	// Indica si la sucursal está activa
	IsActive bool `json:"is_active" example:"true"`

	// ID de la zona a la que pertenece la sucursal
	ZoneID string `json:"zone_id" example:"f8c3e8d7-b6a5-4d3c-9f1e-0a2b4c6d8e0f"`

	// Nombre de la zona
	ZoneName string `json:"zone_name,omitempty" example:"Zona Norte"`

	// Horarios de operación de la sucursal
	OperatingHours string `json:"operating_hours" example:"{\"weekdays\":{\"start\":\"08:00\",\"end\":\"20:00\"},\"weekends\":{\"start\":\"09:00\",\"end\":\"17:00\"}}"`

	// Cuando se creó la sucursal
	CreatedAt time.Time `json:"created_at" format:"date-time"`

	// Última actualización de la sucursal
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`

	// Dirección asociada a la sucursal
	Address *CompanyAddressResponse `json:"address,omitempty"`

	// Métricas de la sucursal (opcional)
	Metrics *BranchMetricsResponse `json:"metrics,omitempty"`
}

// BranchActivationRequest representa la solicitud para activar o desactivar una sucursal
type BranchActivationRequest struct {
	// Indica si la sucursal debe estar activa
	IsActive bool `json:"is_active" example:"true" binding:"required"`

	// Razón del cambio de estado
	Reason string `json:"reason" example:"Renovación de local" binding:"required"`
}

// BranchZoneAssignRequest representa la solicitud para asignar una sucursal a una zona
type BranchZoneAssignRequest struct {
	// ID de la zona a asignar
	ZoneID string `json:"zone_id" example:"f8c3e8d7-b6a5-4d3c-9f1e-0a2b4c6d8e0f" binding:"required"`
}

// BranchMetricsResponse representa las métricas de una sucursal
type BranchMetricsResponse struct {
	// Total de órdenes procesadas
	TotalOrders int64 `json:"total_orders" example:"250"`

	// Total de órdenes completadas exitosamente
	CompletedOrders int64 `json:"completed_orders" example:"220"`

	// Total de órdenes canceladas
	CancelledOrders int64 `json:"cancelled_orders" example:"15"`

	// Tasa de entrega exitosa (porcentaje)
	DeliverySuccessRate float64 `json:"delivery_success_rate" example:"88.0"`

	// Tiempo promedio de entrega (en minutos)
	AverageDeliveryTime float64 `json:"average_delivery_time" example:"40.2"`

	// Ingresos totales generados
	TotalRevenue float64 `json:"total_revenue" example:"5000.75"`

	// Número de repartidores activos
	ActiveDrivers int `json:"active_drivers" example:"10"`

	// Número de clientes únicos
	UniqueCustomers int `json:"unique_customers" example:"120"`

	// Número de pedidos por hora pico
	PeakHourOrderRate float64 `json:"peak_hour_order_rate" example:"12.5"`
}

// BranchListResponse representa un item en la lista de sucursales
type BranchListResponse struct {
	// ID único de la sucursal
	ID string `json:"id" example:"b5f8c3d1-2e59-4c4b-a6e8-e5f3c0c3d1b5"`

	// ID de la empresa a la que pertenece
	CompanyID string `json:"company_id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`

	// Nombre de la empresa
	CompanyName string `json:"company_name" example:"Express Delivery Co."`

	// Nombre de la sucursal
	Name string `json:"name" example:"Sucursal Norte"`

	// Código único de la sucursal
	Code string `json:"code" example:"SUC-NORTE-001"`

	// Nombre del contacto en la sucursal
	ContactName string `json:"contact_name" example:"Gerente Norte"`

	// Email del contacto en la sucursal
	ContactEmail string `json:"contact_email" example:"norte@expressdelivery.com"`

	// Indica si la sucursal está activa
	IsActive bool `json:"is_active" example:"true"`

	// Nombre de la zona
	ZoneName string `json:"zone_name" example:"Zona Norte"`

	// Total de órdenes en proceso
	ActiveOrders int `json:"active_orders" example:"15"`

	// Cuando se creó la sucursal
	CreatedAt time.Time `json:"created_at" format:"date-time"`
}
