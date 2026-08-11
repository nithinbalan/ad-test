package models

import "time"

type IdentityReviewUserAccountDecisionHistory struct {
	ID                   uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IdentityReviewUserID uint      `json:"identity_review_user_id" gorm:"not null;index"`
	AccountDecisionID    uint      `json:"account_decision_id" gorm:"not null"`
	CreatedByID          uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	IdentityReviewUser *IdentityReviewUser `json:"identity_review_user,omitempty" gorm:"foreignKey:IdentityReviewUserID"`
	AccountDecision    *AccountDecision    `json:"account_decision" gorm:"foreignKey:AccountDecisionID"`
	CreatedByUser      *User               `json:"created_by" gorm:"foreignKey:CreatedByID"`
}

func (IdentityReviewUserAccountDecisionHistory) TableName() string {
	return "identity_review_user_account_decision_histories"
}
