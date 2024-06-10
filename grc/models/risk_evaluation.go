package models

import (
	"time"
)

type RiskEvaluation struct {
	ID                          uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Type                        string     `json:"type" binding:"required,max=10,oneof=inherent residual" gorm:"type:varchar(10);not null"`
	RiskRegisterID              uint       `json:"risk_register_id" binding:"required" gorm:"not null"`
	RiskMatrixLikelihoodID      uint       `json:"risk_matrix_likelihood_id" binding:"required" gorm:"not null"`
	RiskMatrixImpactID          uint       `json:"risk_matrix_impact_id" binding:"required" gorm:"not null"`
	RiskMatrixRiskRatingID      uint       `json:"risk_matrix_risk_rating_id" binding:"required" gorm:"not null"`
	RiskMatrixRiskRatingValueID uint       `json:"risk_matrix_risk_rating_value_id" binding:"required" gorm:"not null"`
	CreatedBy                   uint       `json:"created_by" gorm:"not null"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	UpdatedBy                   *uint      `json:"updated_by" gorm:"default:null"`

	// Relationships
	RiskRegister              *RiskRegister              `gorm:"foreignkey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixLikelihood      *RiskMatrixLikelihood      `gorm:"foreignkey:RiskMatrixLikelihoodID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixImpact          *RiskMatrixImpact          `gorm:"foreignkey:RiskMatrixImpactID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixRiskRating      *RiskMatrixRiskRating      `gorm:"foreignkey:RiskMatrixRiskRatingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixRiskRatingValue *RiskMatrixRiskRatingValue `gorm:"foreignkey:RiskMatrixRiskRatingValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskImpactAnalyses        *[]RiskImpactAnalysis      `gorm:"foreignkey:RiskEvaluationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskAppetiteAlignments    *[]RiskAppetiteAlignment   `gorm:"foreignkey:RiskEvaluationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser             *User                      `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser             *User                      `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by RiskEvaluation to `risk_evaluations`
func (RiskEvaluation) TableName() string {
	return "risk_evaluations"
}
