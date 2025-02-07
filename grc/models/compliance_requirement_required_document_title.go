package models

import (
	"time"
)

type ComplianceRequirementRequiredDocumentTitle struct {
	ID                      uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRequirementID uint      `json:"compliance_requirement_id" binding:"required"`
	RequiredDocumentTitleID uint      `json:"required_document_title_id" binding:"required"`
	CreatedBYID             uint      `json:"created_by_id" binding:"required"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceRequirement *ComplianceRequirement `gorm:"foreignKey:ComplianceRequirementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RequiredDocumentTitle *RequiredDocumentTitle `gorm:"foreignKey:RequiredDocumentTitleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedBYID"`
}

func (ComplianceRequirementRequiredDocumentTitle) TableName() string {
	return "compliance_requirement_required_document_titles"
}
