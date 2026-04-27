package models

import (
	"time"
)

type RiskResidualEvaluationHistory struct {
	ID                               uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskEvaluationID                 uint      `json:"risk_evaluation_id" binding:"required" gorm:"not null"`
	RiskRegisterID                   uint      `json:"risk_register_id" binding:"required" gorm:"not null"`
	RiskMatrixLikelihoodID           uint      `json:"risk_matrix_likelihood_id" binding:"required" gorm:"not null"`
	RiskMatrixRiskRatingID           uint      `json:"risk_matrix_risk_rating_id" binding:"required" gorm:"not null"`
	RiskMatrixRiskRatingValueID      uint      `json:"risk_matrix_risk_rating_value_id" binding:"required" gorm:"not null"`
	AverageControlEffectivenessScore *float64  `json:"average_control_effectiveness_score" gorm:"type:decimal(3,2);default:null"`
	CreatedBYID                      uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                        time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskEvaluation            *RiskEvaluation            `gorm:"foreignKey:RiskEvaluationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskRegister              *RiskRegister              `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixLikelihood      *RiskMatrixLikelihood      `gorm:"foreignKey:RiskMatrixLikelihoodID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixRiskRating      *RiskMatrixRiskRating      `gorm:"foreignKey:RiskMatrixRiskRatingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixRiskRatingValue *RiskMatrixRiskRatingValue `gorm:"foreignKey:RiskMatrixRiskRatingValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser             *User                      `gorm:"foreignKey:CreatedBYID"`
}

// TableName overrides the table name used by RiskResidualEvaluationHistory to `risk_residual_evaluation_histories`
func (RiskResidualEvaluationHistory) TableName() string {
	return "risk_residual_evaluation_histories"
}
