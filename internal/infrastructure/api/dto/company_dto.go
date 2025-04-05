package dto

import (
	"time"
)

// CompanyCreateRequest representa la solicitud para crear una nueva empresa
type CompanyCreateRequest struct {
	// Nombre comercial de la empresa
	Name string `json:"name" example:"Express Delivery Co." binding:"required"`

	// Nombre legal completo de la empresa
	LegalName string `json:"legal_name" example:"Express Delivery S.A.S" binding:"required"`

	// Identificador fiscal de la empresa
	TaxID string `json:"tax_id" example:"900123456-7" binding:"required"`

	// Email de contacto principal de la empresa
	ContactEmail string `json:"contact_email" example:"contacto@expressdelivery.com" binding:"required,email"`

	// Teléfono de contacto principal
	ContactPhone string `json:"contact_phone" example:"+573001234567" binding:"required"`

	// Sitio web de la empresa (opcional)
	Website string `json:"website,omitempty" example:"https://www.expressdelivery.com"`

	// Detalles del contrato en formato JSON
	ContractDetails ContractDetailsDTO `json:"contract_details"`

	// Tarifa base de entrega
	DeliveryRate float64 `json:"delivery_rate" example:"20.50" binding:"required"`

	// URL del logo de la empresa (opcional)
	LogoURL string `json:"logo_url,omitempty" example:"https://www.example.com/logo.png"`

	// Fecha de inicio del contrato
	ContractStartDate time.Time `json:"contract_start_date" binding:"required" format:"date-time"`

	// Fecha de fin del contrato (opcional)
	ContractEndDate *time.Time `json:"contract_end_date,omitempty" format:"date-time"`

	// Dirección principal de la empresa
	MainAddress CompanyAddressDTO `json:"main_address" binding:"required"`
}

// CompanyAddressDTO representa la dirección de una empresa
type CompanyAddressDTO struct {
	// Primera línea de dirección
	AddressLine1 string `json:"address_line1" example:"Calle 100 #15-20" binding:"required"`

	// Segunda línea de dirección (opcional)
	AddressLine2 string `json:"address_line2,omitempty" example:"Edificio Centro Empresarial"`

	// Ciudad
	City string `json:"city" example:"Bogotá" binding:"required"`

	// Estado o provincia
	State string `json:"state" example:"Cundinamarca" binding:"required"`

	// Código postal
	PostalCode string `json:"postal_code,omitempty" example:"110121"`

	// Latitud de la ubicación
	Latitude float64 `json:"latitude" example:"4.68"`

	// Longitud de la ubicación
	Longitude float64 `json:"longitude" example:"-74.05"`

	// Indica si es la dirección principal
	IsMain bool `json:"is_main" example:"true"`
}

// ContractDetailsDTO representa los detalles del contrato
type ContractDetailsDTO struct {
	// Tipo de contrato
	ContractType string `json:"contract_type" example:"Standard" binding:"required"`

	// Términos de pago
	PaymentTerms string `json:"payment_terms" example:"Net 30" binding:"required"`

	// Tipo de renovación
	RenewalType string `json:"renewal_type" example:"Automatic" binding:"required"`

	// Período de notificación en días
	NoticePeriod int `json:"notice_period" example:"30" binding:"required"`

	// Cláusulas especiales (opcional)
	SpecialClauses []string `json:"special_clauses,omitempty"`

	// Persona que firmó el contrato
	SignedBy string `json:"signed_by,omitempty" example:"John Doe"`

	// Fecha de firma del contrato
	SignedAt *time.Time `json:"signed_at,omitempty" format:"date-time"`
}

// CompanyUpdateRequest representa la solicitud para actualizar una empresa existente
type CompanyUpdateRequest struct {
	// Nombre comercial de la empresa
	Name string `json:"name,omitempty" example:"Express Delivery Co."`

	// Nombre legal completo de la empresa
	LegalName string `json:"legal_name,omitempty" example:"Express Delivery S.A.S"`

	// Email de contacto principal de la empresa
	ContactEmail string `json:"contact_email,omitempty" example:"contacto@expressdelivery.com" binding:"omitempty,email"`

	// Teléfono de contacto principal
	ContactPhone string `json:"contact_phone,omitempty" example:"+573001234567"`

	// Sitio web de la empresa
	Website string `json:"website,omitempty" example:"https://www.expressdelivery.com"`

	// Detalles del contrato en formato JSON
	ContractDetails *ContractDetailsDTO `json:"contract_details,omitempty"`

	// Tarifa base de entrega
	DeliveryRate *float64 `json:"delivery_rate,omitempty" example:"20.50"`

	// URL del logo de la empresa
	LogoURL string `json:"logo_url,omitempty" example:"https://www.example.com/logo.png"`

	// Fecha de fin del contrato
	ContractEndDate *time.Time `json:"contract_end_date,omitempty" format:"date-time"`
}

