package models

import (
	"time"
)

type ProjectActionPlan struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProjectID          *uint      `json:"project_id" gorm:"null"`
	ProjectMilestoneID *uint      `json:"project_milestone_id" gorm:"null"`
	Title              string     `json:"title" gorm:"type:nvarchar(255);not null"`
	Description        *string    `json:"description" gorm:"type:nvarchar(max);null"`
	ResponsibleUserID  uint       `json:"responsible_user_id" binding:"required" gorm:"not null"`
	Comments           *string    `json:"comments" gorm:"type:nvarchar(max);null"`
	SeverityLevelID    *uint      `json:"severity_level_id" binding:"omitempty" gorm:"not null"`
	ReviewUserID       *uint      `json:"review_user_id" gorm:"default:null"`
	TargetDate         *time.Time `json:"target_date" binding:"omitempty" gorm:"type:date,default:null"`
	ActionPlanStatusID uint       `json:"action_plan_status_id" binding:"required" gorm:"not null"`
	Percentage         *int       `json:"percentage" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`

	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime"`
	CreatedBYID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	Project          *Project          `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ProjectMilestone *ProjectMilestone `gorm:"foreignKey:ProjectMilestoneID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser  *User             `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewUser       *User             `gorm:"foreignKey:ReviewUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel    *SeverityLevel    `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus *ActionPlanStatus `gorm:"foreignKey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	UpdatedByUser    *User             `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName overrides the table name used by ComplianceActionPlan to `project_action_plans`
func (ProjectActionPlan) TableName() string {
	return "project_action_plans"
}
