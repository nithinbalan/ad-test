package models

import (
	"time"
)

type RiskReviewComment struct {
	ID                  uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	RiskRegisterID      uint       `json:"risk_register_id" gorm:"not null"`
	Comment             string     `json:"comment" binding:"required"`
	CommentStatusID     uint       `json:"comment_status_id"`
	Level               int        `json:"level"`
	RiskReviewCommentID *uint      `json:"risk_review_comment_id"`
	CreatedBYID         uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID         *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt           time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	RiskRegister  *RiskRegister        `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ParentComment *RiskReviewComment   `gorm:"foreignKey:RiskReviewCommentID"`
	ChildComment  []*RiskReviewComment `gorm:"foreignKey:RiskReviewCommentID"`
	CreatedByUser *User                `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User                `gorm:"foreignKey:UpdatedBYID"`
}

func (RiskReviewComment) TableName() string {
	return "risk_review_comments"
}
