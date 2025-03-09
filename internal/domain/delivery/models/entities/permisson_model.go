package entities

import "time"

type Permission struct {
	ID          string    `gorm:"column:id;type:char(36);primary_key" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	Resource    string    `gorm:"column:resource;type:varchar(50);not null" json:"resource"`
	Action      string    `gorm:"column:action;type:varchar(50);not null" json:"action"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Permission) TableName() string {
	return "permissions"
}
