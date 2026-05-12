package models

import (
	"time"
)

type RiskRegisterUpdatePlanWorkflowHistory struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterUpdatePlanID uint       `json:"risk_register_update_plan_id" binding:"required" gorm:"not null;index"`
	Level                    int        `json:"level" binding:"required" gorm:"not null"`
	WorkflowStatusID         uint       `json:"workflow_status_id" binding:"required" gorm:"not null"`
	Comment                  string     `json:"comment" gorm:"type:text"`
	CreatedBYID              uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations - No FK to workflow to avoid deletion conflicts
	RiskRegisterUpdatePlan *RiskRegisterUpdatePlan `gorm:"foreignKey:RiskRegisterUpdatePlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	WorkflowStatus         *WorkflowStatus         `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (RiskRegisterUpdatePlanWorkflowHistory) TableName() string {
	return "risk_update_plan_workflow_histories"
}
