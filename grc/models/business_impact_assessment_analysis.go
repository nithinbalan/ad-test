package models

import (
	"time"
)

type BusinessImpactAssessmentAnalysis struct {
	ID                         uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BusinessImpactAssessmentID uint      `json:"business_impact_assessment_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OutageBucketID             *uint     `json:"outage_bucket_id" gorm:"default:null;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ImpactAnalysisLevelID      *uint     `json:"impact_analysis_level_id" gorm:"default:null;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBYID                uint      `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt                  time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relations
	BusinessImpactAssessment *BusinessImpactAssessment `gorm:"foreignKey:BusinessImpactAssessmentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OutageBucket             *BCMBIAOutageBucket       `gorm:"foreignKey:OutageBucketID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ImpactAnalysisLevel      *BCMImpactAnalysisLevel   `gorm:"foreignKey:ImpactAnalysisLevelID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser            *User                     `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (BusinessImpactAssessmentAnalysis) TableName() string {
	return "business_impact_assessment_analysis"
}
