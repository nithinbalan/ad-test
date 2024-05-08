package models

import (
	"time"
)

type IncidentComplianceRequirement struct {
	ID                      uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentRegisterID      uint      `json:"incident_register_id" binding:"required"`
	ComplianceRequirementID uint      `json:"compliance_requirement_id" binding:"required"`
	CreatedBy               uint      `json:"created_by" binding:"required"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	IncidentRegister      *IncidentRegister      `gorm:"foreignkey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRequirement *ComplianceRequirement `gorm:"foreignkey:ComplianceRequirementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignkey:CreatedBy"`
}

func (IncidentComplianceRequirement) TableName() string {
	return "incident_compliance_requirements"
}
