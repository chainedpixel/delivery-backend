package value_objects

import (
	"domain/delivery/constants"
	"strings"
)

type OrderStatus struct {
	value string
}

func NewOrderStatus(value string) *OrderStatus {
	return &OrderStatus{value: strings.ToUpper(value)}
}

func (s *OrderStatus) IsValid() bool {

	for _, status := range constants.ValidOrderStatuses {
		if s.value == status {
			return true
		}
	}
	return false
}

func (s *OrderStatus) ToString() string {
	return string(s.value)
}

func (s *OrderStatus) Equals(value ValidaterObject[string]) bool {
	return s.value == value.GetValue()
}

func (s *OrderStatus) GetValue() string {
	return s.value
}

func (s *OrderStatus) IsPending() bool {
	return s.value == constants.OrderStatusPending
}

func (s *OrderStatus) IsDelivered() bool {
	return s.value == constants.OrderStatusDelivered
}

func (s *OrderStatus) IsCancelled() bool {
	return s.value == constants.OrderStatusCancelled
}

func (s *OrderStatus) IsPickedUp() bool {
	return s.value == constants.OrderStatusPickedUp
}

func (s *OrderStatus) IsAccepted() bool {
	return s.value == constants.OrderStatusAccepted
}

func (s *OrderStatus) IsInWarehouse() bool {
	return s.value == constants.OrderStatusInWarehouse
}

func (s *OrderStatus) IsCompleted() bool {
	return s.value == constants.OrderStatusCompleted
}

func (s *OrderStatus) IsLost() bool {
	return s.value == constants.OrderStatusLost
}

func (s *OrderStatus) CanTransitionTo(nextStatus *OrderStatus) bool {
	validTransitions := map[string][]string{
		constants.OrderStatusPending:     {constants.OrderStatusAccepted, constants.OrderStatusCancelled},
		constants.OrderStatusAccepted:    {constants.OrderStatusPickedUp, constants.OrderStatusCancelled},
		constants.OrderStatusPickedUp:    {constants.OrderStatusInTransit, constants.OrderStatusInWarehouse, constants.OrderStatusCancelled},
		constants.OrderStatusInWarehouse: {constants.OrderStatusInTransit, constants.OrderStatusCancelled},
		constants.OrderStatusInTransit:   {constants.OrderStatusDelivered, constants.OrderStatusReturned, constants.OrderStatusCancelled},
		constants.OrderStatusDelivered:   {constants.OrderStatusReturned},
		constants.OrderStatusReturned:    {},
		constants.OrderStatusCancelled:   {},
	}

	for _, validNext := range validTransitions[s.value] {
		if validNext == nextStatus.value {
			return true
		}
	}
	return false
}
