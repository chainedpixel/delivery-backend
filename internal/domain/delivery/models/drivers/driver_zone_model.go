package drivers

import (
	"domain/delivery/models/zones"
	"time"
)

type Zone struct {
	DriverID            string     `gorm:"column:driver_id;type:char(36);primaryKey"`
	ZoneID              string     `gorm:"column:zone_id;type:char(36);primaryKey"`
	IsPrimary           bool       `gorm:"column:is_primary;type:boolean;default:false"`
	EfficiencyRating    float64    `gorm:"column:efficiency_rating;type:decimal(3,2);default:5.00"`
	DeliveriesCompleted int        `gorm:"column:deliveries_completed;type:int;default:0"`
	LastDelivery        *time.Time `gorm:"column:last_delivery;type:timestamp"`
	IsActive            bool       `gorm:"column:is_active;type:boolean;default:true"`
	CreatedAt           time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Inverse Relationships
	Driver *Driver     `gorm:"foreignKey:DriverID;references:UserID"`
	Zone   *zones.Zone `gorm:"foreignKey:ZoneID;references:ID"`
}

func (Zone) TableName() string {
	return "driver_zones"
}
