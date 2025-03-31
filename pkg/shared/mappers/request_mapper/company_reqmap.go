package request_mapper

import (
	"encoding/json"
	"fmt"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/dto"
	"github.com/google/uuid"
	"time"
)

// CompanyRequestToCompany convierte un DTO de creación de empresa a una entidad de dominio
func CompanyRequestToCompany(req *dto.CompanyCreateRequest) (*entities.Company, error) {
	companyID := uuid.NewString()
	now := time.Now()

	// Crear objeto base de la empresa
	company := &entities.Company{
		ID:                companyID,
		Name:              req.Name,
		LegalName:         req.LegalName,
		TaxID:             req.TaxID,
		ContactEmail:      req.ContactEmail,
		ContactPhone:      req.ContactPhone,
		Website:           req.Website,
		IsActive:          true,
		DeliveryRate:      req.DeliveryRate,
		LogoURL:           req.LogoURL,
		ContractStartDate: req.ContractStartDate,
		ContractEndDate:   req.ContractEndDate,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// Procesar los detalles del contrato a formato JSON
	contractDetails, err := json.Marshal(req.ContractDetails)
	if err != nil {
		return nil, fmt.Errorf("error serializing contract details: %w", err)
	}
	company.ContractDetails = string(contractDetails)

	// Crear la dirección principal de la empresa
	mainAddress := &entities.CompanyAddress{
		ID:           uuid.NewString(),
		CompanyID:    companyID,
		AddressLine1: req.MainAddress.AddressLine1,
		AddressLine2: req.MainAddress.AddressLine2,
		City:         req.MainAddress.City,
		State:        req.MainAddress.State,
		PostalCode:   req.MainAddress.PostalCode,
		IsMain:       true,
		CreatedAt:    now,
	}

	// Establecer latitud y longitud para valores no almacenados en DB directamente
	mainAddress.Latitude = req.MainAddress.Latitude
	mainAddress.Longitude = req.MainAddress.Longitude

	// Agregar la dirección a la empresa
	company.Address = mainAddress

	return company, nil
}

// CompanyUpdateRequestToCompany convierte un DTO de actualización de empresa a una entidad de dominio
func CompanyUpdateRequestToCompany(req *dto.CompanyUpdateRequest) (*entities.Company, error) {
	// Crear objeto base de la empresa para actualización parcial
	company := &entities.Company{
		UpdatedAt: time.Now(),
	}

	// Agregar solo los campos con valores
	if req.Name != "" {
		company.Name = req.Name
	}

	if req.LegalName != "" {
		company.LegalName = req.LegalName
	}

	if req.ContactEmail != "" {
		company.ContactEmail = req.ContactEmail
	}

	if req.ContactPhone != "" {
		company.ContactPhone = req.ContactPhone
	}

	if req.Website != "" {
		company.Website = req.Website
	}

	if req.DeliveryRate != nil {
		company.DeliveryRate = *req.DeliveryRate
	}

	if req.LogoURL != "" {
		company.LogoURL = req.LogoURL
	}

	if req.ContractEndDate != nil {
		company.ContractEndDate = req.ContractEndDate
	}

	// Procesar los detalles del contrato si se proporcionan
	if req.ContractDetails != nil {
		contractDetails, err := json.Marshal(req.ContractDetails)
		if err != nil {
			return nil, fmt.Errorf("error serializing contract details: %w", err)
		}
		company.ContractDetails = string(contractDetails)
	}

	return company, nil
}

// CompanyAddressDTOToEntity convierte un DTO de dirección a una entidad de dominio
func CompanyAddressDTOToEntity(req *dto.CompanyAddressDTO) *entities.CompanyAddress {
	address := &entities.CompanyAddress{
		ID:           uuid.NewString(),
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		PostalCode:   req.PostalCode,
		IsMain:       req.IsMain,
		CreatedAt:    time.Now(),
	}

	// Establecer latitud y longitud para valores no almacenados en DB directamente
	address.Latitude = req.Latitude
	address.Longitude = req.Longitude

	return address
}
