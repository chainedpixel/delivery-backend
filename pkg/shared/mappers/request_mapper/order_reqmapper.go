package request_mapper

import (
	"encoding/json"
	"fmt"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/constants"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/dto"
	"github.com/google/uuid"
	"time"
)

func OrderRequestToOrder(req *dto.OrderCreateRequest, companyAddress *entities.CompanyAddress) (*entities.Order, error) {
	orderID := uuid.NewString()

	// Crear objeto base del pedido
	order := &entities.Order{
		ID:        orderID,
		ClientID:  req.ClientID,
		Status:    constants.OrderStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Crear detalles del pedido (información esencial)
	order.Detail = &entities.Details{
		OrderID:           orderID,
		Price:             req.Price,
		Distance:          req.Distance,
		PickupTime:        req.PickupTime,
		DeliveryDeadline:  req.DeliveryDeadline,
		RequiresSignature: req.RequiresSignature,
		DeliveryNotes:     req.DeliveryNotes,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	var err error
	order.PackageDetail, err = createPackageDetail(req.PackageDetails, orderID)
	if err != nil {
		return nil, fmt.Errorf("error creating package details: %w", err)
	}

	// Datos de dirección de entrega
	order.DeliveryAddress = &entities.DeliveryAddress{
		OrderID:        orderID,
		RecipientName:  req.DeliveryAddress.RecipientName,
		RecipientPhone: req.DeliveryAddress.RecipientPhone,
		AddressLine1:   req.DeliveryAddress.AddressLine1,
		AddressLine2:   req.DeliveryAddress.AddressLine2,
		City:           req.DeliveryAddress.City,
		State:          req.DeliveryAddress.State,
		PostalCode:     req.DeliveryAddress.PostalCode,
		AddressNotes:   req.DeliveryAddress.AddressNotes,
		CreatedAt:      time.Now(),
	}

	// Datos de dirección de recogida
	order.PickupAddress = &entities.PickupAddress{
		OrderID: orderID,
		// Datos específicos de la recogida (de la solicitud)
		ContactName:  req.PickupContactName,
		ContactPhone: req.PickupContactPhone,

		// Datos geográficos (del company_address)
		AddressLine1: companyAddress.AddressLine1,
		AddressLine2: companyAddress.AddressLine2,
		City:         companyAddress.City,
		State:        companyAddress.State,
		PostalCode:   companyAddress.PostalCode,

		// Notas específicas para esta recogida
		AddressNotes: req.PickupNotes,
		CreatedAt:    time.Now(),
	}

	if companyAddress.Latitude != 0 && companyAddress.Longitude != 0 {
		order.PickupAddress.Latitude = companyAddress.Latitude
		order.PickupAddress.Longitude = companyAddress.Longitude
	}

	return order, nil
}

// En el mapper que procesa el DTO
func createPackageDetail(req dto.PackageDetailRequest, orderID string) (*entities.PackageDetail, error) {
	dimensionsJSON := ""
	if req.Length > 0 || req.Width > 0 || req.Height > 0 {
		dimensions := map[string]interface{}{
			"length": req.Length,
			"width":  req.Width,
			"height": req.Height,
			"unit":   "cm",
		}

		// Serializar a JSON
		dimensionsBytes, err := json.Marshal(dimensions)
		if err != nil {
			return nil, fmt.Errorf("error serializing dimensions: %w", err)
		}
		dimensionsJSON = string(dimensionsBytes)
	} else {
		dimensionsJSON = "{}"
	}

	return &entities.PackageDetail{
		OrderID:             orderID,
		IsFragile:           req.IsFragile,
		IsUrgent:            req.IsUrgent,
		Weight:              req.Weight,
		Dimensions:          dimensionsJSON,
		SpecialInstructions: req.SpecialInstructions,
		CreatedAt:           time.Now(),
	}, nil
}

func UpdateOrderFromRequest(orderID string, req *dto.OrderUpdateRequest) (*entities.Order, error) {
	// Crear una nueva orden vacía para actualización parcial
	order := &entities.Order{
		ID:        orderID,
		UpdatedAt: time.Now(),
	}

	// Añadir detail solo si hay campos a actualizar
	if hasDetailFields(req) {
		order.Detail = &entities.Details{
			OrderID:   orderID,
			UpdatedAt: time.Now(),
		}

		// Agregar solo los campos con valores
		if req.Price > 0 {
			order.Detail.Price = req.Price
		}
		if req.Distance > 0 {
			order.Detail.Distance = req.Distance
		}
		if req.PickupTime != nil {
			order.Detail.PickupTime = *req.PickupTime
		}
		if req.DeliveryDeadline != nil {
			order.Detail.DeliveryDeadline = *req.DeliveryDeadline
		}
		if req.RequiresSignature != nil {
			order.Detail.RequiresSignature = *req.RequiresSignature
		}
		if req.DeliveryNotes != "" {
			order.Detail.DeliveryNotes = req.DeliveryNotes
		}
	}

	// Añadir package detail solo si hay campos a actualizar
	if req.PackageDetails != nil && hasPackageDetailFields(req.PackageDetails) {
		order.PackageDetail = &entities.PackageDetail{
			OrderID: orderID,
		}

		if req.PackageDetails.IsFragile != nil {
			order.PackageDetail.IsFragile = *req.PackageDetails.IsFragile
		}
		if req.PackageDetails.IsUrgent != nil {
			order.PackageDetail.IsUrgent = *req.PackageDetails.IsUrgent
		}
		if req.PackageDetails.Weight != nil {
			order.PackageDetail.Weight = *req.PackageDetails.Weight
		}
		if req.PackageDetails.SpecialInstructions != "" {
			order.PackageDetail.SpecialInstructions = req.PackageDetails.SpecialInstructions
		}

		// Procesar dimensiones solo si hay alguna
		if req.PackageDetails.Length != nil || req.PackageDetails.Width != nil || req.PackageDetails.Height != nil {
			dimensions := map[string]interface{}{
				"unit": "cm",
			}

			if req.PackageDetails.Length != nil {
				dimensions["length"] = *req.PackageDetails.Length
			}
			if req.PackageDetails.Width != nil {
				dimensions["width"] = *req.PackageDetails.Width
			}
			if req.PackageDetails.Height != nil {
				dimensions["height"] = *req.PackageDetails.Height
			}

			dimensionsBytes, err := json.Marshal(dimensions)
			if err != nil {
				return nil, fmt.Errorf("error serializing dimensions: %w", err)
			}
			order.PackageDetail.Dimensions = string(dimensionsBytes)
		}
	}

	// Añadir delivery address solo si hay campos a actualizar
	if req.DeliveryAddress != nil && hasDeliveryAddressFields(req.DeliveryAddress) {
		order.DeliveryAddress = &entities.DeliveryAddress{
			OrderID: orderID,
		}

		if req.DeliveryAddress.RecipientName != "" {
			order.DeliveryAddress.RecipientName = req.DeliveryAddress.RecipientName
		}
		if req.DeliveryAddress.RecipientPhone != "" {
			order.DeliveryAddress.RecipientPhone = req.DeliveryAddress.RecipientPhone
		}
		if req.DeliveryAddress.AddressLine1 != "" {
			order.DeliveryAddress.AddressLine1 = req.DeliveryAddress.AddressLine1
		}
		if req.DeliveryAddress.AddressLine2 != "" {
			order.DeliveryAddress.AddressLine2 = req.DeliveryAddress.AddressLine2
		}
		if req.DeliveryAddress.City != "" {
			order.DeliveryAddress.City = req.DeliveryAddress.City
		}
		if req.DeliveryAddress.State != "" {
			order.DeliveryAddress.State = req.DeliveryAddress.State
		}
		if req.DeliveryAddress.PostalCode != "" {
			order.DeliveryAddress.PostalCode = req.DeliveryAddress.PostalCode
		}
		if req.DeliveryAddress.AddressNotes != "" {
			order.DeliveryAddress.AddressNotes = req.DeliveryAddress.AddressNotes
		}
	}

	// Añadir pickup address solo si hay campos a actualizar
	if hasPickupContactFields(req) {
		order.PickupAddress = &entities.PickupAddress{
			OrderID: orderID,
		}

		if req.PickupContactName != "" {
			order.PickupAddress.ContactName = req.PickupContactName
		}
		if req.PickupContactPhone != "" {
			order.PickupAddress.ContactPhone = req.PickupContactPhone
		}
		if req.PickupNotes != "" {
			order.PickupAddress.AddressNotes = req.PickupNotes
		}
	}

	return order, nil
}

// Funciones auxiliares para verificar si hay campos a actualizar

func hasDetailFields(req *dto.OrderUpdateRequest) bool {
	return req.Price > 0 ||
		req.Distance > 0 ||
		req.PickupTime != nil ||
		req.DeliveryDeadline != nil ||
		req.RequiresSignature != nil ||
		req.DeliveryNotes != ""
}

func hasPackageDetailFields(pkg *dto.PackageDetailUpdateRequest) bool {
	return pkg.IsFragile != nil ||
		pkg.IsUrgent != nil ||
		pkg.Weight != nil ||
		pkg.Length != nil ||
		pkg.Width != nil ||
		pkg.Height != nil ||
		pkg.SpecialInstructions != ""
}

func hasDeliveryAddressFields(addr *dto.DeliveryAddressUpdateRequest) bool {
	return addr.RecipientName != "" ||
		addr.RecipientPhone != "" ||
		addr.AddressLine1 != "" ||
		addr.AddressLine2 != "" ||
		addr.City != "" ||
		addr.State != "" ||
		addr.PostalCode != "" ||
		addr.AddressNotes != ""
}

func hasPickupContactFields(req *dto.OrderUpdateRequest) bool {
	return req.PickupContactName != "" ||
		req.PickupContactPhone != "" ||
		req.PickupNotes != ""
}
