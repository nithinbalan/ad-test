package models

import "time"

type ComplianceRegisterIssuingAuthority struct {
	ID                   uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceRegisterID uint      `json:"compliance_register_id" gorm:"not null;index"`
	IssuingAuthorityID   uint      `json:"issuing_authority_id" gorm:"not null;index"`
	CreatedBYID          uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceRegister *ComplianceRegister                         `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IssuingAuthority   *ComplianceIssuingAuthority                 `gorm:"foreignKey:IssuingAuthorityID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User                                       `gorm:"foreignKey:CreatedBYID"`
	Entity             *ComplianceRegisterIssuingAuthorityEntity   `gorm:"foreignKey:ComplianceRegisterIssuingAuthorityID"`
	Document           *ComplianceRegisterIssuingAuthorityDocument `gorm:"foreignKey:ComplianceRegisterIssuingAuthorityID"`
}

// TableName overrides the table name
func (ComplianceRegisterIssuingAuthority) TableName() string {
	return "compliance_register_issuing_authorities"
}
