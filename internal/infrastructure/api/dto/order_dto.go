package dto

import "time"

// OrderCreateRequest represents the request body for creating a new order
// @Description Request structure for creating a delivery order
type OrderCreateRequest struct {
	// Unique identifier of the company
	// @example a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11
	// @required
	CompanyID string `json:"company_id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11" binding:"required"`

	// Unique identifier of the company branch
	// @example b5f8c3d1-2e59-4c4b-a6e8-e5f3c0c3d1b5
	// @required
	BranchID string `json:"branch_id" example:"b5f8c3d1-2e59-4c4b-a6e8-e5f3c0c3d1b5" binding:"required"`

	// Unique identifier of the client
	// @example c7d8e9f0-3f4a-5c6b-7d8e-9f0a1b2c3d4e
	// @required
	ClientID string `json:"client_id" example:"c7d8e9f0-3f4a-5c6b-7d8e-9f0a1b2c3d4e" binding:"required"`

	// Price of the delivery in local currency
	// @minimum 0
	// @example 25.50
	// @required
	Price float64 `json:"price" example:"25.50" binding:"required,min=0"`

	// Distance of the delivery in kilometers
	// @minimum 0
	// @example 7.2
	// @required
	Distance float64 `json:"distance" example:"7.2" binding:"required,min=0"`

	// Scheduled pickup time
	// @example 2023-05-15T14:30:00Z
	// @required
	PickupTime time.Time `json:"pickup_time" example:"2023-05-15T14:30:00Z" binding:"required" format:"date-time"`

	// Deadline for delivery completion
	// @example 2023-05-15T16:30:00Z
	// @required
	DeliveryDeadline time.Time `json:"delivery_deadline" example:"2023-05-15T16:30:00Z" binding:"required" format:"date-time"`

	// Whether recipient signature is required for delivery
	// @example false
	RequiresSignature bool `json:"requires_signature" example:"false"`

	// Additional notes for the delivery
	// @example "Please call recipient 5 minutes before arrival"
	DeliveryNotes string `json:"delivery_notes,omitempty" example:"Please call recipient 5 minutes before arrival"`

	// Details about the package being delivered
	// @required
	PackageDetails PackageDetailRequest `json:"package_details" binding:"required"`

	// Delivery destination address details
	// @required
	DeliveryAddress DeliveryAddressRequest `json:"delivery_address" binding:"required"`

	// Pickup origin address details
	// @required
	PickupAddress PickupAddressRequest `json:"pickup_address" binding:"required"`
}

// PackageDetailRequest contains details about the package
// @Description Package characteristics and handling information
type PackageDetailRequest struct {
	// Whether the package contains fragile items
	// @example true
	IsFragile bool `json:"is_fragile" example:"true"`

	// Whether the package requires urgent handling
	// @example false
	IsUrgent bool `json:"is_urgent" example:"false"`

	// Weight of the package in kilograms
	// @minimum 0
	// @example 2.5
	Weight float64 `json:"weight,omitempty" example:"2.5" binding:"omitempty,min=0"`

	// Special handling instructions
	// @example "Contains glass items, handle with care"
	SpecialInstructions string `json:"special_instructions,omitempty" example:"Contains glass items, handle with care"`

	// Length of the package in centimeters
	// @minimum 0
	// @example 30
	Length float64 `json:"length,omitempty" example:"30" binding:"omitempty,min=0"`

	// Width of the package in centimeters
	// @minimum 0
	// @example 20
	Width float64 `json:"width,omitempty" example:"20" binding:"omitempty,min=0"`

	// Height of the package in centimeters
	// @minimum 0
	// @example 15
	Height float64 `json:"height,omitempty" example:"15" binding:"omitempty,min=0"`
}

// DeliveryAddressRequest contains the destination address details
// @Description Delivery destination address information
type DeliveryAddressRequest struct {
	// Name of the person receiving the package
	// @example "John Doe"
	// @required
	RecipientName string `json:"recipient_name" example:"John Doe" binding:"required"`

	// Contact phone number of the recipient
	// @example "+1234567890"
	// @required
	RecipientPhone string `json:"recipient_phone" example:"+1234567890" binding:"required"`

	// First line of the address
	// @example "123 Main Street"
	// @required
	AddressLine1 string `json:"address_line1" example:"123 Main Street" binding:"required"`

	// Second line of the address (optional)
	// @example "Apartment 4B"
	AddressLine2 string `json:"address_line2,omitempty" example:"Apartment 4B"`

	// City name
	// @example "New York"
	// @required
	City string `json:"city" example:"New York" binding:"required"`

	// State or province name
	// @example "NY"
	// @required
	State string `json:"state" example:"NY" binding:"required"`

	// Postal or ZIP code
	// @example "10001"
	PostalCode string `json:"postal_code,omitempty" example:"10001"`

	// Latitude coordinate
	// @minimum -90
	// @maximum 90
	// @example 40.7128
	// @required
	Latitude float64 `json:"latitude" example:"40.7128" binding:"required,min=-90,max=90"`

	// Longitude coordinate
	// @minimum -180
	// @maximum 180
	// @example -74.0060
	// @required
	Longitude float64 `json:"longitude" example:"-74.0060" binding:"required,min=-180,max=180"`

	// Additional notes about the address
	// @example "Ring doorbell twice"
	AddressNotes string `json:"address_notes,omitempty" example:"Ring doorbell twice"`
}

// PickupAddressRequest contains the origin address details
// @Description Pickup origin address information
type PickupAddressRequest struct {
	// Name of the contact person at pickup location
	// @example "Jane Smith"
	// @required
	ContactName string `json:"contact_name" example:"Jane Smith" binding:"required"`

	// Contact phone number at pickup location
	// @example "+0987654321"
	// @required
	ContactPhone string `json:"contact_phone" example:"+0987654321" binding:"required"`

	// First line of the address
	// @example "456 Business Ave"
	// @required
	AddressLine1 string `json:"address_line1" example:"456 Business Ave" binding:"required"`

	// Second line of the address (optional)
	// @example "Suite 300"
	AddressLine2 string `json:"address_line2,omitempty" example:"Suite 300"`

	// City name
	// @example "Chicago"
	// @required
	City string `json:"city" example:"Chicago" binding:"required"`

	// State or province name
	// @example "IL"
	// @required
	State string `json:"state" example:"IL" binding:"required"`

	// Postal or ZIP code
	// @example "60606"
	PostalCode string `json:"postal_code,omitempty" example:"60606"`

	// Latitude coordinate
	// @minimum -90
	// @maximum 90
	// @example 41.8781
	// @required
	Latitude float64 `json:"latitude" example:"41.8781" binding:"required,min=-90,max=90"`

	// Longitude coordinate
	// @minimum -180
	// @maximum 180
	// @example -87.6298
	// @required
	Longitude float64 `json:"longitude" example:"-87.6298" binding:"required,min=-180,max=180"`

	// Additional notes about the address
	// @example "Enter through loading dock"
	AddressNotes string `json:"address_notes,omitempty" example:"Enter through loading dock"`
}
