package models

import "time"

type ComplianceLibraryClause struct {
	ID                         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceLibraryVersionID uint       `json:"compliance_library_version_id"`
	ClauseNumber               string     `json:"clause_number" gorm:"type:nvarchar(100);not null"`
	ClauseTitleEn              string     `json:"clause_title_en" gorm:"type:nvarchar(255);not null"`
	ClauseTitleAr              *string    `json:"clause_title_ar" gorm:"type:nvarchar(255)"`
	ClauseDescriptionEn        *string    `json:"clause_description_en" gorm:"type:nvarchar(max)"`
	ClauseDescriptionAr        *string    `json:"clause_description_ar" gorm:"type:nvarchar(max)"`
	ComplianceLibraryClauseID  *uint      `json:"compliance_library_clause_id"`
	CreatedAt                  time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                  *time.Time `json:"updated_at" gorm:"type:datetime"`
	CreatedByID                uint       `json:"created_by_id"`
	UpdatedByID                *uint      `json:"updated_by_id"`

	// Relations
	Version ComplianceLibraryVersion `json:"version" gorm:"foreignKey:ComplianceLibraryVersionID;references:ID"`
}

func (ComplianceLibraryClause) TableName() string {
	return "compliance_library_clauses"
}
