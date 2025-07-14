package models

import (
	"time"
)

type DocumentClause struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DocumentID           uint      `json:"document_id" binding:"required"`
	ComplianceRegisterID uint      `json:"compliance_register_id" binding:"required"`
	ComplianceClauseID   uint      `json:"compliance_clause_id" binding:"required"`
	CreatedByID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Document           *Document           `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRegister *ComplianceRegister `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceClause   *ComplianceClause   `gorm:"foreignKey:ComplianceClauseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedByID"`
}

func (DocumentClause) TableName() string {
	return "document_clauses"
}
