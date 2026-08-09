package models

import "time"

type IdentityReviewCampaignWorkflowParticipant struct {
	ID                               uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	IdentityReviewCampaignWorkflowID uint       `json:"identity_review_campaign_workflow_id" gorm:"not null;index"`
	UserID                           uint       `json:"user_id" gorm:"not null"`
	RoleID                           *uint      `json:"role_id"`
	DelegateeUserID                  *uint      `json:"delegatee_user_id"`
	WorkflowStatusID                 *uint      `json:"workflow_status_id"`
	CreatedByID                      uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID                      *uint      `json:"updated_by_id"`
	CreatedAt                        time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                        *time.Time `json:"updated_at"`

	User           *User           `gorm:"foreignKey:UserID"`
	Role           *Role           `gorm:"foreignKey:RoleID"`
	WorkflowStatus *WorkflowStatus `gorm:"foreignKey:WorkflowStatusID"`
}

func (IdentityReviewCampaignWorkflowParticipant) TableName() string {
	return "identity_review_campaign_workflow_participants"
}
