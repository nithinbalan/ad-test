package models

type RiskTreatmentGuideline struct {
	ID                     uint                 `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskMatrixRiskRatingID uint                 `json:"risk_matrix_risk_rating_id" gorm:"not null"`
	DescriptionEn          string               `json:"description_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DescriptionAr          string               `json:"description_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Status                 string               `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active'"`
	RiskMatrixRiskRating   RiskMatrixRiskRating `json:"risk_rating" gorm:"foreignKey:RiskMatrixRiskRatingID;references:ID"`
}

// TableName overrides the table name used by RiskTreatmentGuideline to `risk_treatment_guidelines`
func (RiskTreatmentGuideline) TableName() string {
	return "risk_treatment_guidelines"
}
