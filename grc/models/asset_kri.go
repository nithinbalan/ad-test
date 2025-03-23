package models

import (
	"time"
)

type AssetKri struct {
	ID              uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AssetRegisterID uint      `json:"asset_register_id" binding:"required"`
	KriRegisterID   uint      `json:"kri_register_id" binding:"required"`
	CreatedBYID     uint      `json:"created_by_id" binding:"required"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;not null"`

	AssetRegister *AssetRegister `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	KriRegister   *KRIRegister   `gorm:"foreignKey:KriRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedBYID"`
}

func (AssetKri) TableName() string {
	return "asset_kris"
}
