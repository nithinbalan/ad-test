package models

import (
	"time"
)

type ControlTest struct {
	ID                            uint       `json:"id" gorm:"primaryKey;autoIncrement" binding:"-"`
	RiskRegisterID                *uint      `json:"risk_register_id" gorm:"default:null" binding:"omitempty"`
	ControlID                     uint       `json:"control_id" gorm:"not null" binding:"required"`
	ControlImplementationStatusID uint       `json:"control_implementation_status_id" gorm:"not null" binding:"required"`
	ControlEfficiencyStatusID     *uint      `json:"control_efficiency_status_id,omitempty" gorm:"default:null" binding:"omitempty"`
	Title                         string     `json:"title" gorm:"type:nvarchar(500);not null" binding:"required"`
	Description                   *string    `json:"description,omitempty" gorm:"type:nvarchar(MAX);default:null" binding:"omitempty"`
	Score                         *uint      `json:"score,omitempty" gorm:"default:null" binding:"omitempty"`
	ControlEffectivenessScore     *float64   `json:"control_effectiveness_score,omitempty" gorm:"type:decimal(3,2);default:null;check:control_effectiveness_score >= 0 AND control_effectiveness_score <= 5" binding:"omitempty"`
	ControlEffectivenessRatingID  *uint      `json:"control_effectiveness_rating_id,omitempty" gorm:"default:null" binding:"omitempty"`
	RiskReduction                 *float64   `json:"risk_reduction,omitempty" gorm:"type:decimal(5,2);default:null;check:risk_reduction >= 0 AND risk_reduction <= 100" binding:"omitempty"`
	ControlAutomationID           uint       `json:"control_automation_id" gorm:"default:null" binding:"required"`
	ControlConsistencyID          uint       `json:"control_consistency_id" gorm:"default:null" binding:"required"`
	Comments                      *string    `json:"comments,omitempty" gorm:"type:nvarchar(MAX);default:null" binding:"omitempty"`
	TestedSource                  string     `json:"tested_source" gorm:"type:varchar(45);not null" binding:"required"`
	TestedAt                      time.Time  `json:"tested_at,omitempty" gorm:"type:date;default:null" binding:"omitempty"`
	TestedBy                      uint       `json:"tested_by" gorm:"not null" binding:"required"`
	UpdatedAt                     *time.Time `json:"updated_at,omitempty" gorm:"type:datetime" binding:"omitempty"`
	UpdatedByID                   *uint      `json:"updated_by_id,omitempty" gorm:"default:null" binding:"omitempty"`
	RiskRegisterUpdatePlanID      *uint      `json:"risk_register_update_plan_id" gorm:"default:null"`

	// Foreign key relationships
	Risk                        *RiskRegister                `json:"risk" gorm:"foreignKey:RiskRegisterID"`
	Control                     *Control                     `json:"control" gorm:"foreignKey:ControlID"`
	ControlImplementationStatus *ControlImplementationStatus `json:"control_implementation_status" gorm:"foreignKey:ControlImplementationStatusID"`
	ControlEfficiencyStatus     *ControlEfficiencyStatus     `json:"control_efficiency_status" gorm:"foreignKey:ControlEfficiencyStatusID"`
	ControlEffectivenessRating  *ControlEffectivenessRating  `json:"control_effectiveness_rating" gorm:"foreignKey:ControlEffectivenessRatingID"`
	TestedByUser                *User                        `json:"tested_by_user" gorm:"foreignKey:TestedBy"`
	UpdatedByUser               *User                        `json:"updated_by_id_user" gorm:"foreignKey:UpdatedByID"`
	Documents                   []*ControlTestDocument       `json:"documents" gorm:"foreignKey:ControlTestID"`
	ControlAutomation           *ControlAutomation           `gorm:"foreignKey:ControlAutomationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ControlConsistency          *ControlConsistency          `gorm:"foreignKey:ControlConsistencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskRegisterUpdatePlan      *RiskRegisterUpdatePlan      `gorm:"foreignKey:RiskRegisterUpdatePlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by User to `control_tests`
func (ControlTest) TableName() string {
	return "control_tests"
}
