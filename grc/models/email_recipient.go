package models

import "time"

type EmailRecipient struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	EmailID       uint      `json:"email_id" gorm:"not null;index"`
	UserID        *uint     `json:"user_id" gorm:"default:null"`
	RecipientType string    `json:"recipient_type" gorm:"type:nvarchar(10);not null;default:'to'"`
	EmailAddress  string    `json:"email_address" gorm:"type:nvarchar(500);not null"`
	DisplayName   *string   `json:"display_name" gorm:"type:nvarchar(500);default:null"`
	IsExternal    bool      `json:"is_external" gorm:"not null;default:false"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (EmailRecipient) TableName() string {
	return "email_recipients"
}
