package models

import (
	"time"

	models "jethur/grc/models"
)

type KRIValueUpdate struct {
	ID              uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KRIRegisterID   uint       `json:"kri_register_id" binding:"required" gorm:"not null"`
	ReviewedDate    time.Time  `json:"reviewed_date" binding:"required" gorm:"type:date;not null"`
	Value           string     `json:"value" gorm:"type:varchar(455);not null"`
	LowThreshold    string     `json:"low_threshold" gorm:"type:varchar(455);null"`
	MediumThreshold string     `json:"medium_threshold" gorm:"type:varchar(455);null"`
	HighThreshold   string     `json:"high_threshold" gorm:"type:varchar(455);null"`
	Remarks         *string    `json:"remarks" gorm:"type:nvarchar(max)"`
	CreatedAt       time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBy       uint       `json:"created_by" gorm:"not null"`
	UpdatedBy       *uint      `json:"updated_by" gorm:"default:null"`

	// Relationships
	KRIRegister             *KRIRegister             `gorm:"foreignkey:KRIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	KRIValueUpdateDocuments []KRIValueUpdateDocument `gorm:"foreignKey:KRIValueUpdateID"`
	CreatedByUser           *models.User             `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser           *models.User             `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by KRIValueUpdate to `kri_value_updates`
func (KRIValueUpdate) TableName() string {
	return "kri_value_updates"
}
