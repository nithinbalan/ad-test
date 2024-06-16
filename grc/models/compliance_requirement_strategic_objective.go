package models

import (
	"time"
)

type ComplianceRequirementStrategicObjective struct {
	ID                      uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRequirementID uint      `json:"compliance_requirement_id" binding:"required"`
	StrategicObjectiveID    uint      `json:"strategic_objective_id" binding:"required"`
	CreatedBy               uint      `json:"created_by" binding:"required"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceRequirement *ComplianceRequirement `gorm:"foreignkey:ComplianceRequirementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	StrategicObjective    *StrategicObjective    `gorm:"foreignkey:StrategicObjectiveID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignkey:CreatedBy"`
}

func (ComplianceRequirementStrategicObjective) TableName() string {
	return "compliance_requirement_strategic_objectives"
}
