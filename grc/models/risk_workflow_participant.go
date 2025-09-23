package models

import (
	"time"
)

// RiskWorkflowParticipant represents a participant in a document workflow.
type RiskWorkflowParticipant struct {
	ID              uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskWorkflowID  uint       `json:"risk_workflow_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UserID          uint       `json:"user_id" binding:"required" gorm:"not null"`
	RoleID          *uint      `json:"role_id" gorm:"default:null"`
	DelegateeUserID *uint      `json:"delegatee_user_id"`
	CreatedBYID     uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID     *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt       time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	RiskWorkflow *RiskWorkflow `gorm:"foreignKey:RiskWorkflowID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	//WorkflowLevelType *WorkflowLevelType `gorm:"foreignKey:WorkflowLevelTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User          *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DelegateeUser *User `gorm:"foreignKey:DelegateeUserID"`
	Role          *Role `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedBy     *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy     *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by RiskWorkflowParticipant to `risk_workflow_participants`
func (RiskWorkflowParticipant) TableName() string {
	return "risk_workflow_participants"
}
