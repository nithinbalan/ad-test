package models

import (
	"time"

	organizationModels "jethur/grc/organization/models"
)

type ComplianceReview struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRegisterID uint      `json:"compliance_register_id" gorm:"not null"`
	Comment              string    `json:"comment" binding:"required" gorm:"type:nvarchar(max);not null"`
	ReviewedAt           time.Time `json:"reviewed_at" gorm:"type:datetime;not null"`
	ReviewedBy           uint      `json:"reviewed_by" gorm:"tnot null"`

	// Relationships
	ComplianceRegister *ComplianceRegister      `gorm:"foreignkey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewedByUser     *organizationModels.User `gorm:"foreignKey:ReviewedBy"`
}

// TableName overrides the table name used by ComplianceReview to `compliance_reviews`
func (ComplianceReview) TableName() string {
	return "compliance_reviews"
}
