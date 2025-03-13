package entities

import "time"

type UserSession struct {
	ID           string    `gorm:"column:id;type:char(36);primary_key" json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	UserID       string    `gorm:"column:user_id;type:char(36);not null;index" json:"-"`
	Token        string    `gorm:"column:token;type:varchar(255);not null;index" json:"-"`
	DeviceInfo   string    `gorm:"column:device_info;type:text" json:"device_info,omitempty" example:"{\"os\":\"android\",\"version\":\"10\"}"`
	IPAddress    string    `gorm:"column:ip_address;type:varchar(45)" json:"ip_address" example:"200.43.52.1"`
	LastActivity time.Time `gorm:"column:last_activity;type:timestamp" json:"last_activity" example:"2021-01-01T00:00:00Z"`
	ExpiresAt    time.Time `gorm:"column:expires_at;type:timestamp;not null" json:"expires_at" example:"2021-01-01T00:00:00Z"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at" example:"2021-01-01T00:00:00Z"`

	// Inverse Relationships
	User *User `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
