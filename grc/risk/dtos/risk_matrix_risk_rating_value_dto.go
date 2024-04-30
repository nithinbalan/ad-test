package dtos

type RiskMatrixRiskRatingValueFormDTO struct {
	CalculationType      string                          `json:"calculation_type" binding:"required,max=20,oneof=multiplication addition" gorm:"type:varchar(20);not null"`
	RiskMatrixRiskRating []*RiskMatrixRiskRatingValueDTO `json:"risk_matrix_risk_rating"`
}

type RiskMatrixRiskRatingValueDTO struct {
	ID                     uint   `json:"id"`
	RiskMatrixRiskRatingID uint   `json:"risk_matrix_risk_rating_id" binding:"required"`
	Values                 []uint `json:"values" binding:"required" gorm:"not null"`
}
