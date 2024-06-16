package models

import (
	"time"
)

type ComplianceRequirementStrategicObjective struct {
	ID                      uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRequirementID uint      `json:"compliance_requirement_id" binding:"required"`
	StrategicObjectiveID    uint      `json:"strategic_objective_id" binding:"required"`
	CreatedByID             uint      `json:"created_by_id" binding:"required"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceRequirement *ComplianceRequirement `gorm:"foreignKey:ComplianceRequirementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	StrategicObjective    *StrategicObjective    `gorm:"foreignKey:StrategicObjectiveID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedByID"`
}

func (ComplianceRequirementStrategicObjective) TableName() string {
	return "compliance_requirement_strategic_objectives"
}
