package models

import (
	"time"
)

type KRIRegister struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber   string     `json:"reference_number"`
	KRICategoryID     uint       `json:"kri_category_id" binding:"required" gorm:"not null"`
	ResponsibleUnitID uint       `json:"responsible_unit_id" binding:"required" gorm:"not null"`
	ResponsibleUserID uint       `json:"responsible_user_id" gorm:"not null"`
	ReviewFrequencyID uint       `json:"review_frequency_id" binding:"required" gorm:"not null"`
	RiskLevelID       *uint      `json:"risk_level_id"`
	UnitID            uint       `json:"unit_id" binding:"required" gorm:"not null"`
	Title             string     `json:"title" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	Description       string     `json:"description" gorm:"type:nvarchar(max)"`
	CalculationType   string     `json:"calculation_type" binding:"required,oneof=automatic, manual" gorm:"type:varchar(10);not null"`
	LowThreshold      *int       `json:"low_threshold" gorm:"type:int;null"`
	MediumThreshold   *int       `json:"medium_threshold" gorm:"type:int;null"`
	HighThreshold     *int       `json:"high_threshold" gorm:"type:int;null"`
	CurrentValue      *string    `json:"current_value" gorm:"type:varchar(455);null"`
	LastReviewedDate  *time.Time `json:"last_reviewed_date" gorm:"type:date"`
	NextReviewDate    time.Time  `json:"next_review_date" binding:"required" gorm:"type:date;not null"`
	LastReviewNote    *string    `json:"last_review_note" gorm:"type:nvarchar(max)"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID       uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID       *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	KRICategory     *KRICategory      `gorm:"foreignKey:KRICategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit *OrganizationUnit `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser *User             `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewFrequency *ReviewFrequency  `gorm:"foreignKey:ReviewFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskLevel       *RiskLevel        `gorm:"foreignKey:RiskLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User             `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser   *User             `gorm:"foreignKey:UpdatedBYID"`
	ValueUpdates    []*KRIValueUpdate `gorm:"foreignKey:KRIRegisterID"`
}

// TableName overrides the table name used by KRIRegister to `kri_register`
func (KRIRegister) TableName() string {
	return "kri_register"
}
