package models

import (
	"time"
)

type RiskEvaluation struct {
	ID                             uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID                 uint       `json:"risk_register_id" binding:"required" gorm:"not null"`
	InherentRiskMatrixLikelihoodID uint       `json:"inherent_risk_matrix_likelihood_id" binding:"required" gorm:"not null"`
	InherentRiskMatrixImpactID     uint       `json:"inherent_risk_matrix_impact_id" binding:"required" gorm:"not null"`
	InherentRiskMatrixRiskRatingID uint       `json:"inherent_risk_matrix_risk_rating_id" binding:"required" gorm:"not null"`
	ResidualRiskMatrixLikelihoodID *uint      `json:"residual_risk_matrix_likelihood_id" gorm:"default:null"`
	ResidualRiskMatrixImpactID     *uint      `json:"residual_risk_matrix_impact_id" gorm:"default:null"`
	ResidualRiskMatrixRiskRatingID *uint      `json:"residual_risk_matrix_risk_rating_id" gorm:"default:null"`
	CreatedBy                      uint       `json:"created_by" gorm:"not null"`
	CreatedAt                      time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                      *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	UpdatedBy                      *uint      `json:"updated_by" gorm:"default:null"`

	// Relationships
	RiskRegister                 *RiskRegister            `gorm:"foreignkey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	InherentRiskMatrixLikelihood *RiskMatrixLikelihood    `gorm:"foreignkey:InherentRiskMatrixLikelihoodID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	InherentRiskMatrixImpact     *RiskMatrixImpact        `gorm:"foreignkey:InherentRiskMatrixImpactID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	InherentRiskMatrixRiskRating *RiskMatrixRiskRating    `gorm:"foreignkey:InherentRiskMatrixRiskRatingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResidualRiskMatrixLikelihood *RiskMatrixLikelihood    `gorm:"foreignkey:ResidualRiskMatrixLikelihoodID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResidualRiskMatrixImpact     *RiskMatrixImpact        `gorm:"foreignkey:ResidualRiskMatrixImpactID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResidualRiskMatrixRiskRating *RiskMatrixRiskRating    `gorm:"foreignkey:ResidualRiskMatrixRiskRatingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskImpactEvaluations        *[]RiskImpactEvaluation  `gorm:"foreignkey:RiskEvaluationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskAppetitieAlignments      *[]RiskAppetiteAlignment `gorm:"foreignkey:RiskEvaluationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                *User                    `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser                *User                    `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by RiskEvaluation to `risk_evaluations`
func (RiskEvaluation) TableName() string {
	return "risk_evaluations"
}
