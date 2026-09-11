package models

type ThreatIntelligenceSeverity struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(500);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(500);not null"`
	Score   uint   `json:"score" binding:"required" gorm:"not null"`
	Status  string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by ThreatIntelligenceSeverity to `threat_intelligence_severities`
func (ThreatIntelligenceSeverity) TableName() string {
	return "threat_intelligence_severities"
}
