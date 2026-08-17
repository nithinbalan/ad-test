package models

import (
	"time"
)

type IncidentClosureReviewComment struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	IncidentClosureID uint       `json:"incident_closure_id" gorm:"not null"`
	Comment           string     `json:"comment" binding:"required" gorm:"type:nvarchar(max)"`
	Level             *int       `json:"level"`
	CommentStatusID   uint       `json:"comment_status_id" gorm:"not null"`
	IsEdited          bool       `json:"is_edited" gorm:"not null;default:false"`
	ParentCommentID   *uint      `json:"incident_closure_review_comment_id" gorm:"column:incident_closure_review_comment_id"`
	CreatedBYID       uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID       *uint      `json:"updated_by_id"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	IncidentClosure *IncidentClosure                        `gorm:"foreignKey:IncidentClosureID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ParentComment   *IncidentClosureReviewComment           `gorm:"foreignKey:ParentCommentID"`
	ChildComments   []*IncidentClosureReviewComment         `gorm:"foreignKey:ParentCommentID"`
	CommentStatus   *CommentStatus                          `gorm:"foreignKey:CommentStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User                                   `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser   *User                                   `gorm:"foreignKey:UpdatedBYID"`
	Documents       []*IncidentClosureReviewCommentDocument `gorm:"foreignKey:IncidentClosureReviewCommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (IncidentClosureReviewComment) TableName() string {
	return "incident_closure_review_comments"
}
