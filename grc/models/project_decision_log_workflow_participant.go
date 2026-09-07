package models

import (
	"time"
)

type ProjectDecisionLogWorkflowParticipant struct {
	ID                           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectDecisionLogWorkflowID uint       `json:"project_decision_log_workflow_id" gorm:"not null"`
	RoleID                       *uint      `json:"role_id" gorm:"null"`
	UserID                       uint       `json:"user_id" gorm:"not null"`
	CreatedBYID                  uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt                    time.Time  `json:"created_at" gorm:"type:datetime2;default:CURRENT_TIMESTAMP"`
	UpdatedAt                    *time.Time `json:"updated_at" gorm:"type:datetime2;null"`

	ProjectDecisionLogWorkflow *ProjectDecisionLogWorkflow `gorm:"foreignKey:ProjectDecisionLogWorkflowID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Role                       *Role                       `gorm:"foreignKey:RoleID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	User                       *User                       `gorm:"foreignKey:UserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBy                  *User                       `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectDecisionLogWorkflowParticipant) TableName() string {
	return "project_decision_log_workflow_participants"
}
