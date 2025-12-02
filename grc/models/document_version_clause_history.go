package models

import (
	"time"
)

type DocumentVersionClauseHistory struct {
	ID                      uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	DocumentVersionClauseID uint      `json:"document_version_clause_id" gorm:"not null"`
	OldValue                *string   `json:"old_value"`
	NewValue                *string   `json:"new_value"`
	ChangedField            *string   `json:"changed_field"`
	CreatedBYID             uint      `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt               time.Time `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	// Relationships
	DocumentVersionClause *DocumentVersionClause `gorm:"foreignKey:DocumentVersionClauseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedBYID"`
}

func (DocumentVersionClauseHistory) TableName() string {
	return "document_version_clause_histories"
}
