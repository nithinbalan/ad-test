package models

import (
	"time"
)

type UserLogin struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserID      uint      `json:"user_id" binding:"required" gorm:"type:int;not null"`
	LoginAt     time.Time `json:"login_at" gorm:"type:datetime;not null"`
	LoginStatus string    `json:"login_status" gorm:"type:varchar(20);not null"`
	IPAddress   string    `json:"ip_address" gorm:"type:varchar(45);not null"`
	Location    *string   `json:"location" gorm:"type:varchar(255)"`
	DeviceInfo  *string   `json:"device_info" gorm:"type:varchar(255)"`
	BrowserInfo *string   `json:"browser_info" gorm:"type:varchar(255)"`
}

// TableName overrides the table name used by UserLogin to `user_logins`
func (UserLogin) TableName() string {
	return "user_logins"
}
