package models

type ThreatIntelligenceRiskLevel struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	ColorCode string `json:"color_code" binding:"required,max=10" gorm:"type:varchar(10);not null"`
	Type      string `json:"type" binding:"required,max=10,oneof=low medium high extreme" gorm:"type:varchar(10);not null"`
	Status    string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`

	// Relationships
	ThreatIntelligenceRiskLevelValue []*ThreatIntelligenceRiskLevelValue `json:"values" gorm:"foreignKey:ThreatIntelligenceRiskLevelID"`
}

// TableName overrides the table name used by ThreatIntelligenceRiskLevel to `threat_intelligence_risk_levels`
func (ThreatIntelligenceRiskLevel) TableName() string {
	return "threat_intelligence_risk_levels"
}
