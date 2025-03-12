package dto

import (
	errPackage "domain/error"
	error2 "infrastructure/error"
	"time"
)

// OrderCreateRequest represents the request body for creating a new order
// @Description Request structure for creating a delivery order
type OrderCreateRequest struct {
	// Unique identifier of the company pickup location
	// @example a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11
	// @required
	CompanyPickUpID string `json:"company_pickup_id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11" binding:"required"`

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

	// Contact name for pickup location
	// @example "Jane Smith"
	// @required
	PickupContactName string `json:"pickup_contact_name" binding:"required"`

	// Contact phone number for pickup location
	// @example "+0987654321"
	// @required
	PickupContactPhone string `json:"pickup_contact_phone" binding:"required"`

	// Important notes about the pickup location
	// @example "Enter through loading dock"
	// @required
	PickupNotes string `json:"pickup_notes,omitempty"`

	// Delivery destination address details
	// @required
	DeliveryAddress DeliveryAddressRequest `json:"delivery_address" binding:"required"`
}

func (o *OrderCreateRequest) Validate() error {
	if o.CompanyPickUpID == "" {
		return error2.NewGeneralServiceError("OrderDTO", "Validate", errPackage.ErrCompanyPickUpIDRequired)
	}

	if o.ClientID == "" {
		return error2.NewGeneralServiceError("OrderDTO", "Validate", errPackage.ErrClientIDRequired)
	}

	return nil
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

	// TODO: Temporalmente desactivados hasta implementar correctamente el manejo geoespacial

	// Additional notes about the address
	// @example "Ring doorbell twice"
	AddressNotes string `json:"address_notes,omitempty" example:"Ring doorbell twice"`
}

// OrderResponse represents the response for an order
// @Description Order information with all related details
type OrderResponse struct {
	// Unique identifier of the order
	// @example a1b2c3d4-e5f6-7g8h-9i0j-k1l2m3n4o5p6
	ID string `json:"id" example:"a1b2c3d4-e5f6-7g8h-9i0j-k1l2m3n4o5p6"`

	// Company ID that owns the order
	// @example a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11
	CompanyID string `json:"company_id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`

	// Company name
	// @example "Express Delivery Inc."
	CompanyName string `json:"company_name,omitempty" example:"Express Delivery Inc."`

	// Branch ID where the order originated
	// @example b5f8c3d1-2e59-4c4b-a6e8-e5f3c0c3d1b5
	BranchID string `json:"branch_id" example:"b5f8c3d1-2e59-4c4b-a6e8-e5f3c0c3d1b5"`

	// Branch name
	// @example "Downtown Branch"
	BranchName string `json:"branch_name,omitempty" example:"Downtown Branch"`

	// Client ID who placed the order
	// @example c7d8e9f0-3f4a-5c6b-7d8e-9f0a1b2c3d4e
	ClientID string `json:"client_id" example:"c7d8e9f0-3f4a-5c6b-7d8e-9f0a1b2c3d4e"`

	// Client full name
	// @example "John Smith"
	ClientName string `json:"client_name,omitempty" example:"John Smith"`

	// Driver ID assigned to the order
	// @example d1e2f3g4-h5i6-j7k8-l9m0-n1o2p3q4r5s6
	DriverID *string `json:"driver_id,omitempty" example:"d1e2f3g4-h5i6-j7k8-l9m0-n1o2p3q4r5s6"`

	// Driver full name
	// @example "Michael Johnson"
	DriverName string `json:"driver_name,omitempty" example:"Michael Johnson"`

	// Tracking number for the order
	// @example DEL-230512-7890
	TrackingNumber string `json:"tracking_number" example:"DEL-230512-7890"`

	// Current status of the order
	// @example PENDING
	Status string `json:"status" example:"PENDING"`

	// When the order was created
	// @example 2023-05-15T10:30:00Z
	CreatedAt time.Time `json:"created_at" example:"2023-05-15T10:30:00Z" format:"date-time"`

	// When the order was last updated
	// @example 2023-05-15T10:30:00Z
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
	// @example [{"status":"PENDING","description":"Order has been created","updated_at":"2023-05-15T10:30:00Z"}]
	StatusHistory []OrderStatusHistoryResponse `json:"status_history"`

	// Current tracking status
	// @example IN_TRANSIT
	CurrentStatus string `json:"current_status" example:"IN_TRANSIT"`

	// Last time the order status was updated
	// @example 2023-05-15T12:45:00Z
	LastUpdated time.Time `json:"last_updated,omitempty" example:"2023-05-15T12:45:00Z" format:"date-time"`

	// Estimated time of arrival
	// @example 2023-05-15T16:30:00Z
	EstimatedArrival time.Time `json:"estimated_arrival,omitempty" example:"2023-05-15T16:30:00Z" format:"date-time"`
}

