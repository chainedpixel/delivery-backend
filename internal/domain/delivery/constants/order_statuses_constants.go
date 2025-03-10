package constants

var (
	OrderStatusPending     = "PENDING"
	OrderStatusAccepted    = "ACCEPTED"
	OrderStatusCancelled   = "CANCELLED"
	OrderStatusDelivered   = "DELIVERED"
	OrderStatusPickedUp    = "PICKUP"
	OrderStatusInWarehouse = "IN_WAREHOUSE"
	OrderStatusInTransit   = "IN_TRANSIT"
	OrderStatusReturned    = "RETURNED"
	OrderStatusCompleted   = "COMPLETED"
	OrderStatusLost        = "LOST"
)

var ValidOrderStatuses = []string{
	OrderStatusPending,
	OrderStatusAccepted,
	OrderStatusCancelled,
	OrderStatusDelivered,
	OrderStatusCompleted,
	OrderStatusPickedUp,
	OrderStatusInWarehouse,
	OrderStatusReturned,
	OrderStatusInTransit,
	OrderStatusLost,
}
