package models

import (
	"time"
)

type ComplianceAssessmentControlComment struct {
	ID                                      uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceAssessmentControlID           uint      `json:"compliance_assessment_control_id" gorm:"not null"`
	Comment                                 string    `json:"comment" binding:"required"`
	CreatedBYID                             uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                               time.Time `json:"created_at" gorm:"type:datetime;not null"`
	ComplianceAssessmentRequirementStatusID *uint     `json:"compliance_assessment_requirement_status_id" gorm:"default:null"`

	// Relationships
	ComplianceAssessmentControl           *ComplianceAssessmentControl           `gorm:"foreignKey:ComplianceAssessmentControlID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ComplianceAssessmentRequirementStatus *ComplianceAssessmentRequirementStatus `gorm:"foreignKey:ComplianceAssessmentRequirementStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                         *User                                  `gorm:"foreignKey:CreatedBYID"`
}

func (ComplianceAssessmentControlComment) TableName() string {
	return "compliance_assessment_control_comments"
}
