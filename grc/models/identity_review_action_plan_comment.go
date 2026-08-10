package models

import "time"

type IdentityReviewActionPlanComment struct {
	ID                         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IdentityReviewActionPlanID uint      `json:"identity_review_action_plan_id" gorm:"not null"`
	Comment                    string    `json:"comment" gorm:"type:nvarchar(max);not null"`
	CreatedByID                uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                  time.Time `json:"created_at" gorm:"type:datetime;not null"`

	IdentityReviewActionPlan *IdentityReviewActionPlan                  `json:"identity_review_action_plan" gorm:"foreignKey:IdentityReviewActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser            *User                                      `json:"created_by" gorm:"foreignKey:CreatedByID"`
	Documents                []*IdentityReviewActionPlanCommentDocument `json:"documents" gorm:"foreignKey:IdentityReviewActionPlanCommentID"`
}

func (IdentityReviewActionPlanComment) TableName() string {
	return "identity_review_action_plan_comments"
}
