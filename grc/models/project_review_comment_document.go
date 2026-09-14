package models

import (
	"time"
)

type ProjectReviewCommentDocument struct {
	ID                     uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectReviewCommentID uint      `json:"project_review_comment_id" gorm:"not null"`
	Name                   string    `json:"name" gorm:"type:nvarchar(255);not null"`
	Token                  string    `json:"token" gorm:"type:nvarchar(max);not null"`
	OriginalName           string    `json:"original_name" gorm:"type:nvarchar(255);not null"`
	Type                   string    `json:"type" gorm:"type:nvarchar(255);not null"`
	Size                   int64     `json:"size" gorm:"not null"`
	CreatedBYID            uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt              time.Time `json:"created_at" gorm:"type:datetime2;default:CURRENT_TIMESTAMP"`

	ProjectReviewComment *ProjectReviewComment `gorm:"foreignKey:ProjectReviewCommentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBy            *User                 `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectReviewCommentDocument) TableName() string {
	return "project_review_comment_documents"
}
