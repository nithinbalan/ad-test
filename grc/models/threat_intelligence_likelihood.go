package models

type ThreatIntelligenceLikelihood struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(500);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(500);not null"`
	Score   uint   `json:"score" binding:"required" gorm:"not null"`
	Status  string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by ThreatIntelligenceLikelihood to `threat_intelligence_likelihoods`
func (ThreatIntelligenceLikelihood) TableName() string {
	return "threat_intelligence_likelihoods"
}
