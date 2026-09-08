package models

import (
	"time"
)

type ProjectClosure struct {
	ID                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProjectID             uint       `json:"project_id" gorm:"not null"`
	OverallProjectOutcome *string    `json:"overall_project_outcome" gorm:"type:nvarchar(max);null"`
	ClosureReason         *string    `json:"closure_reason" gorm:"type:nvarchar(max);null"`
	ClosureComments       *string    `json:"closure_comments" gorm:"type:nvarchar(max);null"`
	LessonLearned         *string    `json:"lesson_learned" gorm:"type:nvarchar(max);null"`
	ApproverID            uint       `json:"approver_id" gorm:"not null"`
	CreatedByID           uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID           *uint      `json:"updated_by_id" gorm:"null"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime;null"`

	// Relationships
	Project       *Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ApproverUser  *User    `gorm:"foreignKey:ApproverID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User    `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser *User    `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectClosure) TableName() string {
	return "project_closures"
}
