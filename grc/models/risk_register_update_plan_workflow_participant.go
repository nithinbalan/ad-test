package models

import (
	"time"
)

type RiskRegisterUpdatePlanWorkflowParticipant struct {
	ID                               uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterUpdatePlanWorkflowID uint       `json:"risk_register_update_plan_workflow_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UserID                           uint       `json:"user_id" binding:"required" gorm:"not null"`
	RoleID                           *uint      `json:"role_id" gorm:"default:null"`
	DelegateeUserID                  *uint      `json:"delegatee_user_id"`
	CreatedBYID                      uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	WorkflowStatusID                 uint       `json:"workflow_status_id" binding:"required" gorm:"not null"`
	UpdatedBYID                      *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                        time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                        *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	RiskRegisterUpdatePlanWorkflow *RiskRegisterUpdatePlanWorkflow `gorm:"foreignKey:RiskRegisterUpdatePlanWorkflowID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User                           *User                           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DelegateeUser                  *User                           `gorm:"foreignKey:DelegateeUserID"`
	Role                           *Role                           `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedBy                      *User                           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy                      *User                           `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowStatus                 *WorkflowStatus                 `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (RiskRegisterUpdatePlanWorkflowParticipant) TableName() string {
	return "risk_update_plan_workflow_participants"
}
