package models

import (
	"time"
)

type RiskUpdatePlanRiskWorkshopAssessmentDocument struct {
	ID                       uint   `gorm:"primary_key;AUTO_INCREMENT"`
	RiskUpdatePlanTimelineID uint   `gorm:"not null"`
	URL                      *string
	Name                     *string
	Size                     *int
	Token                    *string
	OriginalName             *string
	Type                     *string
	CreatedBYID              *uint
	UpdatedBYID              *uint
	CreatedAt                time.Time
	UpdatedAt                *time.Time

	// Relationships
	RiskUpdatePlanTimeline *RiskUpdatePlanTimeline `gorm:"foreignKey:RiskUpdatePlanTimelineID"`
	CreatedBy              *User                  `gorm:"foreignKey:CreatedBYID"`
	UpdatedBy              *User                  `gorm:"foreignKey:UpdatedBYID"`
}

func (RiskUpdatePlanRiskWorkshopAssessmentDocument) TableName() string {
	return "risk_update_plan_risk_workshop_assessment_documents"
}
