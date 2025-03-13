package entities

type PaginationQueryParams struct {
	// Paginación
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`

	// Ordenamiento
	SortBy        string `json:"sort_by,omitempty"`
	SortDirection string `json:"sort_direction,omitempty"`
}
