package models

type ThreatIntelligenceRiskLevelValue struct {
	ID                            uint `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ThreatIntelligenceRiskLevelID uint `json:"threat_intelligence_risk_level_id" binding:"required" gorm:"not null"`
	Value                         uint `json:"value" binding:"required" gorm:"not null"`

	// Relationships
	ThreatIntelligenceRiskLevel *ThreatIntelligenceRiskLevel `gorm:"foreignKey:ThreatIntelligenceRiskLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by ThreatIntelligenceRiskLevelValue to `threat_intelligence_risk_level_values`
func (ThreatIntelligenceRiskLevelValue) TableName() string {
	return "threat_intelligence_risk_level_values"
}
