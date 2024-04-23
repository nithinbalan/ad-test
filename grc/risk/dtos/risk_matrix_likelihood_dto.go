package dtos

type RiskMatrixLikelihoodDTO struct {
	ID            uint   `json:"id"`
	RiskMatrixID  uint   `json:"risk_matrix_id"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	Status        string `json:"status"`
}

type RiskMatrixLikelihoodFormDTO struct {
	ID            uint   `json:"id"`
	TitleEn       string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr       string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DescriptionEn string `json:"description_en" binding:"required" gorm:"type:nvarchar(max);not null"`
	DescriptionAr string `json:"description_ar" binding:"required" gorm:"type:nvarchar(max);not null"`
	Status        string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}
