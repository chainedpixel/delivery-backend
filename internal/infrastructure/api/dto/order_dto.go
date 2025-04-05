package dto

import (
	domainErr "github.com/MarlonG1/delivery-backend/internal/domain/error"
	infraErr "github.com/MarlonG1/delivery-backend/internal/infrastructure/error"
	"time"
)

// OrderCreateRequest represents the request body for creating a new order
// @Description Request structure for creating a delivery order
type OrderCreateRequest struct {
	// Unique identifier of the company pickup location
	// @required
	CompanyPickUpID string `json:"company_pickup_id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11" binding:"required"`

	// Unique identifier of the client
	// @required
	ClientID string `json:"client_id" example:"c7d8e9f0-3f4a-5c6b-7d8e-9f0a1b2c3d4e" binding:"required"`

	// Price of the delivery in local currency
	// @minimum 0
	// @required
	Price float64 `json:"price" example:"25.50" binding:"required,min=0"`

	// Distance of the delivery in kilometers
	// @minimum 0
	// @required
	Distance float64 `json:"distance" example:"7.2" binding:"required,min=0"`

	// Scheduled pickup time
	// @required
	PickupTime time.Time `json:"pickup_time" example:"2023-05-15T14:30:00Z" binding:"required" format:"date-time"`

	// Deadline for delivery completion
	// @required
	DeliveryDeadline time.Time `json:"delivery_deadline" example:"2023-05-15T16:30:00Z" binding:"required" format:"date-time"`

	// Whether recipient signature is required for delivery
	RequiresSignature bool `json:"requires_signature" example:"false"`

	// Additional notes for the delivery
	DeliveryNotes string `json:"delivery_notes,omitempty" example:"Please call recipient 5 minutes before arrival"`

	// Details about the package being delivered
	// @required
	PackageDetails PackageDetailRequest `json:"package_details" binding:"required"`

	// Contact name for pickup location
	// @required
	PickupContactName string `json:"pickup_contact_name" binding:"required"`

	// Contact phone number for pickup location
	// @required
	PickupContactPhone string `json:"pickup_contact_phone" binding:"required"`

	// Important notes about the pickup location
	// @required
	PickupNotes string `json:"pickup_notes,omitempty"`

	// Delivery destination address details
	// @required
	DeliveryAddress DeliveryAddressRequest `json:"delivery_address" binding:"required"`
}

func (o *OrderCreateRequest) Validate() error {
	if o.CompanyPickUpID == "" {
		return infraErr.NewGeneralServiceError("OrderDTO", "Validate", domainErr.ErrCompanyPickUpIDRequired)
	}

	if o.ClientID == "" {
		return infraErr.NewGeneralServiceError("OrderDTO", "Validate", domainErr.ErrClientIDRequired)
	}

	return nil
}

// PackageDetailRequest contains details about the package
// @Description Package characteristics and handling information
type PackageDetailRequest struct {
	// Whether the package contains fragile items
	IsFragile bool `json:"is_fragile" example:"true"`

	// Whether the package requires urgent handling
	IsUrgent bool `json:"is_urgent" example:"false"`

	// Weight of the package in kilograms
	// @minimum 0
	Weight float64 `json:"weight,omitempty" example:"2.5" binding:"omitempty,min=0"`

	// Special handling instructions
	SpecialInstructions string `json:"special_instructions,omitempty" example:"Contains glass items, handle with care"`

	// Length of the package in centimeters
	// @minimum 0
	Length float64 `json:"length,omitempty" example:"30" binding:"omitempty,min=0"`

	// Width of the package in centimeters
	// @minimum 0
	Width float64 `json:"width,omitempty" example:"20" binding:"omitempty,min=0"`

	// Height of the package in centimeters
	// @minimum 0
	Height float64 `json:"height,omitempty" example:"15" binding:"omitempty,min=0"`
}

