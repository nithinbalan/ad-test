package models

import "time"

type IdentityReviewCampaignWorkflowReview struct {
	ID                               uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	IdentityReviewCampaignWorkflowID uint       `json:"identity_review_campaign_workflow_id" gorm:"not null;index"`
	WorkflowStatusID                 uint       `json:"workflow_status_id" gorm:"not null"`
	Comment                          string     `json:"comment" gorm:"type:text"`
	CreatedByID                      uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID                      *uint      `json:"updated_by_id"`
	CreatedAt                        time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                        *time.Time `json:"updated_at"`

	WorkflowStatus *WorkflowStatus `gorm:"foreignKey:WorkflowStatusID"`
	CreatedByUser  *User           `gorm:"foreignKey:CreatedByID"`
}

func (IdentityReviewCampaignWorkflowReview) TableName() string {
	return "identity_review_campaign_workflow_reviews"
}
