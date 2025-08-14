package models

import (
	"time"
)

type AssetAssociatedDocument struct {
	ID              uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AssetRegisterID uint      `json:"asset_register_id" binding:"required"`
	DocumentID      uint      `json:"document_id" binding:"required"`
	CreatedByID     uint      `json:"created_by_id" binding:"required"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AssetRegister *AssetRegister `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document      *Document      `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedByID"`
}

func (AssetAssociatedDocument) TableName() string {
	return "asset_associated_documents"
}
