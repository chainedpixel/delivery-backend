package entities

import "time"

type Role struct {
	ID          string       `gorm:"column:id;type:char(36);primary_key" json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	Name        string       `gorm:"column:name;type:varchar(50);not null" json:"name" example:"admin"`
	Description string       `gorm:"column:description;type:text" json:"description" example:"Administrator role"`
	IsActive    bool         `gorm:"column:is_active;type:boolean;default:true" json:"is_active" example:"true"`
	CreatedAt   time.Time    `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time    `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"-"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"-"`
}

func (Role) TableName() string {
	return "roles"
}
