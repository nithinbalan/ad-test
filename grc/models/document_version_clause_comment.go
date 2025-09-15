package models

import (
	"time"
)

type DocumentVersionClauseComment struct {
	ID                      uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	DocumentVersionClauseID uint      `json:"document_version_clause_id" gorm:"not null"`
	Comment                 string    `json:"comment" binding:"required"`
	CreatedBYID             uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	DocumentVersionClause *DocumentVersionClause `gorm:"foreignKey:DocumentVersionClauseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedBYID"`
}

func (DocumentVersionClauseComment) TableName() string {
	return "document_version_clause_comments"
}
