package models

import "time"

type IntegrationCredential struct {
	ID             uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
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

// TableName overrides the table name used by Integration to `integrations`
func (IntegrationCredential) TableName() string {
	return "integration_credentials"
}
