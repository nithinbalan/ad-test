package models

import "time"

type IdentityReviewCampaignWorkflow struct {
	ID                       uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	IdentityReviewCampaignID uint       `json:"identity_review_campaign_id" gorm:"not null;index"`
	Level                    int        `json:"level" gorm:"not null"`
	WorkflowLevelTypeID      uint       `json:"workflow_level_type_id" gorm:"not null"`
	ApprovalTypes            string     `json:"approval_types" gorm:"type:varchar(255);not null;default:'anyone'"`
	WorkflowStatusID         uint       `json:"workflow_status_id" gorm:"not null"`
	CreatedByID              uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at"`

	IdentityReviewCampaign *IdentityReviewCampaign `gorm:"foreignKey:IdentityReviewCampaignID"`
	WorkflowLevelType      *WorkflowLevelType      `gorm:"foreignKey:WorkflowLevelTypeID"`
	WorkflowStatus         *WorkflowStatus         `gorm:"foreignKey:WorkflowStatusID"`
	CreatedBy              *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy              *User                   `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	IdentityReviewCampaignWorkflowParticipants []*IdentityReviewCampaignWorkflowParticipant `gorm:"foreignKey:IdentityReviewCampaignWorkflowID"`
	IdentityReviewCampaignWorkflowReviews      []*IdentityReviewCampaignWorkflowReview      `gorm:"foreignKey:IdentityReviewCampaignWorkflowID"`
}

func (IdentityReviewCampaignWorkflow) TableName() string {
	return "identity_review_campaign_workflows"
}
