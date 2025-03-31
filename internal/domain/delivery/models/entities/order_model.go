package entities

import (
	errPackage "github.com/MarlonG1/delivery-backend/internal/domain/error"
	"time"
)

type Order struct {
	ID             string     `gorm:"column:id;type:char(36);primaryKey"`
	CompanyID      string     `gorm:"column:company_id;type:char(36);not null"`
	BranchID       string     `gorm:"column:branch_id;type:char(36);not null"`
	ClientID       string     `gorm:"column:client_id;type:char(36);not null"`
	DriverID       *string    `gorm:"column:driver_id;type:char(36)"`
	TrackingNumber string     `gorm:"column:tracking_number;type:varchar(50);not null"`
	Status         string     `gorm:"column:status;type:varchar(20);not null"`
	CreatedAt      time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt      *time.Time `gorm:"column:deleted_at;type:timestamp;index"`

	// Inverse Relationships
	Company *Company `gorm:"foreignKey:CompanyID;references:ID"`
	Branch  *Branch  `gorm:"foreignKey:BranchID;references:ID"`
	Client  *User    `gorm:"foreignKey:ClientID;references:ID"`
	Driver  *Driver  `gorm:"foreignKey:DriverID;references:UserID"`

	// Relationships one to one
	Detail          *Details         `gorm:"foreignKey:OrderID"`
	PackageDetail   *PackageDetail   `gorm:"foreignKey:OrderID"`
	DeliveryAddress *DeliveryAddress `gorm:"foreignKey:OrderID"`
	PickupAddress   *PickupAddress   `gorm:"foreignKey:OrderID"`
	Tracking        *Tracking        `gorm:"foreignKey:OrderID"`
	QRCode          *QRCode          `gorm:"foreignKey:OrderID"`

	// Relationships one to many
	StatusHistory      []StatusHistory   `gorm:"foreignKey:OrderID"`
	WarehouseTrackings []PackageTracking `gorm:"foreignKey:OrderID"`
	WarehouseInventory []Inventory       `gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string {
	return "orders"
}

func (o *Order) Validate() error {
	if o.CompanyID == "" || o.BranchID == "" || o.ClientID == "" {
		return errPackage.ErrIDsNotFound
	}

	if o.TrackingNumber == "" {
		return errPackage.ErrTrackingNumber
	}

	if o.Status == "" {
		return errPackage.ErrStatusRequired
	}

	if o.Detail == nil {
		return errPackage.ErrOrderDetailsRequired
	}

	if o.DeliveryAddress == nil {
		return errPackage.ErrDeliveryAddress
	}

	if o.PickupAddress == nil {
		return errPackage.ErrPickupAddress
	}

	if o.PackageDetail == nil {
		return errPackage.ErrPackageDetails
	}

	return nil
}
