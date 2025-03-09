package entities

import "time"

type Zone struct {
	ID              string    `json:"id" gorm:"column:id;type:char(36);primary_key"`
	Name            string    `json:"name" gorm:"column:name;type:varchar(100);not null"`
	Code            string    `json:"code" gorm:"column:code;type:varchar(20);not null"`
	Boundaries      string    `json:"boundaries" gorm:"column:boundaries;type:polygon;not null"`
	CenterPoint     string    `json:"center_point" gorm:"column:center_point;type:point;not null"`
	BaseRate        float64   `json:"base_rate" gorm:"column:base_rate;type:decimal(10,2);not null"`
	MaxDeliveryTime int       `json:"max_delivery_time" gorm:"column:max_delivery_time;type:int;not null"`
	IsActive        bool      `json:"is_active" gorm:"column:is_active;type:boolean;default:true"`
	PriorityLevel   int       `json:"priority_level" gorm:"column:priority_level;type:int;not null;default:1"`
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (Zone) TableName() string {
	return "zones"
}
