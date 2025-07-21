package models

import (
	"time"
)

type ProjectActionPlanStatusUpdate struct {
	ID                  uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProjectActionPlanID uint       `json:"project_action_plan_id" gorm:"not null"`
	ActionPlanStatusID  uint       `json:"action_plan_status_id" gorm:"not null"`
	Percentage          *int       `json:"percentage" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	Comment             string     `json:"comment" binding:"required" gorm:"type:nvarchar(max);not null"`
	CreatedAt           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID         uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID         *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	ProjectActionPlan *ProjectActionPlan                       `gorm:"foreignKey:ProjectActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus  *ActionPlanStatus                        `gorm:"foreignKey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser     *User                                    `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser     *User                                    `gorm:"foreignKey:UpdatedByID"`
	Documents         []*ProjectActionPlanStatusUpdateDocument `gorm:"foreignKey:ProjectActionPlanStatusUpdateID"`
}

// TableName overrides the table name used by ProjectActionPlanStatusUpdate to `project_action_plan_status_updates`
func (ProjectActionPlanStatusUpdate) TableName() string {
	return "project_action_plan_status_updates"
}
