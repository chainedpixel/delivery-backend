package request_mapper

import (
	"domain/delivery/constants"
	"domain/delivery/models/entities"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"infrastructure/api/dto"
	"time"
)

func OrderRequestToOrder(req *dto.OrderCreateRequest, companyAddress *entities.CompanyAddress) (*entities.Order, error) {
	orderID := uuid.NewString()

	// Crear objeto base del pedido
	order := &entities.Order{
		ID:        orderID,
		CompanyID: req.CompanyID,
		BranchID:  req.BranchID,
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
