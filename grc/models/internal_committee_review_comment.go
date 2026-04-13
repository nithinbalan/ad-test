package models

import (
	"time"
)

type InternalCommitteeReviewComment struct {
	ID                         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	InternalCommitteeID        uint       `json:"internal_committee_id" gorm:"not null"`
	InternalCommitteeVersionID uint       `json:"internal_committee_version_id" gorm:"not null"`
	Comment                    string     `json:"comment" binding:"required" gorm:"type:nvarchar(MAX);not null"`
	Level                      *int       `json:"level"`
	IsEdited                   bool       `json:"is_edited" gorm:"default:false"`
	CommentStatusID            uint       `json:"comment_status_id" gorm:"not null"`
	ParentCommentID            *uint      `json:"internal_committee_review_comment_id" gorm:"column:internal_committee_review_comment_id"` // keeping JSON same as DB column
	CreatedBYID                uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID                *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                  time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null;default:GETDATE()"`
	UpdatedAt                  *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	InternalCommittee        *InternalCommittee                        `gorm:"foreignKey:InternalCommitteeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	InternalCommitteeVersion *InternalCommitteeVersion                 `gorm:"foreignKey:InternalCommitteeVersionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ParentComment            *InternalCommitteeReviewComment           `gorm:"foreignKey:ParentCommentID"`
	ChildComments            []*InternalCommitteeReviewComment         `gorm:"foreignKey:ParentCommentID"`
	CommentStatus            *CommentStatus                            `gorm:"foreignKey:CommentStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Documents                []*InternalCommitteeReviewCommentDocument `gorm:"foreignKey:InternalCommitteeReviewCommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser            *User                                     `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser            *User                                     `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (InternalCommitteeReviewComment) TableName() string {
	return "internal_committee_review_comments"
}
