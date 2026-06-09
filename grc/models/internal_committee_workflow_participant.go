package models

import (
	"time"
)

// InternalCommitteeWorkflowParticipant represents a participant in internal committee workflow
type InternalCommitteeWorkflowParticipant struct {
	ID                          uint  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	InternalCommitteeWorkflowID uint  `json:"internal_committee_workflow_id" binding:"required" gorm:"not null;index"`
	UserID                      uint  `json:"user_id" binding:"required" gorm:"not null"`
	RoleID                      *uint `json:"role_id" gorm:"default:null"`
	CustomRoleID                *uint `json:"custom_role_id" gorm:"default:null"`
	//UnitID                      *uint      `json:"unit_id" gorm:"column:unit_id;default:null"`
	DelegateeUserID *uint      `json:"delegatee_user_id"`
	CreatedBYID     uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID     *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt       time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	//Relations
	InternalCommitteeWorkflow *InternalCommitteeWorkflow `gorm:"foreignKey:InternalCommitteeWorkflowID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User                      *User                      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DelegateeUser             *User                      `gorm:"foreignKey:DelegateeUserID"`
	Role                      *Role                      `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedBy                 *User                      `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy                 *User                      `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CustomRole                *WorkflowCustomRole        `gorm:"foreignKey:CustomRoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	//OrganizationUnit          *OrganizationUnit          `gorm:"foreignKey:UnitID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides table name
func (InternalCommitteeWorkflowParticipant) TableName() string {
	return "internal_committee_workflow_participants"
}
