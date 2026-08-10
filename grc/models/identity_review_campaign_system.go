package models

import "time"

type IdentityReviewCampaignSystem struct {
	ID                       uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IdentityReviewCampaignID uint      `json:"identity_review_campaign_id" gorm:"not null"`
	SystemRegisterID         *uint     `json:"system_register_id,omitempty" gorm:"default:null"`
	AssetRegisterID          *uint     `json:"asset_register_id"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`
	CreatedByID              uint      `json:"created_by_id" gorm:"not null"`

	IdentityReviewCampaign *IdentityReviewCampaign `json:"-" gorm:"foreignKey:IdentityReviewCampaignID"`
	SystemRegister         *SystemRegister         `json:"system_register" gorm:"foreignKey:SystemRegisterID"`
	AssetRegister          *AssetRegister          `json:"asset_register,omitempty" gorm:"foreignKey:AssetRegisterID"`
	CreatedByUser          *User                   `json:"created_by" gorm:"foreignKey:CreatedByID"`
}

func (IdentityReviewCampaignSystem) TableName() string {
	return "identity_review_campaign_systems"
}
