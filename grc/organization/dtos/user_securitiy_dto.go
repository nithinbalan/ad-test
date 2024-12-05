package dtos

import (
	"time"
)

type UserSecurityFormDTO struct {
	PasswordHash string     `json:"password_hash" binding:"required"`
	IsMFA        bool       `json:"is_mfa" binding:"required"`
	OTP          *string    `json:"otp" binding:"omitempty"`
	OTPExpiryAt  *time.Time `json:"otp_expiry_at" binding:"omitempty"`
}

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

type ChangePasswordDTO struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}
