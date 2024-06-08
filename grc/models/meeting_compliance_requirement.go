package models

import (
	"time"
)

type MeetingComplianceRequirement struct {
	ID                      uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	MeetingID               uint      `json:"meeting_id" binding:"required"`
	ComplianceRequirementID uint      `json:"compliance_requirement_id" binding:"required"`
	CreatedBy               uint      `json:"created_by" binding:"required"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Meeting               *Meeting               `gorm:"foreignkey:MeetingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRequirement *ComplianceRequirement `gorm:"foreignkey:ComplianceRequirementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignkey:CreatedBy"`
}

func (MeetingComplianceRequirement) TableName() string {
	return "meeting_compliance_requirements"
}
