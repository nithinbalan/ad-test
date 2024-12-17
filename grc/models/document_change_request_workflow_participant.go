package models

import (
	"time"
)

// DocumentChangeRequestWorkflowParticipant represents a participant in a document change request workflow.
type DocumentChangeRequestWorkflowParticipant struct {
	ID                              uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DocumentChangeRequestWorkflowID uint       `json:"document_change_request_workflow_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UserID                          uint       `json:"user_id" binding:"required" gorm:"not null"`
	RoleID                          *uint      `json:"role_id" gorm:"default:null"`
	CreatedBYID                     uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID                     *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                       time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                       *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	DocumentChangeRequestWorkflow *DocumentChangeRequestWorkflow `gorm:"foreignKey:DocumentChangeRequestWorkflowID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User                          *User                          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Role                          *Role                          `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedBy                     *User                          `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy                     *User                          `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by DocumentChangeRequestWorkflowParticipant to `document_change_request_workflow_participants`
func (DocumentChangeRequestWorkflowParticipant) TableName() string {
	return "document_change_request_workflow_participants"
}
