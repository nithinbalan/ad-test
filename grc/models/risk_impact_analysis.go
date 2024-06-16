package models

import (
	"time"
)

type RiskImpactAnalysis struct {
	ID                   uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskEvaluationID     uint       `json:"risk_evaluation_id" binding:"required" gorm:"not null"`
	RiskImpactCategoryID uint       `json:"risk_impact_category_id" binding:"required" gorm:"not null"`
	RiskMatrixImpactID   uint       `json:"risk_matrix_impact_id" binding:"required" gorm:"not null"`
	CreatedByID          uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBy            *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	RiskEvaluation               *RiskEvaluation                `gorm:"foreignKey:RiskEvaluationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskImpactCategory           *RiskImpactCategory            `gorm:"foreignKey:RiskImpactCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixImpact             *RiskMatrixImpact              `gorm:"foreignKey:RiskMatrixImpactID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskImpactAnalysisGuidelines []*RiskImpactAnalysisGuideline `gorm:"foreignKey:RiskImpactAnalysisID"`
	CreatedByUser                *User                          `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser                *User                          `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by RiskImpactAnalysis to `risk_impact_analyses`
func (RiskImpactAnalysis) TableName() string {
	return "risk_impact_analyses"
}