// DeliveryAddressRequest contains the destination address details
// @Description Delivery destination address information
type DeliveryAddressRequest struct {
	// Name of the person receiving the package
	// @required
	RecipientName string `json:"recipient_name" example:"John Doe" binding:"required"`

	// Contact phone number of the recipient
	// @required
	RecipientPhone string `json:"recipient_phone" example:"+1234567890" binding:"required"`

	// First line of the address
	// @required
	AddressLine1 string `json:"address_line1" example:"123 Main Street" binding:"required"`

	// Second line of the address (optional)
	AddressLine2 string `json:"address_line2,omitempty" example:"Apartment 4B"`

	// City name
	// @required
	City string `json:"city" example:"New York" binding:"required"`

	// State or province name
	// @required
	State string `json:"state" example:"NY" binding:"required"`

	// Postal or ZIP code
	PostalCode string `json:"postal_code,omitempty" example:"10001"`

	// TODO: Temporalmente desactivados hasta implementar correctamente el manejo geoespacial

	// Additional notes about the address
	AddressNotes string `json:"address_notes,omitempty" example:"Ring doorbell twice"`
}

// OrderResponse represents the response for an order
// @Description Order information with all related details
type OrderResponse struct {
	// Unique identifier of the order
	ID string `json:"id" example:"a1b2c3d4-e5f6-7g8h-9i0j-k1l2m3n4o5p6"`

	// Company ID that owns the order
	CompanyID string `json:"company_id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`

	// Company name
	CompanyName string `json:"company_name,omitempty" example:"Express Delivery Inc."`

	// Branch ID where the order originated
	BranchID string `json:"branch_id" example:"b5f8c3d1-2e59-4c4b-a6e8-e5f3c0c3d1b5"`

	// Branch name
	BranchName string `json:"branch_name,omitempty" example:"Downtown Branch"`

	// Client ID who placed the order
	ClientID string `json:"client_id" example:"c7d8e9f0-3f4a-5c6b-7d8e-9f0a1b2c3d4e"`

	// Client full name
	ClientName string `json:"client_name,omitempty" example:"John Smith"`

	// Driver ID assigned to the order
	DriverID *string `json:"driver_id,omitempty" example:"d1e2f3g4-h5i6-j7k8-l9m0-n1o2p3q4r5s6"`

	// Driver full name
	DriverName string `json:"driver_name,omitempty" example:"Michael Johnson"`

	// Tracking number for the order
	TrackingNumber string `json:"tracking_number" example:"DEL-230512-7890"`

	// Current status of the order
	Status string `json:"status" example:"PENDING"`

	// When the order was created
	CreatedAt time.Time `json:"created_at" example:"2023-05-15T10:30:00Z" format:"date-time"`

	// When the order was last updated
	UpdatedAt time.Time `json:"updated_at" example:"2023-05-15T10:30:00Z" format:"date-time"`

	// Details of the order
	Detail OrderDetailResponse `json:"detail"`

	// Details about the package
	PackageDetail PackageDetailResponse `json:"package_detail"`

	// Delivery destination address
	DeliveryAddress DeliveryAddressResponse `json:"delivery_address"`

	// Pickup origin address
	PickupAddress PickupAddressResponse `json:"pickup_address"`

	// Order status history
	StatusHistory []OrderStatusHistoryResponse `json:"status_history"`

	// Current tracking status
	CurrentStatus string `json:"current_status" example:"IN_TRANSIT"`

	// Last time the order status was updated
	LastUpdated time.Time `json:"last_updated,omitempty" example:"2023-05-15T12:45:00Z" format:"date-time"`

	// Estimated time of arrival
	EstimatedArrival time.Time `json:"estimated_arrival,omitempty" example:"2023-05-15T16:30:00Z" format:"date-time"`
}

