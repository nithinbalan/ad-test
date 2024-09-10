package models

import (
	"time"
)

type ComplianceRequirementAssociatedDocument struct {
	ID                      uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRegisterID    uint      `json:"compliance_register_id" binding:"required"`
	ComplianceRequirementID uint      `json:"compliance_requirement_id" binding:"required"`
	DocumentID              uint      `json:"document_id" binding:"required"`
	CreatedByID             uint      `json:"created_by_id" binding:"required"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceRegister    *ComplianceRegister    `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRequirement *ComplianceRequirement `gorm:"foreignKey:ComplianceRequirementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document              *Document              `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedByID"`
}

func (ComplianceRequirementAssociatedDocument) TableName() string {
	return "compliance_requirement_associated_documents"
}
