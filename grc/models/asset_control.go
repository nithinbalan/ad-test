package models

import (
	"time"
)

type AssetControl struct {
	ID              uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetRegisterID uint      `json:"asset_register_id" binding:"required"`
	ControlID       uint      `json:"control_id" binding:"required"`
	CreatedBYID     uint      `json:"created_by_id" binding:"required"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AssetRegister *AssetRegister `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control       *Control       `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedBYID"`
}

func (AssetControl) TableName() string {
	return "asset_controls"
}
