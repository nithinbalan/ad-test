package models

import (
	"time"
)

// UserSecurity represents the user_securities table in GORM
type UserSecurity struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint       `gorm:"not null" json:"user_id" binding:"required"`
	PasswordHash string     `gorm:"type:text;not null" json:"password_hash" binding:"required"`
	IsMFA        bool       `gorm:"type:tinyint;not null; default:0" json:"is_mfa"`
	OTP          *string    `gorm:"type:varchar(45);default:null" json:"otp"`
	OTPExpiryAt  *time.Time `gorm:"type:timestamps;default:null" json:"otp_expiry_at"`
	UpdatedAt    time.Time  `gorm:"type:timestamps;not null; default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`

	// User is a reference to the User model for the foreign key constraint
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

// TableName overrides the table name used by UserSecurity to `user_securities`
func (UserSecurity) TableName() string {
	return "user_securities"
}
