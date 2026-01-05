package dtos

import (
	settingDTOs "jethur/grc/settings/dtos"
)

type LoginResponseDTO struct {
	AccessToken  *string `json:"access_token"`
	RefreshToken *string `json:"refresh_token"`
	IsMFAEnabled bool    `json:"is_mfa_enabled"`
}

type LoginFormDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenFormDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type SendOTPFormDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPFormDTO struct {
	Email string `json:"email" binding:"required,email"`
	Otp   string `json:"otp"`
}

type PasswordRecoveryDataDTO struct {
	FirstNameEn string `json:"first_name_en"`
	LastNameEn  string `json:"last_name_en"`
	Otp         string `json:"otp"`
}

type PasswordRecoverySetPasswordFormDTO struct {
	Token           string `json:"token"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type PasswordRecoverySetPasswordEmailDTO struct {
	FirstNameEn string `json:"first_name_en"`
	LastNameEn  string `json:"last_name_en"`
	Email       string `json:"email"`
}

type LoginSettingsDTO struct {
	ActiveDirectoryEnabled bool                       `json:"active_directory_enabled"`
	AzureADEnabled         bool                       `json:"azure_ad_enabled"`
	Theme                  settingDTOs.ThemeDetailDTO `json:"theme"`
}

type MicrosoftOAuthCallbackDTO struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state"`
}
