package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type IncidentClosureFormDTO struct {
	CompletionDate    string    `json:"completion_date" binding:"required,datetime=2006-01-02"` // Date format: "2006-01-02"
	ResolutionSummary *string   `json:"resolution_summary" binding:"omitempty"`
	Notes             *string   `json:"notes" binding:"omitempty"`
	Lessons           []*string `json:"lessons" binding:"required,dive,max=500" gorm:"type:varchar(500);not null"`
}

type IncidentClosureParsedDates struct {
	CompletionDate *time.Time
}

type IncidentClosureDetailDTO struct {
	ID                   int                              `json:"id"`
	IncidentRegister     IncidentRegisterDetailDTO        `json:"incident_register"`
	CompletionDateString *string                          `json:"completion_date"`
	ResolutionSummary    *string                          `json:"resolution_summary"`
	Notes                *string                          `json:"notes"`
	Lessons              []*IncidentClosureLessonDTO      `json:"lessons"`
	CreatedAt            time.Time                        `json:"created_at"`
	UpdatedAt            *time.Time                       `json:"updated_at"`
	CreatedBYID          *uint                            `json:"created_by_id"`
	CreatedByUser        organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedBYID          *uint                            `json:"updated_by_id"`
	UpdatedByUser        organizationDTOs.UserMinimalDTO  `json:"updated_by"`
	SubmittedBYID        *uint                            `json:"submitted_by_id"`
	SubmittedByUser      *organizationDTOs.UserMinimalDTO `json:"submitted_by"`
	SubmittedAt          *time.Time                       `json:"submitted_at"`
	WorkflowStatusID     *uint                            `json:"workflow_status_id"`
	ApprovedBYID         *uint                            `json:"approved_by_id"`
	ApprovedByUser       *organizationDTOs.UserMinimalDTO `json:"approved_by"`
	ApprovedAt           *time.Time                       `json:"approved_at"`
	NextReviewLevel      *int                             `json:"next_review_level"`
	WorkflowStatus       *masterDTOs.WorkflowStatusDTO    `json:"workflow_status"`
}
