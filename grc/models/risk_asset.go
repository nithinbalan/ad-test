package models

import (
	"time"
)

type RiskAsset struct {
	ID              uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID  uint      `json:"risk_register_id" binding:"required"`
	AssetRegisterID uint      `json:"asset_register_id" binding:"required"`
	CreatedByID     uint      `json:"created_by_id" binding:"required"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskRegister  *RiskRegister  `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetRegister *AssetRegister `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedByID"`
}

func (RiskAsset) TableName() string {
	return "risk_assets"
}
