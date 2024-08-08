package models

import (
	"time"
)

type DocumentControl struct {
	ID                      uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DocumentID              uint       `json:"document_id" gorm:"not null" binding:"required"`
	DocumentVersionID       *uint      `json:"document_version_id" gorm:"index"`
	DocumentVersionClauseID *uint      `json:"document_version_clause_id" gorm:"index"`
	ControlID               *uint      `json:"control_id" gorm:"not null" binding:"required"`
	CreatedBYID             uint       `json:"created_by_id" gorm:"not null" binding:"required"`
	CreatedAt               time.Time  `json:"created_at" gorm:"not null" binding:"required"`
	UpdatedBYID             *uint      `json:"updated_by_id"`
	UpdatedAt               *time.Time `json:"updated_at"`

	// Relationships
	Document              *Document              `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRIC"`
	DocumentVersion       *DocumentVersion       `gorm:"foreignKey:DocumentVersionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRIC"`
	DocumentVersionClause *DocumentVersionClause `gorm:"foreignKey:DocumentVersionClauseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRIC"`
	Control               *Control               `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRIC"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRIC"`
	UpdatedByUser         *User                  `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRIC"`
}

// TableName overrides the table name used by DocumentControl to `document_controls`
func (DocumentControl) TableName() string {
	return "document_controls"
}
