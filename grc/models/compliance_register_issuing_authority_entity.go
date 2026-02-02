package models

import "time"

type ComplianceRegisterIssuingAuthorityEntity struct {
	ID                                   uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceRegisterIssuingAuthorityID uint      `json:"compliance_register_issuing_authority_id" gorm:"not null;index"`
	CircularNo                           string    `json:"circular_no" gorm:"type:nvarchar(500);not null"`
	DateOfReceiving                      time.Time `json:"date_of_receiving" gorm:"type:datetime;not null"`
	HijriDateOfReceiving                 time.Time `json:"hijri_date_of_receiving" gorm:"type:datetime;not null"`
	CreatedAt                            time.Time `json:"created_at" gorm:"type:datetime;not null"`
	CreatedBYID                          uint      `json:"created_by_id" gorm:"not null;index"`

	// Relationships
	ComplianceRegisterIssuingAuthority *ComplianceRegisterIssuingAuthority `gorm:"foreignKey:ComplianceRegisterIssuingAuthorityID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                      *User                               `gorm:"foreignKey:CreatedBYID"`
}

// Override table name
func (ComplianceRegisterIssuingAuthorityEntity) TableName() string {
	return "compliance_register_issuing_authority_entities"
}
