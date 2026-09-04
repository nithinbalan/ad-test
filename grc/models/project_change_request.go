package models

import (
	"time"
)

type ProjectChangeRequest struct {
	ID                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber       string     `json:"reference_number" gorm:"type:varchar(50);not null"`
	ProjectID             uint       `json:"project_id" gorm:"not null"`
	Title                 string     `json:"title" gorm:"type:nvarchar(255);not null"`
	Description           *string    `json:"description" gorm:"type:nvarchar(max);null"`
	BusinessJustification *string    `json:"business_justification" gorm:"type:nvarchar(max);null"`
	AnticipatedBenefits   *string    `json:"anticipated_benefits" gorm:"type:nvarchar(max);null"`
	RisksOfDeferral       *string    `json:"risks_of_deferral" gorm:"type:nvarchar(max);null"`
	RequestedByUserID     uint       `json:"requested_by_user_id" gorm:"not null"`
	RequestedOn           *time.Time `json:"requested_on" gorm:"type:date;not null"`
	WorkflowStatusID      *uint      `json:"workflow_status_id" gorm:"null"`
	SubmittedBYID         *uint      `json:"submitted_by_id" gorm:"null"`
	SubmittedAt           *time.Time `json:"submitted_at" gorm:"type:datetime;null"`
	NextReviewLevel       *int       `json:"next_review_level" gorm:"null"`
	ApprovedBYID          *uint      `json:"approved_by_id" gorm:"null"`
	ApprovedAt            *time.Time `json:"approved_at" gorm:"type:datetime;null"`
	RejectedBYID          *uint      `json:"rejected_by_id" gorm:"null"`
	RejectedAt            *time.Time `json:"rejected_at" gorm:"type:datetime;null"`
	CreatedBYID           uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID           *uint      `json:"updated_by_id" gorm:"null"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime;null"`

	// Relationships
	Project           *Project                               `gorm:"foreignKey:ProjectID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RequestedByUser   *User                                  `gorm:"foreignKey:RequestedByUserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus    *WorkflowStatus                        `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SubmittedByUser   *User                                  `gorm:"foreignKey:SubmittedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ApprovedByUser    *User                                  `gorm:"foreignKey:ApprovedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RejectedByUser    *User                                  `gorm:"foreignKey:RejectedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser     *User                                  `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:SET NULL"`
	UpdatedByUser     *User                                  `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:SET NULL"`
	ImpactAssessments []*ProjectImpactAssessment             `json:"impact_assessments" gorm:"foreignKey:ProjectChangeRequestID"`
	Documents         []*ProjectChangeRequestDocument        `json:"documents" gorm:"foreignKey:ProjectChangeRequestID"`
	Workflows         []*ProjectChangeRequestWorkflow        `json:"workflows" gorm:"foreignKey:ProjectChangeRequestID"`
	Histories         []*ProjectChangeRequestWorkflowHistory `json:"histories" gorm:"foreignKey:ProjectChangeRequestID"`
	ReviewComments    []*ProjectChangeRequestReviewComment   `json:"review_comments" gorm:"foreignKey:ProjectChangeRequestID"`
}

func (ProjectChangeRequest) TableName() string {
	return "project_change_requests"
}
