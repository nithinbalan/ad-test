package models

import (
	"time"
)

type AssetMaintenance struct {
	ID              uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetRegisterID uint       `json:"asset_register_id" binding:"required" gorm:"not null"`
	Type            string     `json:"type" binding:"required" gorm:"type:nvarchar(500);not null"`
	PerformedDate   time.Time  `json:"performed_date" binding:"required" gorm:"type:date;not null"`
	Title           string     `json:"title" binding:"required" gorm:"type:nvarchar(500);not null"`
	Description     *string    `json:"description" gorm:"type:nvarchar(max)"`
	PerformedBYID   uint       `json:"performed_by_id" binding:"required" gorm:"not null"`
	Outcome         *string    `json:"outcome" gorm:"type:nvarchar(max)"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBYID     uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID     *uint      `json:"updated_by_id"`

	// Relationships
	AssetRegister *AssetRegister `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PerformedBy   *User          `gorm:"foreignKey:PerformedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser *User          `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the default table name
func (AssetMaintenance) TableName() string {
	return "asset_maintenance"
}
