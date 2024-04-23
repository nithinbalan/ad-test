package dtos

type RiskMatrixRiskRatingDTO struct {
	ID           uint   `json:"id"`
	RiskMatrixID uint   `json:"risk_matrix_id"`
	TitleEn      string `json:"title_en"`
	TitleAr      string `json:"title_ar"`
	ColorCode    string `json:"color_code"`
	Type         string `json:"type"`
	Status       string `json:"status"`
}

type RiskMatrixRiskRatingFormDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	ColorCode string `json:"color_code" binding:"required,max=10" gorm:"type:varchar(10);not null"`
	Type      string `json:"type" binding:"required,max=10,oneof=very_low low medium high very_high extream" gorm:"type:varchar(10);not null"`
	Status    string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}
