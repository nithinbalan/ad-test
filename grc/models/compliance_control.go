package models

import (
	"time"
)

type ComplianceControl struct {
	ID                      uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRegisterID    uint      `json:"compliance_register_id" binding:"required"`
	ComplianceRequirementID uint      `json:"compliance_requirement_id"`
	ControlID               uint      `json:"control_id" binding:"required"`
	CreatedBy               uint      `json:"created_by" binding:"required"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceRegister    *ComplianceRegister    `gorm:"foreignkey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRequirement *ComplianceRequirement `gorm:"foreignkey:ComplianceRequirementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control               *Control               `gorm:"foreignkey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignkey:CreatedBy"`
}

func (ComplianceControl) TableName() string {
	return "compliance_controls"
}