// OrderDetailResponse contains detailed information about the order
// @Description Order details including price, schedule and delivery requirements
type OrderDetailResponse struct {
	// Price of the delivery
	// @example 25.50
	Price float64 `json:"price" example:"25.50"`

	// Distance to be traveled in kilometers
	// @example 7.2
	Distance float64 `json:"distance" example:"7.2"`

	// Scheduled pickup time
	// @example 2023-05-15T14:30:00Z
	PickupTime time.Time `json:"pickup_time" example:"2023-05-15T14:30:00Z" format:"date-time"`

	// Deadline for delivery
	// @example 2023-05-15T16:30:00Z
	DeliveryDeadline time.Time `json:"delivery_deadline" example:"2023-05-15T16:30:00Z" format:"date-time"`

	// When the order was actually delivered
	// @example 2023-05-15T16:15:00Z
	DeliveredAt *time.Time `json:"delivered_at,omitempty" example:"2023-05-15T16:15:00Z" format:"date-time"`

	// Whether recipient signature is required
	// @example false
	RequiresSignature bool `json:"requires_signature" example:"false"`

	// Additional notes for delivery
	// @example "Please call recipient 5 minutes before arrival"
	DeliveryNotes string `json:"delivery_notes,omitempty" example:"Please call recipient 5 minutes before arrival"`
}

// PackageDetailResponse contains information about the package
// @Description Package characteristics and handling requirements
type PackageDetailResponse struct {
	// Whether the package contains fragile items
	// @example true
	IsFragile bool `json:"is_fragile" example:"true"`

	// Whether the package requires urgent handling
	// @example false
	IsUrgent bool `json:"is_urgent" example:"false"`

	// Weight of the package in kilograms
	// @example 2.5
	Weight float64 `json:"weight,omitempty" example:"2.5"`

	// Package dimensions in JSON format
	// @example {"length":30,"width":20,"height":15,"unit":"cm"}
	Dimensions string `json:"dimensions,omitempty" example:"{\"length\":30,\"width\":20,\"height\":15,\"unit\":\"cm\"}"`

	// Special handling instructions
	// @example "Contains glass items, handle with care"
	SpecialInstructions string `json:"special_instructions,omitempty" example:"Contains glass items, handle with care"`
}

type OrderStatusHistoryResponse struct {
	// Name of the status
	// @example ACCEPTED
	Status string `json:"status" example:"PENDING"`
	// Description of the status change
	// @example "Driver has accepted the order and is heading to pickup location"
	Description string `json:"description,omitempty" example:"Driver has accepted the order and is heading to pickup location"`
	// Updated at time
	// @example 2023-05-15T12:45:00Z
	UpdatedAt string `json:"updated_at" example:"2023-05-15T12:45:00Z" format:"date-time"`
}

