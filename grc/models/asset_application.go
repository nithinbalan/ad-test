package models

import (
	"time"
)

type AssetApplication struct {
	ID uint `json:"id" gorm:"primary_key;AUTO_INCREMENT"`

	AssetRegisterID uint `json:"asset_register_id" gorm:"not null"`
	ApplicationID   uint `json:"application_id" gorm:"not null"`

	CreatedByID uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AssetRegister *AssetRegister `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	Application   *Application   `gorm:"foreignKey:ApplicationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
}

func (AssetApplication) TableName() string {
	return "asset_applications"
}