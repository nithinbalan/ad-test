package models

import "time"

type WorkflowLevel struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Level              int        `json:"level" binding:"required" gorm:"not null"`
	WorkflowID         uint       `json:"workflow_id" binding:"required" gorm:"not null"`
	WorkflowRoleTypeID uint       `json:"workflow_role_type_id" binding:"required" gorm:"not null"`
	ApprovalType       string     `json:"approval_type" binding:"required,oneof=anyone all" gorm:"type:varchar(50);not null;check:approval_type IN ('anyone','all')"`
	CreatedByID        uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID        *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	Workflow         *Workflow                   `gorm:"foreignKey:WorkflowID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowRoleType *WorkflowRoleType           `gorm:"foreignKey:WorkflowRoleTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User                       `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser    *User                       `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Participants     []*WorkflowLevelParticipant `gorm:"foreignKey:WorkflowLevelID"`
}

func (WorkflowLevel) TableName() string {
	return "workflow_levels"
}
