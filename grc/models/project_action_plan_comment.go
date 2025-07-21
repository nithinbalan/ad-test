package models

import (
	"time"
)

type ProjectActionPlanComment struct {
	ID                  uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectActionPlanID uint      `json:"project_action_plan_id" gorm:"not null"`
	Comment             string    `json:"comment" binding:"required"`
	CreatedBYID         uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt           time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ProjectActionPlan *ProjectActionPlan                  `gorm:"foreignKey:ProjectActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser     *User                               `gorm:"foreignKey:CreatedBYID"`
	Documents         []*ProjectActionPlanCommentDocument `gorm:"foreignKey:ProjectActionPlanCommentID"`
}

func (ProjectActionPlanComment) TableName() string {
	return "project_action_plan_comments"
}
