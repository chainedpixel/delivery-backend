package dto

// PaginatedResponse represents a paginated response with metadata
// @Description Paginated data response with metadata about pagination
type PaginatedResponse struct {
	// The actual data items (array of objects)
	Data interface{} `json:"data"`

	// Total number of items across all pages
	TotalItems int64 `json:"total_items" example:"100"`

	// Current page number
	Page int `json:"page" example:"1"`

	// Number of items per page
	PageSize int `json:"page_size" example:"20"`

	// Total number of pages
	TotalPages int `json:"total_pages" example:"5"`
}
