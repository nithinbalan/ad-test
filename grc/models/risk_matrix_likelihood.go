package models

type RiskMatrixLikelihood struct {
	ID            uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskMatrixID  uint   `json:"risk_matrix_id" binding:"required" gorm:"not null"`
	TitleEn       string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr       string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Score         uint   `json:"score" binding:"required" gorm:"not null"`
	DescriptionEn string `json:"description_en" binding:"required" gorm:"type:nvarchar(max);not null"`
	DescriptionAr string `json:"description_ar" binding:"required" gorm:"type:nvarchar(max);not null"`
	Status        string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`

	// Relationships
	RiskMatrix *RiskMatrix `gorm:"foreignkey:RiskMatrixID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by RiskMatrixLikelihood to `risk_matrix_likelihoods`
func (RiskMatrixLikelihood) TableName() string {
	return "risk_matrix_likelihoods"
}
