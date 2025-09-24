package models

import (
	"time"
)

// RiskWorkflow represents the structure of our resource
type RiskWorkflow struct {
	ID                  uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"` // Primary key, auto-incremented
	RiskRegisterID      uint       `json:"risk_register_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Level               int        `json:"level" binding:"required" gorm:"not null"`
	WorkflowLevelTypeID uint       `json:"workflow_level_type_id" binding:"required" gorm:"not null"`
	ApprovalTypes       string     `json:"approval_types" binding:"required,oneof='anyone' 'all'" gorm:"type:varchar(255);not null;default:'anyone';check:approval_types IN ('anyone', 'all')"`
	WorkflowStatusID    uint       `json:"workflow_status_id" binding:"required" gorm:"not null"`
	CreatedBYID         uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID         *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt           time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	RiskRegister      *RiskRegister      `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowLevelType *WorkflowLevelType `gorm:"foreignKey:WorkflowLevelTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowStatus    *WorkflowStatus    `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedBy         *User              `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy         *User              `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	RiskWorkflowParticipants []*RiskWorkflowParticipant `gorm:"foreignKey:RiskWorkflowID"`
	RiskWorkflowReviews      []*RiskWorkflowReview      `gorm:"foreignKey:RiskWorkflowID"`
}

// TableName overrides the table name used by User to `risk_workflows`
func (RiskWorkflow) TableName() string {
	return "risk_workflows"
}
