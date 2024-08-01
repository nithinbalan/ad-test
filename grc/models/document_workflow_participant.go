package models

import (
	"time"
)

// DocumentWorkflowParticipant represents a participant in a document workflow.
type DocumentWorkflowParticipant struct {
	ID                 uint `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DocumentWorkflowID uint `json:"document_workflow_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	//WorkflowLevelTypeID uint       `json:"workflow_level_type_id" binding:"required" gorm:"not null"`
	UserID      uint       `json:"user_id" binding:"required" gorm:"not null"`
	RoleID      *uint      `json:"role_id" gorm:"default:null"`
	CreatedBYID uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt   time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	DocumentWorkflow *DocumentWorkflow `gorm:"foreignKey:DocumentWorkflowID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	//WorkflowLevelType *WorkflowLevelType `gorm:"foreignKey:WorkflowLevelTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User      *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Role      *Role `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedBy *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by DocumentWorkflowParticipant to `document_workflow_participants`
func (DocumentWorkflowParticipant) TableName() string {
	return "document_workflow_participants"
}
