package models

import "time"

type KRIRegisterReviewComment struct {
	ID              uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	KRIRegisterID   uint       `json:"kri_register_id" gorm:"not null"`
	Comment         string     `json:"comment" binding:"required" gorm:"type:nvarchar(max)"`
	CommentStatusID uint       `json:"comment_status_id" gorm:"not null"`
	Level           *int       `json:"level"`
	ParentCommentID *uint      `json:"parent_comment_id"`
	IsEdited        bool       `json:"is_edited" gorm:"column:isEdited;not null;default:false"`
	CreatedBYID     uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID     *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt       time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	KRIRegister   *KRIRegister                `gorm:"foreignKey:KRIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ParentComment *KRIRegisterReviewComment   `gorm:"foreignKey:ParentCommentID"`
	ChildComments []*KRIRegisterReviewComment `gorm:"foreignKey:ParentCommentID"`
	CommentStatus *CommentStatus              `gorm:"foreignKey:CommentStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User                       `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User                       `gorm:"foreignKey:UpdatedBYID"`
}

func (KRIRegisterReviewComment) TableName() string {
	return "kri_register_review_comments"
}
