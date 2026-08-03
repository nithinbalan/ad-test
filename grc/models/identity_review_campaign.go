package models

import "time"

type IdentityReviewCampaign struct {
	ID                   uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber      string     `json:"reference_number" gorm:"type:varchar(45);not null"`
	CampaignName         string     `json:"campaign_name" gorm:"type:nvarchar(250);not null"`
	ReviewTypeID         uint       `json:"review_type_id" gorm:"not null"`
	Description          *string    `json:"description" gorm:"type:nvarchar(max)"`
	StartDate            time.Time  `json:"start_date" gorm:"type:date;not null"`
	DueDate              time.Time  `json:"due_date" gorm:"type:date;not null"`
	ReviewerID           uint       `json:"reviewer_id" gorm:"not null"`
	ApproverID           uint       `json:"approver_id" gorm:"not null"`
	CampaignStatusID     uint       `json:"campaign_status_id" gorm:"not null"`
	ReminderDaysOffset   *int       `json:"reminder_days_offset"`
	CompletionPercentage int        `json:"completion_percentage" gorm:"not null;default:0"`
	CreatedByID          uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID          *uint      `json:"updated_by_id"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime"`

	ReviewType     *ReviewType                       `json:"review_type" gorm:"foreignKey:ReviewTypeID"`
	CampaignStatus *ReviewCampaignStatus             `json:"campaign_status" gorm:"foreignKey:CampaignStatusID"`
	Reviewer       *User                             `json:"reviewer" gorm:"foreignKey:ReviewerID"`
	Approver       *User                             `json:"approver" gorm:"foreignKey:ApproverID"`
	CreatedByUser  *User                             `json:"created_by" gorm:"foreignKey:CreatedByID"`
	UpdatedByUser  *User                             `json:"updated_by" gorm:"foreignKey:UpdatedByID"`
	Systems        []*IdentityReviewCampaignSystem   `json:"systems" gorm:"foreignKey:IdentityReviewCampaignID"`
	Documents      []*IdentityReviewCampaignDocument `json:"documents" gorm:"foreignKey:IdentityReviewCampaignID"`
}

func (IdentityReviewCampaign) TableName() string {
	return "identity_review_campaigns"
}