// CompanyResponse representa la respuesta con información de una empresa
type CompanyResponse struct {
	// ID único de la empresa
	ID string `json:"id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`

	// Nombre comercial de la empresa
	Name string `json:"name" example:"Express Delivery Co."`

	// Nombre legal completo de la empresa
	LegalName string `json:"legal_name" example:"Express Delivery S.A.S"`

	// Identificador fiscal de la empresa
	TaxID string `json:"tax_id" example:"900123456-7"`

	// Email de contacto principal de la empresa
	ContactEmail string `json:"contact_email" example:"contacto@expressdelivery.com"`

	// Teléfono de contacto principal
	ContactPhone string `json:"contact_phone" example:"+573001234567"`

	// Sitio web de la empresa
	Website string `json:"website,omitempty" example:"https://www.expressdelivery.com"`

	// Indica si la empresa está activa
	IsActive bool `json:"is_active" example:"true"`

	// Detalles del contrato en formato JSON
	ContractDetails string `json:"contract_details,omitempty"`

	// Tarifa base de entrega
	DeliveryRate float64 `json:"delivery_rate" example:"20.50"`

	// URL del logo de la empresa
	LogoURL string `json:"logo_url,omitempty" example:"https://www.example.com/logo.png"`

	// Fecha de inicio del contrato
	ContractStartDate time.Time `json:"contract_start_date" format:"date-time"`

	// Fecha de fin del contrato
	ContractEndDate *time.Time `json:"contract_end_date,omitempty" format:"date-time"`

	// Cuando se creó la empresa
	CreatedAt time.Time `json:"created_at" format:"date-time"`

	// Última actualización de la empresa
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`

	// Direcciones asociadas a la empresa
	Addresses []CompanyAddressResponse `json:"addresses,omitempty"`

	// Sucursales asociadas a la empresa
	Branches []BranchResponse `json:"branches,omitempty"`

	// Métricas de la empresa (opcional)
	Metrics *CompanyMetricsResponse `json:"metrics,omitempty"`
}

// CompanyAddressResponse representa la respuesta con información de una dirección de empresa
type CompanyAddressResponse struct {
	// ID único de la dirección
	ID string `json:"id" example:"e1b09d38-e71f-415f-b3eb-ffeb8dd3b493"`

	// ID de la empresa a la que pertenece
	CompanyID string `json:"company_id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`

	// Primera línea de dirección
	AddressLine1 string `json:"address_line1" example:"Calle 100 #15-20"`

	// Segunda línea de dirección
	AddressLine2 string `json:"address_line2,omitempty" example:"Edificio Centro Empresarial"`

	// Ciudad
	City string `json:"city" example:"Bogotá"`

	// Estado o provincia
	State string `json:"state" example:"Cundinamarca"`

	// Código postal
	PostalCode string `json:"postal_code,omitempty" example:"110121"`

	// Latitud de la ubicación
	Latitude float64 `json:"latitude" example:"4.68"`

	// Longitud de la ubicación
	Longitude float64 `json:"longitude" example:"-74.05"`

	// Indica si es la dirección principal
	IsMain bool `json:"is_main" example:"true"`

	// Cuando se creó la dirección
	CreatedAt time.Time `json:"created_at" format:"date-time"`
}

// CompanyActivationRequest representa la solicitud para activar o desactivar una empresa
type CompanyActivationRequest struct {
	// Indica si la empresa debe estar activa
	IsActive bool `json:"is_active" example:"true" binding:"required"`

	// Razón del cambio de estado
	Reason string `json:"reason" example:"Contrato vencido" binding:"required"`
}

// CompanyMetricsResponse representa las métricas de una empresa
type CompanyMetricsResponse struct {
	// Total de órdenes procesadas
	TotalOrders int64 `json:"total_orders" example:"1250"`

	// Total de órdenes completadas exitosamente
	CompletedOrders int64 `json:"completed_orders" example:"1100"`

	// Total de órdenes canceladas
	CancelledOrders int64 `json:"cancelled_orders" example:"50"`

	// Tasa de entrega exitosa (porcentaje)
	DeliverySuccessRate float64 `json:"delivery_success_rate" example:"88.0"`

	// Tiempo promedio de entrega (en minutos)
	AverageDeliveryTime float64 `json:"average_delivery_time" example:"45.5"`

	// Ingresos totales generados
	TotalRevenue float64 `json:"total_revenue" example:"25000.50"`

	// Número de sucursales activas
	ActiveBranches int `json:"active_branches" example:"5"`

	// Número de clientes únicos
	UniqueCustomers int `json:"unique_customers" example:"500"`
}

// CompanyListResponse representa un item en la lista de empresas
type CompanyListResponse struct {
	// ID único de la empresa
	ID string `json:"id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`

	// Nombre comercial de la empresa
	Name string `json:"name" example:"Express Delivery Co."`

	// Nombre legal completo de la empresa
	LegalName string `json:"legal_name" example:"Express Delivery S.A.S"`

	// Identificador fiscal de la empresa
	TaxID string `json:"tax_id" example:"900123456-7"`

	// Email de contacto principal de la empresa
	ContactEmail string `json:"contact_email" example:"contacto@expressdelivery.com"`

	// Teléfono de contacto principal
	ContactPhone string `json:"contact_phone" example:"+573001234567"`

	// Indica si la empresa está activa
	IsActive bool `json:"is_active" example:"true"`

	// Número de sucursales
	BranchCount int `json:"branch_count" example:"3"`

	// Fecha de inicio del contrato
	ContractStartDate time.Time `json:"contract_start_date" format:"date-time"`

	// Fecha de fin del contrato
	ContractEndDate *time.Time `json:"contract_end_date,omitempty" format:"date-time"`

	// Cuando se creó la empresa
	CreatedAt time.Time `json:"created_at" format:"date-time"`
}

// CompanySimpleListResponse representa un item simplificado en la lista de empresas
type CompanySimpleListResponse struct {
	// ID único de la empresa
	ID string `json:"id" example:"b5f8c3d1-2e59-4c4b-a6e8-e5f3c0c3d1b5"`

	// Nombre comercial de la empresa
	Name string `json:"name" example:"Express Delivery Co."`

	// Indica si la empresa está activa
	IsActive bool `json:"is_active" example:"true"`

	// Número de sucursales
	Count int `json:"count" example:"12"`
}