// OrderDetailResponse contains detailed information about the order
// @Description Order details including price, schedule and delivery requirements
type OrderDetailResponse struct {
	// Price of the delivery
	Price float64 `json:"price" example:"25.50"`

	// Distance to be traveled in kilometers
	Distance float64 `json:"distance" example:"7.2"`

	// Scheduled pickup time
	PickupTime time.Time `json:"pickup_time" example:"2023-05-15T14:30:00Z" format:"date-time"`

	// Deadline for delivery
	DeliveryDeadline time.Time `json:"delivery_deadline" example:"2023-05-15T16:30:00Z" format:"date-time"`

	// When the order was actually delivered
	DeliveredAt *time.Time `json:"delivered_at,omitempty" example:"2023-05-15T16:15:00Z" format:"date-time"`

	// Whether recipient signature is required
	RequiresSignature bool `json:"requires_signature" example:"false"`

	// Additional notes for delivery
	DeliveryNotes string `json:"delivery_notes,omitempty" example:"Please call recipient 5 minutes before arrival"`
}

// PackageDetailResponse contains information about the package
// @Description Package characteristics and handling requirements
type PackageDetailResponse struct {
	// Whether the package contains fragile items
	IsFragile bool `json:"is_fragile" example:"true"`

	// Whether the package requires urgent handling
	IsUrgent bool `json:"is_urgent" example:"false"`

	// Weight of the package in kilograms
	Weight float64 `json:"weight,omitempty" example:"2.5"`

	// Package dimensions in JSON format
	Dimensions string `json:"dimensions,omitempty" example:"{\"length\":30,\"width\":20,\"height\":15,\"unit\":\"cm\"}"`

	// Special handling instructions
	SpecialInstructions string `json:"special_instructions,omitempty" example:"Contains glass items, handle with care"`
}

type OrderStatusHistoryResponse struct {
	// Name of the status
	Status string `json:"status" example:"PENDING"`
	// Description of the status change
	Description string `json:"description,omitempty" example:"Driver has accepted the order and is heading to pickup location"`
	// Updated at time
	UpdatedAt string `json:"updated_at" example:"2023-05-15T12:45:00Z" format:"date-time"`
}

// DeliveryAddressResponse contains the destination address details
// @Description Delivery address information
type DeliveryAddressResponse struct {
	// Name of the recipient
	RecipientName string `json:"recipient_name" example:"John Doe"`

	// Contact phone number of the recipient
	RecipientPhone string `json:"recipient_phone" example:"+1234567890"`

	// First line of the address
	AddressLine1 string `json:"address_line1" example:"123 Main Street"`

	// Second line of the address (optional)
	AddressLine2 string `json:"address_line2,omitempty" example:"Apartment 4B"`

	// City name
	City string `json:"city" example:"New York"`

	// State or province name
	State string `json:"state" example:"NY"`

	// Postal or ZIP code
	PostalCode string `json:"postal_code,omitempty" example:"10001"`

	// Latitude coordinate
	Latitude float64 `json:"latitude" example:"40.7128"`

	// Longitude coordinate
	Longitude float64 `json:"longitude" example:"-74.0060"`
}

// PickupAddressResponse contains the origin address details
// @Description Pickup address information
type PickupAddressResponse struct {
	// Name of the contact person at pickup location
	ContactName string `json:"contact_name" example:"Jane Smith"`

	// Contact phone number at pickup location
	ContactPhone string `json:"contact_phone" example:"+0987654321"`

	// First line of the address
	AddressLine1 string `json:"address_line1" example:"456 Business Ave"`

	// Second line of the address (optional)
	AddressLine2 string `json:"address_line2,omitempty" example:"Suite 300"`

	// City name
	City string `json:"city" example:"Chicago"`

	// State or province name
	State string `json:"state" example:"IL"`

	// Postal or ZIP code
	PostalCode string `json:"postal_code,omitempty" example:"60606"`

	// Latitude coordinate
	Latitude float64 `json:"latitude" example:"41.8781"`

	// Longitude coordinate
	Longitude float64 `json:"longitude" example:"-87.6298"`

	// Additional notes about the address
	AddressNotes string `json:"address_notes,omitempty" example:"Enter through loading dock"`

	// Full formatted address
	FormattedAddress string `json:"formatted_address,omitempty" example:"456 Business Ave, Suite 300, Chicago, IL 60606"`
}

