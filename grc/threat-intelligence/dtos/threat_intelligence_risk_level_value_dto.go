package dtos

type ThreatIntelligenceRiskLevelValueFormDTO struct {
	RiskLevels []*ThreatIntelligenceRiskLevelValueArrayDTO `json:"threat_intelligence_risk_levels"`
}

type ThreatIntelligenceRiskLevelValueArrayDTO struct {
	ID                            uint   `json:"id"`
	ThreatIntelligenceRiskLevelID uint   `json:"threat_intelligence_risk_level_id" binding:"required"`
	Values                        []uint `json:"values" binding:"required" gorm:"not null"`
}

type ThreatIntelligenceRiskLevelValueListResponseDTO struct {
	ID                            uint `json:"id"`
	ThreatIntelligenceRiskLevelID uint `json:"threat_intelligence_risk_level_id" binding:"required"`
	Value                         uint `json:"value" binding:"required" gorm:"not null"`
}

type ThreatIntelligenceRiskLevelDataDTO struct {
	ThreatIntelligenceRiskLevelID      uint   `json:"threat_intelligence_risk_level_id"`
	ThreatIntelligenceRiskLevelValueID uint   `json:"threat_intelligence_risk_level_value_id"`
	RiskLevelTitleEN                   string `json:"risk_level_title_en"`
	RiskLevelTitleAR                   string `json:"risk_level_title_ar"`
	RiskLevelType                      string `json:"risk_level_type"`
	ColorCode                          string `json:"color_code"`
	Value                              uint   `json:"value"`
}

type ThreatIntelligenceRiskLevelMatrixResponseDTO struct {
	ThreatIntelligenceLikelihoodID      uint                               `json:"threat_intelligence_likelihood_id"`
	ThreatIntelligenceLikelihoodTitleEN string                             `json:"threat_intelligence_likelihood_title_en"`
	ThreatIntelligenceLikelihoodTitleAR string                             `json:"threat_intelligence_likelihood_title_ar"`
	ThreatIntelligenceSeverityID        uint                               `json:"threat_intelligence_severity_id"`
	ThreatIntelligenceSeverityTitleEN   string                             `json:"threat_intelligence_severity_title_en"`
	ThreatIntelligenceSeverityTitleAR   string                             `json:"threat_intelligence_severity_title_ar"`
	ThreatIntelligenceScore             uint                               `json:"threat_intelligence_score"`
	ThreatIntelligenceRiskLevel         ThreatIntelligenceRiskLevelDataDTO `json:"threat_intelligence_risk_level"`
}

type ThreatIntelligenceRiskLevelValueDTO struct {
	ID                            uint `json:"id"`
	ThreatIntelligenceRiskLevelID uint `json:"threat_intelligence_risk_level_id" binding:"required"`
	Value                         uint `json:"value" binding:"required" gorm:"not null"`
}
