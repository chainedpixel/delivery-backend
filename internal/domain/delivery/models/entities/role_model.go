package entities

import "time"

type Role struct {
	ID          string       `gorm:"column:id;type:char(36);primary_key" json:"id"`
	Name        string       `gorm:"column:name;type:varchar(50);not null" json:"name"`
	Description string       `gorm:"column:description;type:text" json:"description"`
	IsActive    bool         `gorm:"column:is_active;type:boolean;default:true" json:"is_active"`
	CreatedAt   time.Time    `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"-"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"-"`
}

func (Role) TableName() string {
	return "roles"
}
