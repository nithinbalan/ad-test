package models

import "time"

type ComplianceLibraryClauseControl struct {
	ID                        uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceLibraryClauseID uint       `json:"compliance_library_clause_id" gorm:"not null"`
	ControlLibraryID          uint       `json:"control_library_id" gorm:"not null"`
	CreatedAt                 time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                 *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relations
	Clause  *ComplianceLibraryClause `json:"clause" gorm:"foreignKey:ComplianceLibraryClauseID;references:ID"`
	Control *ControlLibrary          `json:"control" gorm:"foreignKey:ControlLibraryID;references:ID"`
}

func (ComplianceLibraryClauseControl) TableName() string {
	return "compliance_library_clauses_controls"
}
