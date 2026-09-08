package models

import (
	"time"
)

type ProjectChangeRequestReviewCommentDocument struct {
	ID                                  uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectChangeRequestReviewCommentID uint      `json:"project_change_request_review_comment_id" gorm:"not null"`
	Name                                string    `json:"name" gorm:"type:nvarchar(255);not null"`
	Token                               string    `json:"token" gorm:"type:nvarchar(max);not null"`
	OriginalName                        string    `json:"original_name" gorm:"type:nvarchar(255);not null"`
	Type                                string    `json:"type" gorm:"type:nvarchar(255);not null"`
	Size                                int64     `json:"size" gorm:"not null"`
	CreatedBYID                         uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                           time.Time `json:"created_at" gorm:"type:datetime2;default:CURRENT_TIMESTAMP"`

	ProjectChangeRequestReviewComment *ProjectChangeRequestReviewComment `gorm:"foreignKey:ProjectChangeRequestReviewCommentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBy                         *User                              `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectChangeRequestReviewCommentDocument) TableName() string {
	return "project_change_request_review_comment_documents"
}
