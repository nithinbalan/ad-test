package models

import (
	"time"
)

type ThirdPartyAsset struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ThirdPartyRegisterID uint      `json:"third_party_register_id" binding:"required"`
	AssetRegisterID      uint      `json:"asset_register_id" binding:"required"`
	CreatedBYID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	ThirdPartyRegister *ThirdPartyRegister `gorm:"foreignKey:ThirdPartyRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetRegister      *AssetRegister      `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBYID"`
}

func (ThirdPartyAsset) TableName() string {
	return "third_party_assets"
}
