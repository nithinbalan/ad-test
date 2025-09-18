package models

import (
	"time"
)

type DocumentVersionClauseComment struct {
	ID                             uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	DocumentVersionClauseID        uint       `json:"document_version_clause_id" gorm:"not null"`
	Comment                        string     `json:"comment" binding:"required"`
	CommentStatusID                uint       `json:"comment_status_id"`
	Level                          int        `json:"level"`
	DocumentVersionClauseCommentID *uint      `json:"document_version_clause_comment_id"`
	CreatedBYID                    uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID                    *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                      time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                      *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	DocumentVersionClause *DocumentVersionClause          `gorm:"foreignKey:DocumentVersionClauseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ParentComment         *DocumentVersionClauseComment   `gorm:"foreignKey:DocumentVersionClauseCommentID"`
	ChildComment          []*DocumentVersionClauseComment `gorm:"foreignKey:DocumentVersionClauseCommentID"`
	CreatedByUser         *User                           `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser         *User                           `gorm:"foreignKey:UpdatedBYID"`
}

func (DocumentVersionClauseComment) TableName() string {
	return "document_version_clause_comments"
}
