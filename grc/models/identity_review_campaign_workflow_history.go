package models

import "time"

type IdentityReviewCampaignWorkflowHistory struct {
	ID                       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IdentityReviewCampaignID uint      `json:"identity_review_campaign_id" gorm:"not null;index"`
	Level                    int       `json:"level" gorm:"not null"`
	WorkflowStatusID         uint      `json:"workflow_status_id" gorm:"not null"`
	Comment                  string    `json:"comment" gorm:"type:text"`
	CreatedByID              uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`

	WorkflowStatus *WorkflowStatus `gorm:"foreignKey:WorkflowStatusID"`
	CreatedByUser  *User           `gorm:"foreignKey:CreatedByID"`
}

func (IdentityReviewCampaignWorkflowHistory) TableName() string {
	return "identity_review_campaign_workflow_histories"
}
