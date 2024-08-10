package models

import (
	"time"
)

type DocumentVersionClause struct {
	ID                      uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DocumentVersionID       uint       `json:"document_version_id" gorm:"not null;index" binding:"required"`
	DocumentVersionClauseID *uint      `json:"document_version_clause_id" gorm:"index; default:null"`
	ClauseNumber            *string    `json:"clause_number" gorm:"size:45;default:null"`
	ClauseTitleEn           string     `json:"clause_title_en" gorm:"type:nvarchar(255);size:255;not null" binding:"required"`
	ClauseTitleAr           string     `json:"clause_title_ar" gorm:"type:nvarchar(255);size:255;not null" binding:"required"`
	Description             *string    `json:"description" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID             uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID             *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt               time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt               *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	DocumentVersion  *DocumentVersion         `gorm:"foreignKey:DocumentVersionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ParantClause     *DocumentVersionClause   `gorm:"foreignKey:DocumentVersionClauseID"`
	ChildClauses     []*DocumentVersionClause `gorm:"foreignKey:DocumentVersionClauseID"`
	CreatedByUser    *User                    `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRIC"`
	UpdatedByUser    *User                    `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRIC"`
	DocumentControls []*DocumentControl       `gorm:"foreignKey:DocumentVersionClauseID"`
}

// TableName overrides the table name used by DocumentVersionClause to `document_version_clauses`
func (DocumentVersionClause) TableName() string {
	return "document_version_clauses"
}
