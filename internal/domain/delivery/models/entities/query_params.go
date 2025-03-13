package entities

import "time"

type OrderQueryParams struct {
	// Filtros
	Status         string     `json:"status,omitempty"`
	TrackingNumber string     `json:"tracking_number,omitempty"`
	Location       string     `json:"location,omitempty"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	IncludeDeleted bool       `json:"include_deleted,omitempty"`

	PaginationQueryParams
}

type UserQueryParams struct {
	// Filtros
	Status         bool       `json:"status,omitempty"`
	CreationDate   *time.Time `json:"creation_date,omitempty"`
	Phone          string     `json:"phone,omitempty"`
	Name           string     `json:"name,omitempty"`
	Email          string     `json:"email,omitempty"`
	IncludeDeleted bool       `json:"include_deleted,omitempty"`

	PaginationQueryParams
}
