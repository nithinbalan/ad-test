package models

import (
	"time"
)

type BusinessImpactAssessment struct {
	ID               uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProcessID        *uint      `json:"process_id"`
	AssetRegisterID  *uint      `json:"asset_register_id"`
	ImpactCategoryID uint       `json:"impact_category_id" binding:"required"`
	RiskLevelID      *uint      `json:"risk_level_id"`
	Cost             *uint      `json:"cost"`
	Comments         *string    `json:"comments" gorm:"type:nvarchar(max)"`
	IsAvailable      bool       `json:"is_available"`
	CreatedByID      uint       `json:"created_by_id" binding:"required"`
	UpdatedByID      *uint      `json:"updated_by_id"`
	CreatedAt        time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt        *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	Process        *Process        `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetRegister  *AssetRegister  `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ImpactCategory *ImpactCategory `gorm:"foreignKey:ImpactCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskLevel      *RiskLevel      `gorm:"foreignKey:RiskLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser  *User           `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser  *User           `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by BusinessFunction to `business_functions`
func (BusinessImpactAssessment) TableName() string {
	return "business_impact_assessments"
}