// DeliveryAddressResponse contains the destination address details
// @Description Delivery address information
type DeliveryAddressResponse struct {
	// Name of the recipient
	// @example "John Doe"
	RecipientName string `json:"recipient_name" example:"John Doe"`

	// Contact phone number of the recipient
	// @example "+1234567890"
	RecipientPhone string `json:"recipient_phone" example:"+1234567890"`

	// First line of the address
	// @example "123 Main Street"
	AddressLine1 string `json:"address_line1" example:"123 Main Street"`

	// Second line of the address (optional)
	// @example "Apartment 4B"
	AddressLine2 string `json:"address_line2,omitempty" example:"Apartment 4B"`

	// City name
	// @example "New York"
	City string `json:"city" example:"New York"`

	// State or province name
	// @example "NY"
	State string `json:"state" example:"NY"`

	// Postal or ZIP code
	// @example "10001"
	PostalCode string `json:"postal_code,omitempty" example:"10001"`

	// Latitude coordinate
	// @example 40.7128
	Latitude float64 `json:"latitude" example:"40.7128"`

	// Longitude coordinate
	// @example -74.0060
	Longitude float64 `json:"longitude" example:"-74.0060"`
}

// PaginatedResponse represents a paginated response with metadata
// @Description Paginated data response with metadata about pagination
type PaginatedResponse struct {
	// The actual data items
	// @example [{"id":"a1b2c3d4","tracking_number":"DEL-230512-7890"}]
	Data interface{} `json:"data"`

	// Total number of items across all pages
	// @example 100
	TotalItems int64 `json:"total_items" example:"100"`

	// Current page number
	// @example 1
	Page int `json:"page" example:"1"`

	// Number of items per page
	// @example 20
	PageSize int `json:"page_size" example:"20"`

	// Total number of pages
	// @example 5
	TotalPages int `json:"total_pages" example:"5"`
}

// PickupAddressResponse contains the origin address details
// @Description Pickup address information
type PickupAddressResponse struct {
	// Name of the contact person at pickup location
	// @example "Jane Smith"
	ContactName string `json:"contact_name" example:"Jane Smith"`

	// Contact phone number at pickup location
	// @example "+0987654321"
	ContactPhone string `json:"contact_phone" example:"+0987654321"`

	// First line of the address
	// @example "456 Business Ave"
	AddressLine1 string `json:"address_line1" example:"456 Business Ave"`

	// Second line of the address (optional)
	// @example "Suite 300"
	AddressLine2 string `json:"address_line2,omitempty" example:"Suite 300"`

	// City name
	// @example "Chicago"
	City string `json:"city" example:"Chicago"`

	// State or province name
	// @example "IL"
	State string `json:"state" example:"IL"`

	// Postal or ZIP code
	// @example "60606"
	PostalCode string `json:"postal_code,omitempty" example:"60606"`

	// Latitude coordinate
	// @example 41.8781
	Latitude float64 `json:"latitude" example:"41.8781"`

	// Longitude coordinate
	// @example -87.6298
	Longitude float64 `json:"longitude" example:"-87.6298"`

	// Additional notes about the address
	// @example "Enter through loading dock"
	AddressNotes string `json:"address_notes,omitempty" example:"Enter through loading dock"`

	// Full formatted address
	// @example "456 Business Ave, Suite 300, Chicago, IL 60606"
	FormattedAddress string `json:"formatted_address,omitempty" example:"456 Business Ave, Suite 300, Chicago, IL 60606"`
}

