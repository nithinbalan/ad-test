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
	CreatedByID                 uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	UpdatedBy                   *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	RiskRegister              *RiskRegister              `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixLikelihood      *RiskMatrixLikelihood      `gorm:"foreignKey:RiskMatrixLikelihoodID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixImpact          *RiskMatrixImpact          `gorm:"foreignKey:RiskMatrixImpactID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixRiskRating      *RiskMatrixRiskRating      `gorm:"foreignKey:RiskMatrixRiskRatingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixRiskRatingValue *RiskMatrixRiskRatingValue `gorm:"foreignKey:RiskMatrixRiskRatingValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskImpactAnalyses        *[]RiskImpactAnalysis      `gorm:"foreignKey:RiskEvaluationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskAppetiteAlignments    *[]RiskAppetiteAlignment   `gorm:"foreignKey:RiskEvaluationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser             *User                      `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser             *User                      `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by RiskEvaluation to `risk_evaluations`
func (RiskEvaluation) TableName() string {
	return "risk_evaluations"
}
