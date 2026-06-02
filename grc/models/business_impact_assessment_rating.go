package models

import (
	"time"
)

type BusinessImpactAssessmentRating struct {
	ID                         uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProcessID                  uint      `json:"process_id" binding:"required" gorm:"column:process_id;not null"`
	OverallImpactLevelID       uint      `json:"overall_impact_level_id" binding:"required" gorm:"column:overall_impact_level_id;not null"`
	BCMBIAImpactLevelID        uint      `json:"bcm_bia_impact_level_id" binding:"required" gorm:"column:bcm_bia_impact_level_id;not null"`
	CreatedBYID                *uint     `json:"created_by_id" gorm:"column:created_by_id"`
	CreatedAt                  time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Process                  *Process                  `gorm:"foreignKey:ProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	OverallImpactLevel       *BCMImpactAnalysisLevel   `gorm:"foreignKey:OverallImpactLevelID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	BCMBIAImpactLevel        *BCMBIAImpactLevel        `gorm:"foreignKey:BCMBIAImpactLevelID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser            *User                     `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by BusinessImpactAssessmentRating to `business_impact_assessment_ratings`
func (BusinessImpactAssessmentRating) TableName() string {
	return "business_impact_assessment_ratings"
}