// OrderListResponse is a simplified order representation for listings
// @Description Simplified order information for list views
type OrderListResponse struct {
	// Unique identifier of the order
	// @example a1b2c3d4-e5f6-7g8h-9i0j-k1l2m3n4o5p6
	ID string `json:"id" example:"a1b2c3d4-e5f6-7g8h-9i0j-k1l2m3n4o5p6"`

	// Tracking number for the order
	// @example DEL-230512-7890
	TrackingNumber string `json:"tracking_number" example:"DEL-230512-7890"`

	// Full name of the client
	// @example "John Smith"
	ClientName string `json:"client_name" example:"John Smith"`

	// Shortened delivery address for display
	// @example "123 Main St, New York, NY"
	DeliveryAddress string `json:"delivery_address" example:"123 Main St, New York, NY"`

	// Deadline for delivery
	// @example 2023-05-15T16:30:00Z
	DeliveryDeadline time.Time `json:"delivery_deadline" example:"2023-05-15T16:30:00Z" format:"date-time"`

	// Price of the delivery
	// @example 25.50
	Price float64 `json:"price" example:"25.50"`

	// Current status of the order
	// @example PENDING
	// @enum [PENDING,ACCEPTED,PICKED_UP,IN_TRANSIT,DELIVERED,CANCELLED]
	Status string `json:"status" example:"PENDING" enums:"PENDING,ACCEPTED,PICKED_UP,IN_TRANSIT,DELIVERED,CANCELLED"`

	// Driver ID assigned to the order (if any)
	// @example d1e2f3g4-h5i6-j7k8-l9m0-n1o2p3q4r5s6
	DriverID *string `json:"driver_id,omitempty" example:"d1e2f3g4-h5i6-j7k8-l9m0-n1o2p3q4r5s6"`

	// Full name of the assigned driver
	// @example "Michael Johnson"
	DriverName string `json:"driver_name,omitempty" example:"Michael Johnson"`

	// When the order was created
	// @example 2023-05-15T10:30:00Z
	CreatedAt time.Time `json:"created_at" example:"2023-05-15T10:30:00Z" format:"date-time"`
}

// OrderDriverAssignRequest represents the request to assign a driver to an order
// @Description Request to assign a driver to an order
type OrderDriverAssignRequest struct {
	// ID of the driver to assign
	// @example d1e2f3g4-h5i6-j7k8-l9m0-n1o2p3q4r5s6
	// @required
	DriverID string `json:"driver_id" example:"d1e2f3g4-h5i6-j7k8-l9m0-n1o2p3q4r5s6" binding:"required,uuid"`
}

// OrderUpdateRequest represents the request body for updating an existing order
// @Description Request structure for updating a delivery order
type OrderUpdateRequest struct {
	// Price of the delivery in local currency
	// @minimum 0
	// @example 25.50
	Price float64 `json:"price,omitempty" example:"25.50" binding:"omitempty,min=0"`

	// Distance of the delivery in kilometers
	// @minimum 0
	// @example 7.2
	Distance float64 `json:"distance,omitempty" example:"7.2" binding:"omitempty,min=0"`

	// Scheduled pickup time - only modifiable if order is still in PENDING state
	// @example 2023-05-15T14:30:00Z
	PickupTime *time.Time `json:"pickup_time,omitempty" example:"2023-05-15T14:30:00Z" format:"date-time"`

	// Deadline for delivery completion - only modifiable if order is still in PENDING state
	// @example 2023-05-15T16:30:00Z
	DeliveryDeadline *time.Time `json:"delivery_deadline,omitempty" example:"2023-05-15T16:30:00Z" format:"date-time"`

	// Whether recipient signature is required for delivery
	// @example false
	RequiresSignature *bool `json:"requires_signature,omitempty" example:"false"`

	// Additional notes for the delivery
	// @example "Please call recipient 5 minutes before arrival"
	DeliveryNotes string `json:"delivery_notes,omitempty" example:"Please call recipient 5 minutes before arrival"`

	// Details about the package being delivered
	PackageDetails *PackageDetailUpdateRequest `json:"package_details,omitempty"`

	// Delivery destination address details
	DeliveryAddress *DeliveryAddressUpdateRequest `json:"delivery_address,omitempty"`

	// Pickup contact information and notes
	// Contact name for pickup location
	// @example "Jane Smith"
	PickupContactName string `json:"pickup_contact_name,omitempty" example:"Jane Smith"`

	// Contact phone number for pickup location
	// @example "+0987654321"
	PickupContactPhone string `json:"pickup_contact_phone,omitempty" example:"+0987654321"`

	// Important notes about the pickup location
	// @example "Enter through loading dock"
	PickupNotes string `json:"pickup_notes,omitempty" example:"Enter through loading dock"`
}

