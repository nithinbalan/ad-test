package models

import (
	"time"
)

type MeetingAsset struct {
	ID              uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	MeetingID       uint      `json:"meeting_id" binding:"required"`
	AssetRegisterID uint      `json:"asset_register_id" binding:"required"`
	CreatedBYID     uint      `json:"created_by_id" binding:"required"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;not null"`

	Meeting       *Meeting       `gorm:"foreignKey:MeetingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetRegister *AssetRegister `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedBYID"`
}

func (MeetingAsset) TableName() string {
	return "meeting_assets"
}
