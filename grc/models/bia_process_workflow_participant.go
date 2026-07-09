package models

import (
	"time"
)

type BiaProcessWorkflowParticipant struct {
	ID                   uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BiaProcessWorkflowID uint       `json:"bia_process_workflow_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UserID               uint       `json:"user_id" binding:"required" gorm:"not null"`
	RoleID               *uint      `json:"role_id" gorm:"default:null"`
	DelegateeUserID      *uint      `json:"delegatee_user_id" gorm:"default:null"`
	CreatedBYID          uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID          *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt            time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	BiaProcessWorkflow *BiaProcessWorkflow `gorm:"foreignKey:BiaProcessWorkflowID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	User               *User               `gorm:"foreignKey:UserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	DelegateeUser      *User               `gorm:"foreignKey:DelegateeUserID"`
	Role               *Role               `gorm:"foreignKey:RoleID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBy          *User               `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedBy          *User               `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BiaProcessWorkflowParticipant) TableName() string {
	return "bia_process_workflow_participants"
}
