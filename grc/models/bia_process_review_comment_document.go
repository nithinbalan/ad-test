package models

import (
	"time"
)

type BiaProcessReviewCommentDocument struct {
	ID                        uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	BiaProcessReviewCommentID uint       `json:"bia_process_review_comment_id" gorm:"not null"`
	URL                       *string    `json:"url"`
	Name                      string     `json:"name" binding:"required"`
	Size                      *int       `json:"size"`
	Token                     *string    `json:"token"`
	OriginalName              *string    `json:"original_name"`
	Type                      *string    `json:"type"`
	CreatedByID               uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID               *uint      `json:"updated_by_id"`
	CreatedAt                 time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                 *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	ReviewComment *BiaProcessReviewComment `gorm:"foreignKey:BiaProcessReviewCommentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User                    `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser *User                    `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BiaProcessReviewCommentDocument) TableName() string {
	return "bia_process_review_comment_documents"
}
