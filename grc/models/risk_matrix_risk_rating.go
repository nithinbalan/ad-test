package models

type RiskMatrixRiskRating struct {
	ID           uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskMatrixID uint   `json:"risk_matrix_id" binding:"required" gorm:"not null"`
	TitleEn      string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr      string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	ColorCode    string `json:"color_code" binding:"required,max=10" gorm:"type:varchar(10);not null"`
	Type         string `json:"type" binding:"required,max=10,oneof=very_low low medium high very_high" gorm:"type:varchar(10);not null"`
	Status       string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`

	// Relationships
	RiskMatrix                *RiskMatrix                  `gorm:"foreignkey:RiskMatrixID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixRiskRatingValue []*RiskMatrixRiskRatingValue `json:"values" gorm:"foreignkey:RiskMatrixRiskRatingID"`
}

// TableName overrides the table name used by RiskMatrixRiskRating to `risk_matrix_risk_ratings`
func (RiskMatrixRiskRating) TableName() string {
	return "risk_matrix_risk_ratings"
}
