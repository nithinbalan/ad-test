package models

type RiskMatrixRiskRatingValue struct {
	ID                     uint `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskMatrixRiskRatingID uint `json:"risk_matrix_risk_rating_id" binding:"required" gorm:"not null"`
	Value                  uint `json:"value" binding:"required" gorm:"not null"`

	// Relationships
	RiskMatrixRiskRating *RiskMatrixRiskRating `gorm:"foreignkey:RiskMatrixRiskRatingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by RiskMatrixRiskRatingValue to `risk_matrix_risk_rating_values`
func (RiskMatrixRiskRatingValue) TableName() string {
	return "risk_matrix_risk_rating_values"
}
