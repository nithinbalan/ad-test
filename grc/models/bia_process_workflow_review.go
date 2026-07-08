package models

import (
	"time"
)

type BiaProcessWorkflowReview struct {
	ID                   uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BiaProcessWorkflowID uint       `json:"bia_process_workflow_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatusID     uint       `json:"workflow_status_id" binding:"required" gorm:"not null"`
	Comment              string     `json:"comment" gorm:"type:nvarchar(max)"`
	CreatedBYID          uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID          *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt            time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	BiaProcessWorkflow *BiaProcessWorkflow `gorm:"foreignKey:BiaProcessWorkflowID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus     *WorkflowStatus     `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser      *User               `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BiaProcessWorkflowReview) TableName() string {
	return "bia_process_workflow_reviews"
}
