package models

import (
	"time"
)

type ProjectImpactAssessment struct {
	ID                            uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProjectChangeRequestID        uint       `json:"project_change_request_id" gorm:"not null"`
	ProjectImpactAssessmentTypeID uint       `json:"project_impact_assessment_type_id" gorm:"not null"`
	CurrentDate                   *time.Time `json:"current_date" gorm:"type:date;null"`
	RevisedDate                   *time.Time `json:"revised_date" gorm:"type:date;null"`
	CurrentBudget                 *float64   `json:"current_budget" gorm:"type:decimal(18,2);null"`
	RevisedBudget                 *float64   `json:"revised_budget" gorm:"type:decimal(18,2);null"`
	CurrentResourceCount          *int       `json:"current_resource_count" gorm:"null"`
	RevisedResourceCount          *int       `json:"revised_resource_count" gorm:"null"`
	TechnicalComplexity           *string    `json:"technical_complexity" gorm:"type:nvarchar(max);null"`

	// Relationships
	ProjectChangeRequest        *ProjectChangeRequest        `gorm:"foreignKey:ProjectChangeRequestID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ProjectImpactAssessmentType *ProjectImpactAssessmentType `gorm:"foreignKey:ProjectImpactAssessmentTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectImpactAssessment) TableName() string {
	return "project_change_request_impact_assessments"
}
