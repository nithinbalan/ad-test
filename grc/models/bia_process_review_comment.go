package models

import "time"

type BiaProcessReviewComment struct {
	ID              uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	BiaProcessID    uint      `json:"bia_process_id" gorm:"not null"`
	Comment         string    `json:"comment" binding:"required" gorm:"type:nvarchar(max)"`
	CommentStatusID uint      `json:"comment_status_id" gorm:"not null"`
	Level           *int      `json:"level"`
	ParentCommentID *uint     `json:"parent_comment_id"`
	IsEdited        bool      `json:"is_edited" gorm:"not null;default:false"`
	CreatedBYID     uint      `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID     *uint     `json:"updated_by_id" gorm:"default:null"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	BiaProcess    *BiaProcess                          `gorm:"foreignKey:BiaProcessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ParentComment *BiaProcessReviewComment             `gorm:"foreignKey:ParentCommentID"`
	ChildComments []*BiaProcessReviewComment           `gorm:"foreignKey:ParentCommentID"`
	CommentStatus *CommentStatus                       `gorm:"foreignKey:CommentStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User                                `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User                                `gorm:"foreignKey:UpdatedBYID"`
	Documents     []*BiaProcessReviewCommentDocument   `gorm:"foreignKey:BiaProcessReviewCommentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BiaProcessReviewComment) TableName() string {
	return "bia_process_review_comments"
}