func (o *OrderUpdateRequest) Validate() error {
	// Time validations - ensure delivery deadline is after pickup time if both are provided
	if o.PickupTime != nil && o.DeliveryDeadline != nil {
		if o.DeliveryDeadline.Before(*o.PickupTime) {
			return error2.NewGeneralServiceError("OrderUpdateDTO", "Validate", errPackage.ErrDeliveryDeadlineBeforePickup)
		}
	}

	return nil
}

// PackageDetailUpdateRequest contains details about the package for updates
// @Description Package characteristics and handling information for updates
type PackageDetailUpdateRequest struct {
	// Whether the package contains fragile items
	// @example true
	IsFragile *bool `json:"is_fragile,omitempty" example:"true"`

	// Whether the package requires urgent handling
	// @example false
	IsUrgent *bool `json:"is_urgent,omitempty" example:"false"`

	// Weight of the package in kilograms
	// @minimum 0
	// @example 2.5
	Weight *float64 `json:"weight,omitempty" example:"2.5" binding:"omitempty,min=0"`

	// Special handling instructions
	// @example "Contains glass items, handle with care"
	SpecialInstructions string `json:"special_instructions,omitempty" example:"Contains glass items, handle with care"`

	// Length of the package in centimeters
	// @minimum 0
	// @example 30
	Length *float64 `json:"length,omitempty" example:"30" binding:"omitempty,min=0"`

	// Width of the package in centimeters
	// @minimum 0
	// @example 20
	Width *float64 `json:"width,omitempty" example:"20" binding:"omitempty,min=0"`

	// Height of the package in centimeters
	// @minimum 0
	// @example 15
	Height *float64 `json:"height,omitempty" example:"15" binding:"omitempty,min=0"`
}

// DeliveryAddressUpdateRequest contains the destination address details for updates
// @Description Delivery destination address information for updates
type DeliveryAddressUpdateRequest struct {
	// Name of the person receiving the package
	// @example "John Doe"
	RecipientName string `json:"recipient_name,omitempty" example:"John Doe"`

	// Contact phone number of the recipient
	// @example "+1234567890"
	RecipientPhone string `json:"recipient_phone,omitempty" example:"+1234567890"`

	// First line of the address
	// @example "123 Main Street"
	AddressLine1 string `json:"address_line1,omitempty" example:"123 Main Street"`

	// Second line of the address (optional)
	// @example "Apartment 4B"
	AddressLine2 string `json:"address_line2,omitempty" example:"Apartment 4B"`

	// City name
	// @example "New York"
	City string `json:"city,omitempty" example:"New York"`

	// State or province name
	// @example "NY"
	State string `json:"state,omitempty" example:"NY"`

	// Postal or ZIP code
	// @example "10001"
	PostalCode string `json:"postal_code,omitempty" example:"10001"`

	// Additional notes about the address
	// @example "Ring doorbell twice"
	AddressNotes string `json:"address_notes,omitempty" example:"Ring doorbell twice"`
}

// OrderStatusUpdateRequest represents the request to update an order's status
// @Description Request to change the status of an order
type OrderStatusUpdateRequest struct {
	// New status for the order
	// @example ACCEPTED
	// @enum [PENDING,ACCEPTED,PICKED_UP,IN_TRANSIT,DELIVERED,CANCELLED]
	// @required
	Status string `json:"status" example:"ACCEPTED" binding:"required" enums:"PENDING,ACCEPTED,PICKED_UP,IN_TRANSIT,DELIVERED,CANCELLED"`

	// Optional description about the status change
	// @example "Driver has accepted the order and is heading to pickup location"
	Description string `json:"description,omitempty" example:"Driver has accepted the order and is heading to pickup location"`
}
