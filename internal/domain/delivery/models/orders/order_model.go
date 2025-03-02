package orders

import (
	"domain/delivery/models/companies"
	"domain/delivery/models/drivers"
	"domain/delivery/models/users"
	"domain/delivery/models/warehouses"
	"time"
)

type Order struct {
	ID             string    `gorm:"column:id;type:char(36);primaryKey"`
	CompanyID      string    `gorm:"column:company_id;type:char(36);not null"`
	BranchID       string    `gorm:"column:branch_id;type:char(36);not null"`
	ClientID       string    `gorm:"column:client_id;type:char(36);not null"`
	DriverID       *string   `gorm:"column:driver_id;type:char(36)"`
	TrackingNumber string    `gorm:"column:tracking_number;type:varchar(50);not null"`
	Status         string    `gorm:"column:status;type:varchar(20);not null"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Inverse Relationships
	Company *companies.Company `gorm:"foreignKey:CompanyID;references:ID"`
	Branch  *companies.Branch  `gorm:"foreignKey:BranchID;references:ID"`
	Client  *users.User        `gorm:"foreignKey:ClientID;references:ID"`
	Driver  *drivers.Driver    `gorm:"foreignKey:DriverID;references:UserID"`

	// Relationships one to one
	Detail          *Details       `gorm:"foreignKey:OrderID"`
	PackageDetail   *PackageDetail `gorm:"foreignKey:OrderID"`
	DeliveryAddress *Address       `gorm:"foreignKey:OrderID"`
	PickupAddress   *PickupAddress `gorm:"foreignKey:OrderID"`
	Tracking        *Tracking      `gorm:"foreignKey:OrderID"`
	QRCode          *QRCode        `gorm:"foreignKey:OrderID"`

	// Relationships one to many
	StatusHistory      []StatusHistory              `gorm:"foreignKey:OrderID"`
	WarehouseTrackings []warehouses.PackageTracking `gorm:"foreignKey:OrderID"`
	WarehouseInventory *warehouses.Warehouse        `gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string {
	return "orders"
}
