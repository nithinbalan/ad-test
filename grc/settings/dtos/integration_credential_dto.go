package dtos

import "time"

type IntegrationCredentialDTO struct {
	ID             int        `json:"id"`
	IntegrationID  uint       `json:"integration_id"`
	ApiKey         *string    `json:"api_key"`
	ApiSecret      *string    `json:"api_secret"`
	AccessToken    *string    `json:"access_token"`
	RedirectUrl    *string    `json:"redirect_url"`
	TokenType      *string    `json:"token_type"`
	Scope          *string    `json:"scope"`
	ExpireAt       *string    `json:"expire_at"`
	RefreshToken   *string    `json:"refresh_token"`
	IsActive       *bool      `json:"is_active"`
	AdditionalInfo *string    `json:"additional_info"`
	CreatedByID    uint       `json:"created_by_id"`
	UpdatedByID    *uint      `json:"updated_by_id"`
	CreatedAt      time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt      *time.Time `json:"updated_at" gorm:"type:datetime"`
}

type IntegrationCredentialFormDTO struct {
	IntegrationID  uint       `json:"integration_id" binding:"required"`
	ApiKey         *string    `json:"api_key"`
	ApiSecret      *string    `json:"api_secret"`
	AccessToken    *string    `json:"access_token"`
	RedirectUrl    *string    `json:"redirect_url"`
	TokenType      *string    `json:"token_type"`
	Scope          *string    `json:"scope"`
	ExpireAt       *string    `json:"expire_at"`
	RefreshToken   *string    `json:"refresh_token"`
	IsActive       *bool      `json:"is_active"`
	AdditionalInfo *string    `json:"additional_info"`
	CreatedByID    uint       `json:"created_by_id"`
	UpdatedByID    *uint      `json:"updated_by_id"`
	CreatedAt      time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt      *time.Time `json:"updated_at" gorm:"type:datetime"`
}
