package models

import (
	"time"
)

type ProjectDecisionLogReviewComment struct {
	ID                   uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Comment              string     `json:"comment" gorm:"type:nvarchar(max);not null"`
	IsEdited             bool       `json:"is_edited" gorm:"default:0;not null"`
	CommentStatusID      uint       `json:"comment_status_id" gorm:"not null"`
	ProjectDecisionLogID uint       `json:"project_decision_log_id" gorm:"not null"`
	Level                *int       `json:"level" gorm:"null"`
	ParentCommentID      *uint      `json:"parent_comment_id" gorm:"null"`
	CreatedBYID          uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID          *uint      `json:"updated_by_id" gorm:"null"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime;null"`

	CommentStatus      *CommentStatus                             `gorm:"foreignKey:CommentStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ProjectDecisionLog *ProjectDecisionLog                        `gorm:"foreignKey:ProjectDecisionLogID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ParentComment      *ProjectDecisionLogReviewComment           `gorm:"foreignKey:ParentCommentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Children           []*ProjectDecisionLogReviewComment         `json:"children,omitempty" gorm:"foreignKey:ParentCommentID"`
	CreatedByUser      *User                                      `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser      *User                                      `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Documents          []*ProjectDecisionLogReviewCommentDocument `json:"documents" gorm:"foreignKey:ProjectDecisionLogReviewCommentID"`
}

func (ProjectDecisionLogReviewComment) TableName() string {
	return "project_decision_log_review_comments"
}
