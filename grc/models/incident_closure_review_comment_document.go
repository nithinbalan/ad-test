package models

import (
	"time"
)

type IncidentClosureReviewCommentDocument struct {
	ID                            uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	IncidentClosureReviewCommentID uint      `json:"incident_closure_review_comment_id" gorm:"not null"`
	URL                           *string    `json:"url"`
	Name                          string     `json:"name" binding:"required"`
	Size                          *int       `json:"size"`
	Token                         *string    `json:"token"`
	OriginalName                  *string    `json:"original_name"`
	Type                          *string    `json:"type"`
	CreatedByID                   uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID                   *uint      `json:"updated_by_id"`
	CreatedAt                     time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                     *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	ReviewComment *IncidentClosureReviewComment `gorm:"foreignKey:IncidentClosureReviewCommentID"`
	CreatedByUser *User                         `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser *User                         `gorm:"foreignKey:UpdatedByID"`
}

func (IncidentClosureReviewCommentDocument) TableName() string {
	return "incident_closure_review_comment_documents"
}
