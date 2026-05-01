package models

import (
	"time"
)

type RiskResponseReviewCommentDocument struct {
	ID                             uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	RiskResponseReviewCommentID    uint       `json:"risk_response_review_comment_id" gorm:"not null"`
	URL                            *string    `json:"url"`
	Name                           string     `json:"name" binding:"required"`
	Size                           *int       `json:"size"`
	Token                          *string    `json:"token"`
	OriginalName                   *string    `json:"original_name"`
	Type                           *string    `json:"type"`
	CreatedByID                    uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID                    *uint      `json:"updated_by_id"`
	CreatedAt                      time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                      *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	ReviewComment *RiskResponseReviewComment `gorm:"foreignKey:RiskResponseReviewCommentID"`
	CreatedByUser *User                      `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser *User                      `gorm:"foreignKey:UpdatedByID"`
}

func (RiskResponseReviewCommentDocument) TableName() string {
	return "risk_response_review_comment_documents"
}
