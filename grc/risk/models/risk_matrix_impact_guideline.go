package models

import (
	organizationModels "jethur/grc/organization/models"
	"time"
)

type RiskMatrixImpactGuideline struct {
	ID                   uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskImpactCategoryID uint       `json:"risk_impact_category_id" binding:"required" gorm:"not null"`
	RiskMatrixImpactID   uint       `json:"risk_matrix_impact_id" binding:"required" gorm:"not null"`
	Guideline            string     `json:"guideline" binding:"required,max=500" gorm:"type:varchar(500);not null"`
	Status               string     `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBy            uint       `json:"created_by" gorm:"not null"`
	UpdatedBy            *uint      `json:"updated_by" gorm:"default:null"`

	// Relationships
	RiskImpactCategory *RiskImpactCategory      `gorm:"foreignkey:RiskImpactCategoryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskMatrixImpact   *RiskMatrixImpact        `gorm:"foreignkey:RiskMatrixImpactID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser      *organizationModels.User `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser      *organizationModels.User `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by RiskMatrixImpactGuideline to `risk_matrix_impact_guidelines`
func (RiskMatrixImpactGuideline) TableName() string {
	return "risk_matrix_impact_guidelines"
}
