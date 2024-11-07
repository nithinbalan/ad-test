package models

import (
	"time"
)

type AssetCriticalityContribution struct {
	ID                                   uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetCriticalityMatrixFactorID       uint       `json:"asset_criticality_matrix_factor_id" gorm:"not null" binding:"required"`
	AssetCriticalityMatrixFactorRatingID uint       `json:"asset_criticality_matrix_factor_rating_id" gorm:"not null" binding:"required"`
	Score                                float64    `json:"score" gorm:"not null"`
	CreatedByID                          uint       `json:"created_by_id" gorm:"not null" binding:"required"`
	CreatedAt                            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedByID                          *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt                            *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	AssetCriticalityMatrixFactor       *AssetCriticalityMatrixFactor       `gorm:"foreignKey:AssetCriticalityMatrixFactorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetCriticalityMatrixFactorRating *AssetCriticalityMatrixFactorRating `gorm:"foreignKey:AssetCriticalityMatrixFactorRatingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                      *User                               `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser                      *User                               `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the default table name.
func (AssetCriticalityContribution) TableName() string {
	return "asset_criticality_contributions"
}
