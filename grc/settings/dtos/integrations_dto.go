package dtos

import "time"

type IntegrationDTO struct {
	ID                         int                       `json:"id"`
	IntegrationCategoryID      uint                      `json:"integration_category_id"`
	IntegrationCategoryTitleEn string                    `json:"integration_category_title_en"`
	IntegrationCategoryTitleAr string                    `json:"integration_category_title_ar"`
	TitleEn                    string                    `json:"title_en"`
	TitleAr                    string                    `json:"title_ar"`
	DescriptionEn              *string                   `json:"description_en"`
	DescriptionAr              *string                   `json:"description_ar"`
	ApiBaseUrl                 *string                   `json:"api_base_url"`
	LogoUrl                    *string                   `json:"logo_url"`
	IsEnable                   bool                      `json:"is_enable"`
	CreatedByID                uint                      `json:"created_by_id"`
	UpdatedByID                *uint                     `json:"updated_by_id"`
	CreatedAt                  time.Time                 `json:"created_at"`
	UpdatedAt                  *time.Time                `json:"updated_at"`
	IntegrationCredential      *IntegrationCredentialDTO `json:"credential" gorm:"-"`
}

type IntegrationListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}

type IntegrationDetailDTO struct {
	ID                    uint                     `json:"id"`
	TitleEn               string                   `json:"title_en"`
	TitleAr               string                   `json:"title_ar"`
	DescriptionEn         *string                  `json:"description_en"`
	DescriptionAr         *string                  `json:"description_ar"`
	ApiBaseUrl            *string                  `json:"api_base_url"`
	LogoUrl               *string                  `json:"logo_url"`
	IsEnable              bool                     `json:"is_enable"`
	IntegrationCredential IntegrationCredentialDTO `json:"integration_credential"`
	CreatedAt             time.Time                `json:"created_at"`
	UpdatedAt             *time.Time               `json:"updated_at"`
	CreatedBYUser         UserMinimalDTO           `json:"created_by"`
	UpdatedBYUser         UserMinimalDTO           `json:"updated_by"`
}

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

type IntegrationCredentialDetailDTO struct {
	ID             uint           `json:"id"`
	ApiKey         *string        `json:"api_key"`
	ApiSecret      *string        `json:"api_secret"`
	AccessToken    *string        `json:"access_token"`
	RedirectUrl    *string        `json:"redirect_url"`
	TokenType      *string        `json:"token_type"`
	Scope          *string        `json:"scope"`
	ExpireAt       *string        `json:"expire_at"`
	RefreshToken   *string        `json:"refresh_token"`
	IsActive       *bool          `json:"is_active"`
	AdditionalInfo *string        `json:"additional_info"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      *time.Time     `json:"updated_at"`
	CreatedBYUser  UserMinimalDTO `json:"created_by"`
	UpdatedBYUser  UserMinimalDTO `json:"updated_by"`
}
