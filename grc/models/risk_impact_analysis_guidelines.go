package models

type RiskImpactAnalysisGuideline struct {
	ID                          uint `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskImpactAnalysisID        uint `json:"risk_impact_analysis_id" binding:"required" gorm:"not null"`
	RiskMatrixImpactGuidelineID uint `json:"risk_matrix_impact_guideline_id" binding:"required" gorm:"not null"`

	// Relationships
	RiskImpactAnalysis        *RiskImpactAnalysis        `gorm:"foreignKey:RiskImpactAnalysisID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixImpactGuideline *RiskMatrixImpactGuideline `gorm:"foreignKey:RiskMatrixImpactGuidelineID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by RiskImpactAnalysisGuideline to `risk_impact_analysis_guidelines`
func (RiskImpactAnalysisGuideline) TableName() string {
	return "risk_impact_analysis_guidelines"
}
