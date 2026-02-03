package models

import (
	"time"
)

type ComplianceRegisterIssuingAuthorityDocument struct {
	ID                                   uint      `json:"id" gorm:"primary_key;auto_increment"`
	ComplianceRegisterIssuingAuthorityID uint      `json:"compliance_register_issuing_authority_id"`
	URL                                  string    `json:"url"`
	Name                                 string    `json:"name" binding:"required"`
	Size                                 int       `json:"size"`
	Token                                string    `json:"token"`
	OriginalName                         string    `json:"original_name"`
	Type                                 string    `json:"type"`
	CreatedByID                          uint      `json:"created_by_id" gorm:"not null"`
	UpdatedByID                          uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                            time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                            time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceRegisterIssuingAuthority *ComplianceRegisterIssuingAuthority `gorm:"foreignKey:ComplianceRegisterIssuingAuthorityID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                      *User                               `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser                      *User                               `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by ComplianceRegisterIssuingAuthorityDocument to `compliance_register_issuing_authority_documents`
func (ComplianceRegisterIssuingAuthorityDocument) TableName() string {
	return "compliance_register_issuing_authority_documents"
}
