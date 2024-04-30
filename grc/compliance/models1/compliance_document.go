package models

import (
	organizationModel "jethur/grc/models"
	"time"
)

type ComplianceDocument struct {
	ID                   uint      `json:"id" gorm:"primary_key;auto_increment"`
	ComplianceRegisterID uint      `json:"compliance_register_id"`
	Language             string    `json:"language" binding:"required"`
	URL                  string    `json:"url"`
	Name                 string    `json:"name" binding:"required"`
	Size                 int       `json:"size"`
	Token                string    `json:"token"`
	Type                 string    `json:"type"`
	Icon                 string    `json:"icon" binding:"required"`
	CreatedBy            uint      `json:"created_by" gorm:"not null"`
	UpdatedBy            uint      `json:"updated_by" gorm:"default:null"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceRegister *ComplianceRegister     `gorm:"foreignkey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *organizationModel.User `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser      *organizationModel.User `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by SLARegister to `sla_registers`
func (ComplianceDocument) TableName() string {
	return "compliance_documents"
}
