package models

import (
	"time"
)

type AssetDocument struct {
	ID              uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetRegisterID uint      `json:"asset_register_id" gorm:"not null"`
	URL             string    `json:"url"`
	Name            string    `json:"name" binding:"required"`
	Size            int       `json:"size"`
	Token           string    `json:"token"`
	OriginalName    string    `json:"original_name"`
	Type            string    `json:"type"`
	CreatedBYID     uint      `json:"created_by_id" gorm:"not null"`
	UpdatedBYID     uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	AssetRegister *AssetRegister `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User          `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by AssetDocument to `asset_documents`
func (AssetDocument) TableName() string {
	return "asset_documents"
}
