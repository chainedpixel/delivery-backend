package entities

import (
	"time"
)

type Driver struct {
	UserID              string     `gorm:"column:user_id;type:char(36);primaryKey"`
	LicenseNumber       string     `gorm:"column:license_number;type:varchar(50);not null"`
	LicenseExpiry       time.Time  `gorm:"column:license_expiry;type:date;not null"`
	VehicleType         string     `gorm:"column:vehicle_type;type:varchar(50);not null"`
	VehiclePlate        string     `gorm:"column:vehicle_plate;type:varchar(20);not null"`
	VehicleModel        string     `gorm:"column:vehicle_model;type:varchar(100);not null"`
	VehicleColor        string     `gorm:"column:vehicle_color;type:varchar(50);not null"`
	IsActive            bool       `gorm:"column:is_active;type:boolean;default:true"`
	VehicleDetails      string     `gorm:"column:vehicle_details;type:json"`
	Documentation       string     `gorm:"column:documentation;type:json"`
	Rating              float64    `gorm:"column:rating;type:decimal(3,2);default:5.00"`
	CompletedDeliveries int        `gorm:"column:completed_deliveries;type:int;default:0"`
	LastDelivery        *time.Time `gorm:"column:last_delivery;type:timestamp"`
	CreatedAt           time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt           time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Inverse Relationships
	User *User `gorm:"foreignKey:UserID;references:ID"`

	// Relationships
	DriverZones  []DriverZone  `gorm:"foreignKey:DriverID"`
	Availability *Availability `gorm:"foreignKey:DriverID"`
	Orders       []Order       `gorm:"foreignKey:DriverID"`
}

func (Driver) TableName() string {
	return "drivers"
}
