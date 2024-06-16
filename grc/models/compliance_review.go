package models

import (
	"time"
)

type ComplianceReview struct {
	ID                   uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRegisterID uint       `json:"compliance_register_id" gorm:"not null"`
	Comment              string     `json:"comment" binding:"required" gorm:"type:nvarchar(max);not null"`
	ReviewedAt           time.Time  `json:"reviewed_at" gorm:"type:datetime;not null"`
	ReviewedBy           uint       `json:"reviewed_by" gorm:"tnot null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	UpdatedByID          *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	ComplianceRegister *ComplianceRegister `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewedByUser     *User               `gorm:"foreignKey:ReviewedBy"`
	UpdatedByUser      *User               `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by ComplianceReview to `compliance_reviews`
func (ComplianceReview) TableName() string {
	return "compliance_reviews"
}
