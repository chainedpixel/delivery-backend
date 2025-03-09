package response_mapper

import (
	"domain/delivery/models/entities"
	"infrastructure/api/dto"
)

// OrderToResponseDTO - Mapper esencial para respuestas
func OrderToResponseDTO(order *entities.Order) *dto.OrderResponse {
	// Datos básicos para la respuesta
	response := &dto.OrderResponse{
		ID:             order.ID,
		CompanyID:      order.CompanyID,
		BranchID:       order.BranchID,
		ClientID:       order.ClientID,
		DriverID:       order.DriverID,
		TrackingNumber: order.TrackingNumber,
		Status:         order.Status,
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
	}

	// Mapear información de company/branch/cliente si está disponible
	if order.Company != nil {
		response.CompanyName = order.Company.Name
	}

	if order.Branch != nil {
		response.BranchName = order.Branch.Name
	}

	if order.Client != nil {
		response.ClientName = order.Client.FullName
	}

	// Mapear detalles esenciales del pedido
	if order.Detail != nil {
		response.Detail = dto.OrderDetailResponse{
			Price:             order.Detail.Price,
			Distance:          order.Detail.Distance,
			PickupTime:        order.Detail.PickupTime,
			DeliveryDeadline:  order.Detail.DeliveryDeadline,
			DeliveredAt:       order.Detail.DeliveredAt,
			RequiresSignature: order.Detail.RequiresSignature,
			DeliveryNotes:     order.Detail.DeliveryNotes,
		}
	}

	// Mapear información esencial de direcciones
	if order.DeliveryAddress != nil {
		response.DeliveryAddress = dto.DeliveryAddressResponse{
			RecipientName:  order.DeliveryAddress.RecipientName,
			RecipientPhone: order.DeliveryAddress.RecipientPhone,
			AddressLine1:   order.DeliveryAddress.AddressLine1,
			AddressLine2:   order.DeliveryAddress.AddressLine2,
			City:           order.DeliveryAddress.City,
			State:          order.DeliveryAddress.State,
			PostalCode:     order.DeliveryAddress.PostalCode,
			Latitude:       order.DeliveryAddress.Latitude,
			Longitude:      order.DeliveryAddress.Longitude,
		}
	}

	// Mapear información básica de estado de seguimiento
	response.CurrentStatus = order.Status
	if order.Tracking != nil {
		response.CurrentStatus = order.Tracking.CurrentStatus
		response.LastUpdated = order.Tracking.LastUpdated
	}

	return response
}
