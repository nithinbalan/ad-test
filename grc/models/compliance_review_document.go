package models

import (
	"time"
)

type ComplianceReviewDocument struct {
	ID                 uint      `json:"id" gorm:"primary_key;auto_increment"`
	ComplianceReviewID uint      `json:"compliance_review_id"`
	URL                string    `json:"url"`
	Name               string    `json:"name" binding:"required"`
	Size               int       `json:"size"`
	Token              string    `json:"token"`
	OriginalName       string    `json:"original_name"`
	Type               string    `json:"type"`
	CreatedBy          uint      `json:"created_by_id" gorm:"not null"`
	UpdatedBy          uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceReview *ComplianceReview `gorm:"foreignKey:ComplianceReviewID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser    *User             `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by ComplianceReviewDocument to `compliance_review_documents`
func (ComplianceReviewDocument) TableName() string {
	return "compliance_review_documents"
}
