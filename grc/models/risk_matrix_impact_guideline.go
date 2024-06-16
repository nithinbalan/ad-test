package models

import (
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
	CreatedByID          uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBy            *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	RiskImpactCategory *RiskImpactCategory `gorm:"foreignKey:RiskImpactCategoryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskMatrixImpact   *RiskMatrixImpact   `gorm:"foreignKey:RiskMatrixImpactID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser      *User               `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by RiskMatrixImpactGuideline to `risk_matrix_impact_guidelines`
func (RiskMatrixImpactGuideline) TableName() string {
	return "risk_matrix_impact_guidelines"
}
