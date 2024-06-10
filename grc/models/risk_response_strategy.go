package models

type RiskResponseStrategy struct {
	ID            uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn       string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr       string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DescriptionEn string `json:"description_en" gorm:"type:text"`
	DescriptionAr string `json:"description_ar" gorm:"type:text"`
	Status        string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by RiskResponseStrategy to `risk_response_strategies`
func (RiskResponseStrategy) TableName() string {
	return "risk_response_strategies"
}
