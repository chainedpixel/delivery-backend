package request_mapper

import (
	"encoding/json"
	"fmt"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/dto"
	"github.com/google/uuid"
	"time"
)

// BranchRequestToBranch convierte un DTO de creación de sucursal a una entidad de dominio
func BranchRequestToBranch(companyID string, req *dto.BranchCreateRequest) (*entities.Branch, error) {
	now := time.Now()

	// Crear objeto base de la sucursal
	branch := &entities.Branch{
		ID:           uuid.NewString(),
		CompanyID:    companyID,
		Name:         req.Name,
		Code:         req.Code,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		ContactEmail: req.ContactEmail,
		IsActive:     true,
		ZoneID:       req.ZoneID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Procesar los horarios de operación a formato JSON
	operatingHours, err := json.Marshal(req.OperatingHours)
	if err != nil {
		return nil, fmt.Errorf("error serializing operating hours: %w", err)
	}
	branch.OperatingHours = string(operatingHours)

	return branch, nil
}

// BranchUpdateRequestToBranch convierte un DTO de actualización de sucursal a una entidad de dominio
func BranchUpdateRequestToBranch(id string, req *dto.BranchUpdateRequest) (*entities.Branch, error) {
	// Crear objeto base de la sucursal para actualización parcial
	branch := &entities.Branch{
		ID:        id,
		UpdatedAt: time.Now(),
	}

	// Agregar solo los campos con valores
	if req.Name != "" {
		branch.Name = req.Name
	}

	if req.Code != "" {
		branch.Code = req.Code
	}

	if req.ContactName != "" {
		branch.ContactName = req.ContactName
	}

	if req.ContactPhone != "" {
		branch.ContactPhone = req.ContactPhone
	}

	if req.ContactEmail != "" {
		branch.ContactEmail = req.ContactEmail
	}

	if req.ZoneID != "" {
		branch.ZoneID = req.ZoneID
	}

	// Procesar los horarios de operación si se proporcionan
	if req.OperatingHours != nil {
		operatingHours, err := json.Marshal(req.OperatingHours)
		if err != nil {
			return nil, fmt.Errorf("error serializing operating hours: %w", err)
		}
		branch.OperatingHours = string(operatingHours)
	}

	return branch, nil
}
