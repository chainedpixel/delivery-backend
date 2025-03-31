package entities

import "time"

// CompanyQueryParams define los parámetros para consultar empresas
type CompanyQueryParams struct {
	// Filtros específicos para empresas
	Name              string     `json:"name,omitempty"`
	TaxID             string     `json:"tax_id,omitempty"`
	ContactEmail      string     `json:"contact_email,omitempty"`
	ContactPhone      string     `json:"contact_phone,omitempty"`
	IsActive          *bool      `json:"is_active,omitempty"`
	ContractStartDate *time.Time `json:"contract_start_date,omitempty"`
	ContractEndDate   *time.Time `json:"contract_end_date,omitempty"`
	IncludeDeleted    bool       `json:"include_deleted,omitempty"`

	// Parámetros de paginación heredados
	PaginationQueryParams
}

// BranchQueryParams define los parámetros para consultar sucursales
type BranchQueryParams struct {
	// Filtros específicos para sucursales
	CompanyID    string `json:"company_id,omitempty"`
	Name         string `json:"name,omitempty"`
	Code         string `json:"code,omitempty"`
	ContactName  string `json:"contact_name,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	ZoneID       string `json:"zone_id,omitempty"`
	IsActive     *bool  `json:"is_active,omitempty"`

	// Parámetros de paginación heredados
	PaginationQueryParams
}
