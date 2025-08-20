package models

import (
	"time"
)

type ComplianceAssessmentRequirementComment struct {
	ID                                               uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceAssessmentRequirementResponsibleUnitID uint      `json:"compliance_assessment_requirement_responsible_unit_id" gorm:"not null"`
	Comment                                          string    `json:"comment" binding:"required"`
	CreatedBYID                                      uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                                        time.Time `json:"created_at" gorm:"type:datetime;not null"`
	ComplianceAssessmentRequirementStatusID          *uint     `json:"compliance_assessment_requirement_status_id" gorm:"default:null"`

	// Relationships
	ComplianceAssessmentRequirementResponsibleUnit *ComplianceAssessmentRequirementResponsibleUnit `gorm:"foreignKey:ComplianceAssessmentRequirementResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ComplianceAssessmentRequirementStatus          *ComplianceAssessmentRequirementStatus          `gorm:"foreignKey:ComplianceAssessmentRequirementStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                                  *User                                           `gorm:"foreignKey:CreatedBYID"`
}

func (ComplianceAssessmentRequirementComment) TableName() string {
	return "compliance_assessment_requirement_comments"
}
