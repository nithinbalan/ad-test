package models

import (
	"time"
)

type AssetCriticalityMatrixAssetCategory struct {
	ID                       uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AssetCriticalityMatrixID uint      `json:"asset_criticality_matrix_id" binding:"required" gorm:"not null"`
	AssetCategoryID          uint      `json:"asset_category_id" binding:"required" gorm:"not null"`
	CreatedBYID              uint      `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AssetCriticalityMatrix *AssetCriticalityMatrix `gorm:"foreignKey:AssetCriticalityMatrixID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetCategory          *AssetCategory          `gorm:"foreignKey:AssetCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID"`
}

func (AssetCriticalityMatrixAssetCategory) TableName() string {
	return "asset_criticality_matrix_asset_categories"
}
