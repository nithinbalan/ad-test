package models

import (
	"time"
)

type ProjectAsset struct {
	ID              uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProjectID       uint      `json:"project_id" binding:"required"`
	AssetRegisterID uint      `json:"asset_register_id" binding:"required"`
	CreatedBYID     uint      `json:"created_by_id" binding:"required"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;not null"`

	Project       *Project       `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetRegister *AssetRegister `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedBYID"`
}

func (ProjectAsset) TableName() string {
	return "project_assets"
}
