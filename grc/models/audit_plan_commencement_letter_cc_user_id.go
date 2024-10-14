package models

import (
	"time"
)

type AuditPlanCommencementLetterCCUserID struct {
	ID                            uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditPlanCommencementLetterID uint       `json:"audit_plan_commencement_letter_id" binding:"required" gorm:"not null"`
	UserID                        uint       `json:"user_id" binding:"required" gorm:"not null"`
	CreatedByID                   uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt                     time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedByID                   *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt                     *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	AuditPlanCommencementLetter *AuditPlanCommencementLetter `gorm:"foreignKey:AuditPlanCommencementLetterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User                        *User                        `gorm:"foreignKey:UserID"`
	CreatedByUser               *User                        `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser               *User                        `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by AuditPlanCommencementLetterCCUserID to `audit_plan_commencement_letter_cc_user_ids`
func (AuditPlanCommencementLetterCCUserID) TableName() string {
	return "audit_plan_commencement_letter_cc_user_ids"
}
