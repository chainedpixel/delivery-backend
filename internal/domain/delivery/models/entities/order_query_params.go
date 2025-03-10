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

	// Paginación
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`

	// Ordenamiento
	SortBy        string `json:"sort_by,omitempty"`
	SortDirection string `json:"sort_direction,omitempty"`
}
