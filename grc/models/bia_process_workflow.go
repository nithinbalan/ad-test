package models

import (
	"time"
)

type BiaProcessWorkflow struct {
	ID                  uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BiaProcessID        uint       `json:"bia_process_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Level               int        `json:"level" binding:"required" gorm:"not null"`
	WorkflowLevelTypeID uint       `json:"workflow_level_type_id" binding:"required" gorm:"not null"`
	ApprovalTypes       string     `json:"approval_types" binding:"required,oneof=anyone all" gorm:"type:varchar(255);not null;default:'anyone';check:approval_types IN ('anyone', 'all')"`
	WorkflowStatusID    uint       `json:"workflow_status_id" binding:"required" gorm:"not null"`
	CreatedBYID         uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID         *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt           time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	BiaProcess        *BiaProcess        `gorm:"foreignKey:BiaProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowLevelType *WorkflowLevelType `gorm:"foreignKey:WorkflowLevelTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus    *WorkflowStatus    `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBy         *User              `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedBy         *User              `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`

	BiaProcessWorkflowParticipants []*BiaProcessWorkflowParticipant `gorm:"foreignKey:BiaProcessWorkflowID"`
	BiaProcessWorkflowReviews      []*BiaProcessWorkflowReview      `gorm:"foreignKey:BiaProcessWorkflowID"`
}

func (BiaProcessWorkflow) TableName() string {
	return "bia_process_workflows"
}
