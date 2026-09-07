package models

import (
	"time"
)

type ProjectDecisionLog struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber   string     `json:"reference_number" gorm:"type:varchar(50);not null"`
	ProjectID         uint       `json:"project_id" gorm:"not null"`
	Decision          string     `json:"decision" gorm:"type:nvarchar(max);not null"`
	DecisionNotes     *string    `json:"decision_notes" gorm:"type:nvarchar(max);null"`
	Reason            *string    `json:"reason" gorm:"type:nvarchar(max);null"`
	RequestedByUserID uint       `json:"requested_by_user_id" gorm:"not null"`
	DecisionDate      time.Time  `json:"decision_date" gorm:"type:date;not null"`
	WorkflowStatusID  *uint      `json:"workflow_status_id" gorm:"null"`
	SubmittedBYID     *uint      `json:"submitted_by_id" gorm:"null"`
	SubmittedAt       *time.Time `json:"submitted_at" gorm:"type:datetime;null"`
	NextReviewLevel   *int       `json:"next_review_level" gorm:"null"`
	ApprovedBYID      *uint      `json:"approved_by_id" gorm:"null"`
	ApprovedAt        *time.Time `json:"approved_at" gorm:"type:datetime;null"`
	RejectedBYID      *uint      `json:"rejected_by_id" gorm:"null"`
	RejectedAt        *time.Time `json:"rejected_at" gorm:"type:datetime;null"`
	CreatedByID       uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID       *uint      `json:"updated_by_id" gorm:"null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;null"`

	// Relationships
	Project         *Project                             `gorm:"foreignKey:ProjectID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RequestedByUser *User                                `gorm:"foreignKey:RequestedByUserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus  *WorkflowStatus                      `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SubmittedByUser *User                                `gorm:"foreignKey:SubmittedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ApprovedByUser  *User                                `gorm:"foreignKey:ApprovedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RejectedByUser  *User                                `gorm:"foreignKey:RejectedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser   *User                                `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:SET NULL"`
	UpdatedByUser   *User                                `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:NO ACTION,OnDelete:SET NULL"`
	Documents       []*ProjectDecisionLogDocument        `json:"documents" gorm:"foreignKey:ProjectDecisionLogID"`
	Workflows       []*ProjectDecisionLogWorkflow        `json:"workflows" gorm:"foreignKey:ProjectDecisionLogID"`
	Histories       []*ProjectDecisionLogWorkflowHistory `json:"histories" gorm:"foreignKey:ProjectDecisionLogID"`
	ReviewComments  []*ProjectDecisionLogReviewComment   `json:"review_comments" gorm:"foreignKey:ProjectDecisionLogID"`
}

func (ProjectDecisionLog) TableName() string {
	return "project_decision_logs"
}
