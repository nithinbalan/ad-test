package models

import (
	"time"
)

type IncidentClosure struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentRegisterID uint       `json:"incident_register_id" binding:"required" gorm:"not null"`
	CompletionDate     *time.Time `json:"completion_date" binding:"required" gorm:"type:date;not null"`
	ResolutionSummary  *string    `json:"resolution_summary" gorm:"type:nvarchar(max);null"`
	Notes              *string    `json:"notes" gorm:"type:nvarchar(max);null"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime"`
	CreatedBYID        uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID        *uint      `json:"updated_by_id" gorm:"default:null"`
	SubmittedBYID      *uint      `json:"submitted_by_id" gorm:"default:null"`
	ApprovedBYID       *uint      `json:"approved_by_id" gorm:"default:null"`
	SubmittedAt        *time.Time `json:"submitted_at" gorm:"type:datetime;default:null"`
	ApprovedAt         *time.Time `json:"approved_at" gorm:"type:datetime;default:null"`
	NextReviewLevel    *int       `json:"next_review_level" gorm:"default:null"`
	WorkflowStatusID   *uint      `json:"workflow_status_id" gorm:"default:null"`

	// Relationships
	IncidentRegister         *IncidentRegister          `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Lessons                  []*IncidentClosureLesson   `gorm:"foreignKey:IncidentClosureID"`
	IncidentClosureWorkflows []*IncidentClosureWorkflow `gorm:"foreignKey:IncidentClosureID"`
	CreatedByUser            *User                      `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser            *User                      `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SubmittedByUser          *User                      `gorm:"foreignKey:SubmittedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ApprovedByUser           *User                      `gorm:"foreignKey:ApprovedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus           *WorkflowStatus            `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by IncidentClosure to `incident_closures`
func (IncidentClosure) TableName() string {
	return "incident_closures"
}
