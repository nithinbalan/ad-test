package models

import (
	"time"
)

type IdentityReviewCampaignReviewComment struct {
	ID                                    uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	IdentityReviewCampaignID              uint       `json:"identity_review_campaign_id" gorm:"not null"`
	Comment                               string     `json:"comment" binding:"required"`
	CommentStatusID                       uint       `json:"comment_status_id"`
	Level                                 int        `json:"level"`
	IdentityReviewCampaignReviewCommentID *uint      `json:"identity_review_campaign_review_comment_id"`
	CreatedBYID                           uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID                           *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                             time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                             *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	IdentityReviewCampaign *IdentityReviewCampaign                `gorm:"foreignKey:IdentityReviewCampaignID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ParentComment          *IdentityReviewCampaignReviewComment   `gorm:"foreignKey:IdentityReviewCampaignReviewCommentID"`
	ChildComment           []*IdentityReviewCampaignReviewComment `gorm:"foreignKey:IdentityReviewCampaignReviewCommentID"`
	CreatedByUser          *User                                  `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser          *User                                  `gorm:"foreignKey:UpdatedBYID"`
}

func (IdentityReviewCampaignReviewComment) TableName() string {
	return "identity_review_campaign_review_comments"
}
