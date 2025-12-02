package models

import (
	"time"
)

type AuditFindingAsset struct {
	ID              uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditFindingID  uint      `json:"audit_finding_id" binding:"required"`
	AssetRegisterID uint      `json:"asset_register_id" binding:"required"`
	CreatedBYID     uint      `json:"created_by_id" binding:"required"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;not null"`

	AuditFinding  *AuditFinding  `gorm:"foreignKey:AuditFindingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetRegister *AssetRegister `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedBYID"`
}

func (AuditFindingAsset) TableName() string {
	return "audit_finding_assets"
}
