package models

import (
	"time"
)

type ProjectActionPlan struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProjectMilestoneID uint       `json:"project_milestone_id" binding:"required" gorm:"not null"`
	Title              string     `json:"title" gorm:"type:nvarchar(255);not null"`
	Description        *string    `json:"description" gorm:"type:nvarchar(max);null"`
	ResponsibleUserID  uint       `json:"responsible_user_id" binding:"required" gorm:"not null"`
	Comments           *string    `json:"comments" gorm:"type:nvarchar(max);null"`
	TargetDate         *time.Time `json:"target_date" binding:"omitempty" gorm:"type:date,default:null"`
	ActionPlanStatusID uint       `json:"action_plan_status_id" binding:"required" gorm:"not null"`
	Percentage         *int       `json:"percentage" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`

	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime"`
	CreatedByID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships

	ProjectMilestone ProjectMilestone  `gorm:"foreignKey:ProjectMilestoneID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser  *User             `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus *ActionPlanStatus `gorm:"foreignKey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser    *User             `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by ComplianceActionPlan to `project_action_plans`
func (ProjectActionPlan) TableName() string {
	return "project_action_plans"
}
