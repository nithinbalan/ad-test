package models

import "time"

type BiaProcess struct {
	ID               uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProcessID        uint       `json:"process_id" binding:"required" gorm:"not null"`
	WorkflowStatusID *uint      `json:"workflow_status_id" gorm:"default:null"`
	NextReviewLevel  *int       `json:"next_review_level" gorm:"default:null"`
	SubmittedByID    *uint      `json:"submitted_by_id" gorm:"default:null"`
	SubmittedAt      *time.Time `json:"submitted_at" gorm:"type:datetime;default:null"`
	ApprovedByID     *uint      `json:"approved_by_id" gorm:"default:null"`
	ApprovedAt       *time.Time `json:"approved_at" gorm:"type:datetime;default:null"`
	CreatedBYID      *uint      `json:"created_by_id" gorm:"default:null"`
	UpdatedBYID      *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt        time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt        *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	Process                     *Process                     `gorm:"foreignKey:ProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus              *WorkflowStatus              `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SubmittedBy                 *User                        `gorm:"foreignKey:SubmittedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ApprovedBy                  *User                        `gorm:"foreignKey:ApprovedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBy                   *User                        `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedBy                   *User                        `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	BiaProcessWorkflows         []*BiaProcessWorkflow        `gorm:"foreignKey:BiaProcessID"`
	BiaProcessReviewComments    []*BiaProcessReviewComment   `gorm:"foreignKey:BiaProcessID"`
	BiaProcessWorkflowHistories []*BiaProcessWorkflowHistory `gorm:"foreignKey:BiaProcessID"`
}

func (BiaProcess) TableName() string {
	return "bia_processes"
}
