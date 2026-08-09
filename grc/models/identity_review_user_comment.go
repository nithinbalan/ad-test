package models

import "time"

type IdentityReviewUserComment struct {
	ID                   uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	IdentityReviewUserID uint       `json:"identity_review_user_id" gorm:"not null"`
	AccountDecisionID    *uint      `json:"account_decision_id" gorm:"default:null"`
	Comment              string     `json:"comment" gorm:"type:nvarchar(max);not null"`
	CreatedBYID          uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID          *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	AccountDecision *AccountDecision `json:"account_decision" gorm:"foreignKey:AccountDecisionID"`
	CreatedByUser   *User            `json:"created_by" gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser   *User            `json:"updated_by" gorm:"foreignKey:UpdatedBYID"`
}

func (IdentityReviewUserComment) TableName() string {
	return "identity_review_user_comments"
}
