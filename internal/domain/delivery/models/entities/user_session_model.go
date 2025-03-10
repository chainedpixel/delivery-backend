package entities

import "time"

type UserSession struct {
	ID           string    `gorm:"column:id;type:char(36);primary_key" json:"id"`
	UserID       string    `gorm:"column:user_id;type:char(36);not null;index" json:"user_id"`
	Token        string    `gorm:"column:token;type:varchar(255);not null;index" json:"token"`
	DeviceInfo   string    `gorm:"column:device_info;type:text" json:"device_info,omitempty"`
	IPAddress    string    `gorm:"column:ip_address;type:varchar(45)" json:"ip_address"`
	LastActivity time.Time `gorm:"column:last_activity;type:timestamp" json:"last_activity"`
	ExpiresAt    time.Time `gorm:"column:expires_at;type:timestamp;not null" json:"expires_at"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Inverse Relationships
	User *User `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
