package models

import (
	"time"
)

type ProcessAsset struct {
	ID              uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProcessID       uint      `json:"process_id" binding:"required"`
	AssetRegisterID uint      `json:"asset_register_id" binding:"required"`
	CreatedByID     uint      `json:"created_by_id" binding:"required"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Process       *Process       `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Asset         *AssetRegister `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedByID"`
}

func (ProcessAsset) TableName() string {
	return "process_assets"
}
