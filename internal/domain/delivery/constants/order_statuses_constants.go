package constants

var (
	OrderStatusPending     = "PENDING"
	OrderStatusAccepted    = "ACCEPTED"
	OrderStatusCancelled   = "CANCELLED"
	OrderStatusDelivered   = "DELIVERED"
	OrderStatusPickedUp    = "PICKED_UP"
	OrderStatusInWarehouse = "+"
	OrderStatusInTransit   = "IN_TRANSIT"
	OrderStatusReturned    = "RETURNED"
	OrderStatusCompleted   = "COMPLETED"
	OrderStatusLost        = "LOST"
	OrderStatusDeleted     = "DELETED"
	OrderStatusRestored    = "RESTORED"
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

var AllowedStatesToDelete = map[string]bool{
	OrderStatusPending:   true,
	OrderStatusCancelled: true,
	OrderStatusRestored:  true,
}

var AllowedStatesToUpdate = map[string]bool{
	OrderStatusPending:     true,
	OrderStatusAccepted:    true,
	OrderStatusPickedUp:    true,
	OrderStatusInWarehouse: true,
	OrderStatusInTransit:   true,
}
