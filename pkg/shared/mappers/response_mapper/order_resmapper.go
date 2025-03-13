package response_mapper

import (
	"domain/delivery/constants"
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

	if order.PackageDetail != nil {
		response.PackageDetail = dto.PackageDetailResponse{
			Weight:              order.PackageDetail.Weight,
			IsFragile:           order.PackageDetail.IsFragile,
			IsUrgent:            order.PackageDetail.IsUrgent,
			Dimensions:          order.PackageDetail.Dimensions,
			SpecialInstructions: order.PackageDetail.SpecialInstructions,
		}
	}

	if order.StatusHistory != nil {
		response.StatusHistory = make([]dto.OrderStatusHistoryResponse, len(order.StatusHistory))
		for i, status := range order.StatusHistory {
			response.StatusHistory[i] = dto.OrderStatusHistoryResponse{
				Status:      status.Status,
				Description: status.Description,
				UpdatedAt:   status.CreatedAt.Format("2006-01-02 15:04:05"),
			}
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

	if order.PickupAddress != nil {
		response.PickupAddress = dto.PickupAddressResponse{
			ContactName:  order.PickupAddress.ContactName,
			ContactPhone: order.PickupAddress.ContactPhone,
			AddressLine1: order.PickupAddress.AddressLine1,
			AddressLine2: order.PickupAddress.AddressLine2,
			City:         order.PickupAddress.City,
			State:        order.PickupAddress.State,
			PostalCode:   order.PickupAddress.PostalCode,
			Latitude:     order.PickupAddress.Latitude,
			Longitude:    order.PickupAddress.Longitude,
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

// MapOrdersToResponse mapea las órdenes a DTOs de respuesta
func MapOrdersToResponse(orders []entities.Order, params *entities.OrderQueryParams, total int64) *dto.PaginatedResponse {
	response := make([]dto.OrderListResponse, len(orders))

	for i, order := range orders {
		var driverName string
		if order.Driver != nil {
			driverName = order.Driver.User.FullName
		}

		response[i] = dto.OrderListResponse{
			ID:               order.ID,
			TrackingNumber:   order.TrackingNumber,
			ClientName:       order.DeliveryAddress.RecipientName,
			DeliveryAddress:  order.DeliveryAddress.AddressLine1,
			DeliveryDeadline: order.Detail.DeliveryDeadline,
			Price:            order.Detail.Price,
			Status:           order.Status,
			DriverID:         order.DriverID,
			DriverName:       driverName,
			CreatedAt:        order.CreatedAt,
		}
	}

	return &dto.PaginatedResponse{
		Data:       response,
		TotalItems: total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: calculateTotalPages(total, params.PageSize),
	}
}

// MapUsersToResponse mapea los usuarios a DTOs de respuesta
func MapUsersToResponse(users []entities.User, params *entities.UserQueryParams, total int64) *dto.PaginatedResponse {
	response := make([]dto.UserListResponse, len(users))

	for i, user := range users {
		response[i] = dto.UserListResponse{
			ID:             user.ID,
			FullName:       user.FullName,
			Email:          user.Email,
			Phone:          user.Phone,
			DocumentType:   user.Profile.DocumentType,
			DocumentNumber: user.Profile.DocumentNumber,
			CreatedAt:      user.CreatedAt,
		}

		response[i].Role = constants.FinalUser
		if len(user.Roles) != 0 {
			response[i].Role = user.Roles[0].Role.Name
		}
	}

	return &dto.PaginatedResponse{
		Data:       response,
		TotalItems: total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: calculateTotalPages(total, params.PageSize),
	}
}

// calculateTotalPages calcula el número total de páginas
func calculateTotalPages(totalItems int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}

	pages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		pages++
	}

	return pages
}