// OrderListResponse is a simplified order representation for listings
// @Description Simplified order information for list views
type OrderListResponse struct {
	// Unique identifier of the order
	ID string `json:"id" example:"a1b2c3d4-e5f6-7g8h-9i0j-k1l2m3n4o5p6"`

	// Tracking number for the order
	TrackingNumber string `json:"tracking_number" example:"DEL-230512-7890"`

	// Full name of the client
	ClientName string `json:"client_name" example:"John Smith"`

	// Shortened delivery address for display
	DeliveryAddress string `json:"delivery_address" example:"123 Main St, New York, NY"`

	// Deadline for delivery
	DeliveryDeadline time.Time `json:"delivery_deadline" example:"2023-05-15T16:30:00Z" format:"date-time"`

	// Price of the delivery
	Price float64 `json:"price" example:"25.50"`

	// Current status of the order
	// @enum [PENDING,ACCEPTED,PICKED_UP,IN_TRANSIT,DELIVERED,CANCELLED]
	Status string `json:"status" example:"PENDING" enums:"PENDING,ACCEPTED,PICKED_UP,IN_TRANSIT,DELIVERED,CANCELLED"`

	// Driver ID assigned to the order (if any)
	DriverID *string `json:"driver_id,omitempty" example:"d1e2f3g4-h5i6-j7k8-l9m0-n1o2p3q4r5s6"`

	// Full name of the assigned driver
	DriverName string `json:"driver_name,omitempty" example:"Michael Johnson"`

	// When the order was created
	CreatedAt time.Time `json:"created_at" example:"2023-05-15T10:30:00Z" format:"date-time"`
}

type UserListResponse struct {
	// Unique identifier of the user
	ID string `json:"id" example:"a1b2c3d4-e5f6-7g8h-9i0j-k1l2m3n4o5p6"`

	// Full name of the user
	FullName string `json:"full_name" example:"John Smith"`

	// Phone number of the user
	Phone string `json:"phone" example:"21212828"`

	// Email address of the user
	Email string `json:"email" example:"example@example.com"`

	// Role of the user
	Role string `json:"role" example:"Admin"`

	// Document type of the user
	DocumentType string `json:"document_type" example:"DUI"`

	// Document number of the user
	DocumentNumber string `json:"document_number" example:"1234567890"`

	// Created at time
	CreatedAt time.Time `json:"created_at" example:"2023-05-15T10:30:00Z" format:"date-time"`
}

// OrderDriverAssignRequest represents the request to assign a driver to an order
// @Description Request to assign a driver to an order
type OrderDriverAssignRequest struct {
	// ID of the driver to assign
	// @required
	DriverID string `json:"driver_id" example:"d1e2f3g4-h5i6-j7k8-l9m0-n1o2p3q4r5s6" binding:"required,uuid"`
}

// OrderUpdateRequest represents the request body for updating an existing order
// @Description Request structure for updating a delivery order
type OrderUpdateRequest struct {
	// Price of the delivery in local currency
	// @minimum 0
	Price float64 `json:"price,omitempty" example:"25.50" binding:"omitempty,min=0"`

	// Distance of the delivery in kilometers
	// @minimum 0
	Distance float64 `json:"distance,omitempty" example:"7.2" binding:"omitempty,min=0"`

	// Scheduled pickup time - only modifiable if order is still in PENDING state
	PickupTime *time.Time `json:"pickup_time,omitempty" example:"2023-05-15T14:30:00Z" format:"date-time"`

	// Deadline for delivery completion - only modifiable if order is still in PENDING state
	DeliveryDeadline *time.Time `json:"delivery_deadline,omitempty" example:"2023-05-15T16:30:00Z" format:"date-time"`

	// Whether recipient signature is required for delivery
	RequiresSignature *bool `json:"requires_signature,omitempty" example:"false"`

	// Additional notes for the delivery
	DeliveryNotes string `json:"delivery_notes,omitempty" example:"Please call recipient 5 minutes before arrival"`

	// Details about the package being delivered
	PackageDetails *PackageDetailUpdateRequest `json:"package_details,omitempty"`

	// Delivery destination address details
	DeliveryAddress *DeliveryAddressUpdateRequest `json:"delivery_address,omitempty"`

	// Pickup contact information and notes
	// Contact name for pickup location
	PickupContactName string `json:"pickup_contact_name,omitempty" example:"Jane Smith"`

	// Contact phone number for pickup location
	PickupContactPhone string `json:"pickup_contact_phone,omitempty" example:"+0987654321"`

	// Important notes about the pickup location
	PickupNotes string `json:"pickup_notes,omitempty" example:"Enter through loading dock"`
}

