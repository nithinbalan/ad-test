package models

import "time"

type BcpStrategicObjective struct {
	ID                       uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	StrategicObjectiveID     uint      `json:"strategic_objective_id" gorm:"not null"`
	BusinessContinuityPlanID uint      `json:"business_continuity_plan_id" gorm:"not null"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`
	CreatedBYID              uint      `json:"created_by_id" gorm:"not null"`

	// Relationships
	StrategicObjective     *StrategicObjective     `gorm:"foreignKey:StrategicObjectiveID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	BusinessContinuityPlan *BusinessContinuityPlan `gorm:"foreignKey:BusinessContinuityPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID"`
}

// TableName overrides the default table name
func (BcpStrategicObjective) TableName() string {
	return "bcp_strategic_objectives"
}
