package models

import (
	"time"
)

type ViolationReportAsset struct {
	ID                uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ViolationReportID uint      `json:"violation_report_id" binding:"required"`
	AssetRegisterID   uint      `json:"asset_register_id" binding:"required"`
	CreatedBYID       uint      `json:"created_by_id" binding:"required"`
	CreatedAt         time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ViolationReport *ViolationReport `gorm:"foreignKey:ViolationReportID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetRegister   *AssetRegister   `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User            `gorm:"foreignKey:CreatedBYID"`
}

func (ViolationReportAsset) TableName() string {
	return "violation_report_assets"
}
