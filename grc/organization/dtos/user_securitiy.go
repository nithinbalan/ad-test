package dtos

import (
	"time"
)

// UserSecurity represents the user_securities
type UserSecurity struct {
	ID           uint       `json:"id"`
	UserID       uint       `json:"user_id"`
	PasswordHash string     `json:"password_hash"`
	IsMFA        bool       `json:"is_mfa"`
	OTP          *string    `json:"otp"`
	OTPExpiryAt  *time.Time `json:"otp_expiry_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
