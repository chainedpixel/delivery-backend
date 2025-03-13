package constants

var (
	AdminRole      = "ADMIN"
	CompanyUser    = "COMPANY_USER"
	Driver         = "DRIVER"
	WarehouseStaff = "WAREHOUSE_STAFF"
	Collector      = "COLLECTOR"
	FinalUser      = "FINAL_USER"
)

var ValidRoles = map[string]bool{
	AdminRole:      true,
	CompanyUser:    true,
	Driver:         true,
	WarehouseStaff: true,
	Collector:      true,
	FinalUser:      true,
}