func (o *OrderUpdateRequest) Validate() error {
	// Time validations - ensure delivery deadline is after pickup time if both are provided
	if o.PickupTime != nil && o.DeliveryDeadline != nil {
		if o.DeliveryDeadline.Before(*o.PickupTime) {
			return infraErr.NewGeneralServiceError("OrderUpdateDTO", "Validate", domainErr.ErrDeliveryDeadlineBeforePickup)
		}
	}

	return nil
}

// PackageDetailUpdateRequest contains details about the package for updates
// @Description Package characteristics and handling information for updates
type PackageDetailUpdateRequest struct {
	// Whether the package contains fragile items
	IsFragile *bool `json:"is_fragile,omitempty" example:"true"`

	// Whether the package requires urgent handling
	IsUrgent *bool `json:"is_urgent,omitempty" example:"false"`

	// Weight of the package in kilograms
	// @minimum 0
	Weight *float64 `json:"weight,omitempty" example:"2.5" binding:"omitempty,min=0"`

	// Special handling instructions
	SpecialInstructions string `json:"special_instructions,omitempty" example:"Contains glass items, handle with care"`

	// Length of the package in centimeters
	// @minimum 0
	Length *float64 `json:"length,omitempty" example:"30" binding:"omitempty,min=0"`

	// Width of the package in centimeters
	// @minimum 0
	Width *float64 `json:"width,omitempty" example:"20" binding:"omitempty,min=0"`

	// Height of the package in centimeters
	// @minimum 0
	Height *float64 `json:"height,omitempty" example:"15" binding:"omitempty,min=0"`
}

// DeliveryAddressUpdateRequest contains the destination address details for updates
// @Description Delivery destination address information for updates
type DeliveryAddressUpdateRequest struct {
	// Name of the person receiving the package
	RecipientName string `json:"recipient_name,omitempty" example:"John Doe"`

	// Contact phone number of the recipient
	RecipientPhone string `json:"recipient_phone,omitempty" example:"+1234567890"`

	// First line of the address
	AddressLine1 string `json:"address_line1,omitempty" example:"123 Main Street"`

	// Second line of the address (optional)
	AddressLine2 string `json:"address_line2,omitempty" example:"Apartment 4B"`

	// City name
	City string `json:"city,omitempty" example:"New York"`

	// State or province name
	State string `json:"state,omitempty" example:"NY"`

	// Postal or ZIP code
	PostalCode string `json:"postal_code,omitempty" example:"10001"`

	// Additional notes about the address
	AddressNotes string `json:"address_notes,omitempty" example:"Ring doorbell twice"`
}

// OrderStatusUpdateRequest represents the request to update an order's status
// @Description Request to change the status of an order
type OrderStatusUpdateRequest struct {
	// New status for the order
	// @enum [PENDING,ACCEPTED,PICKED_UP,IN_TRANSIT,DELIVERED,CANCELLED]
	// @required
	Status string `json:"status" example:"ACCEPTED" binding:"required" enums:"PENDING,ACCEPTED,PICKED_UP,IN_TRANSIT,DELIVERED,CANCELLED"`

	// Optional description about the status change
	Description string `json:"description,omitempty" example:"Driver has accepted the order and is heading to pickup location"`
}
