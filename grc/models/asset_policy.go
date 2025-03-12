package models

import (
	"time"
)

type AssetPolicy struct {
	ID              uint      `json:"id" gorm:"primary_key;auto_increment"`
	AssetRegisterID uint      `json:"asset_register_id"`
	DocumentID      uint      `json:"document_id"`
	CreatedBYID     uint      `json:"created_by_id" gorm:"not null"`
	UpdatedBYID     uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	AssetRegister         *AssetRegister          `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document              *Document               `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetPolicyChecklists *[]AssetPolicyChecklist `gorm:"foreignKey:AssetPolicyID"`
	CreatedByUser         *User                   `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser         *User                   `gorm:"foreignKey:UpdatedBYID"`
}

func (AssetPolicy) TableName() string {
	return "asset_policies"
}
