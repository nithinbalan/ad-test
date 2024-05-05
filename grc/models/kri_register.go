package models

import (
	"time"
)

type KRIRegister struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber   string     `json:"reference_number"`
	KRICategoryID     uint       `json:"kri_category_id" binding:"required" gorm:"not null"`
	OwnerUnitID       uint       `json:"owner_unit_id" binding:"required" gorm:"not null"`
	OwnerID           uint       `json:"owner_id" gorm:"not null"`
	ReviewFrequencyID uint       `json:"review_frequency_id" binding:"required" gorm:"not null"`
	RiskLevelID       uint       `json:"risk_level_id" binding:"required" gorm:"not null"`
	Title             string     `json:"title" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	Description       string     `json:"description" gorm:"type:nvarchar(max)"`
	CalculationType   string     `json:"calculation_type" binding:"required,oneof=automatic, manual" gorm:"type:varchar(10);not null"`
	LowThreshold      *string    `json:"low_threshold" gorm:"type:varchar(455);null"`
	MediumThreshold   *string    `json:"medium_threshold" gorm:"type:varchar(455);null"`
	HighThreshold     *string    `json:"high_threshold" gorm:"type:varchar(455);null"`
	CurrentValue      *string    `json:"current_value" gorm:"type:varchar(455);null"`
	LastReviewedDate  time.Time  `json:"last_reviewed_date" gorm:"type:date"`
	NextReviewDate    time.Time  `json:"next_review_date" binding:"required" gorm:"type:date;not null"`
	LastReviewNote    *string    `json:"last_review_note" gorm:"type:nvarchar(max)"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBy         uint       `json:"created_by" gorm:"not null"`
	UpdatedBy         *uint      `json:"updated_by" gorm:"default:null"`

	// Relationships
	KRICategory     *KRICategory      `gorm:"foreignkey:KRICategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OwnerUnit       *OrganizationUnit `gorm:"foreignkey:OwnerUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Owner           *User             `gorm:"foreignkey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewFrequency *ReviewFrequency  `gorm:"foreignkey:ReviewFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskLevel       *RiskLevel        `gorm:"foreignkey:RiskLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User             `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser   *User             `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by KRIRegister to `kri_register`
func (KRIRegister) TableName() string {
	return "kri_register"
}
