package entities

// CompanyMetrics contiene las métricas de una empresa
type CompanyMetrics struct {
	// Total de órdenes procesadas
	TotalOrders int64 `json:"total_orders"`

	// Total de órdenes completadas exitosamente
	CompletedOrders int64 `json:"completed_orders"`

	// Total de órdenes canceladas
	CancelledOrders int64 `json:"cancelled_orders"`

	// Tasa de entrega exitosa (porcentaje)
	DeliverySuccessRate float64 `json:"delivery_success_rate"`

	// Tiempo promedio de entrega (en minutos)
	AverageDeliveryTime float64 `json:"average_delivery_time"`

	// Ingresos totales generados
	TotalRevenue float64 `json:"total_revenue"`

	// Número de sucursales activas
	ActiveBranches int `json:"active_branches"`

	// Número de clientes únicos
	UniqueCustomers int `json:"unique_customers"`
}

// BranchMetrics contiene las métricas de una sucursal
type BranchMetrics struct {
	// Total de órdenes procesadas
	TotalOrders int64 `json:"total_orders"`

	// Total de órdenes completadas exitosamente
	CompletedOrders int64 `json:"completed_orders"`

	// Total de órdenes canceladas
	CancelledOrders int64 `json:"cancelled_orders"`

	// Tasa de entrega exitosa (porcentaje)
	DeliverySuccessRate float64 `json:"delivery_success_rate"`

	// Tiempo promedio de entrega (en minutos)
	AverageDeliveryTime float64 `json:"average_delivery_time"`

	// Ingresos totales generados
	TotalRevenue float64 `json:"total_revenue"`

	// Número de repartidores activos
	ActiveDrivers int `json:"active_drivers"`

	// Número de clientes únicos
	UniqueCustomers int `json:"unique_customers"`

	// Número de pedidos por hora pico
	PeakHourOrderRate float64 `json:"peak_hour_order_rate"`
}
