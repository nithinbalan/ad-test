package models

import (
	"time"
)

type KRIValueUpdate struct {
	ID                  uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KRIRegisterID       uint       `json:"kri_register_id" binding:"required" gorm:"not null"`
	ReviewedDate        time.Time  `json:"reviewed_date" binding:"required" gorm:"type:date;not null"`
	UnitID              uint       `json:"unit_id" binding:"required" gorm:"not null"`
	Value               *float64   `json:"value" gorm:"type:decimal(18,2);not null"`
	LowThreshold        *int       `json:"low_threshold" gorm:"type:int;null"`
	MediumThreshold     *int       `json:"medium_threshold" gorm:"type:int;null"`
	HighThreshold       *int       `json:"high_threshold" gorm:"type:int;null"`
	Remarks             *string    `json:"remarks" gorm:"type:nvarchar(max)"`
	IsLatest            bool       `json:"is_latest" gorm:"type:bit;not null"`
	IsPrevious          bool       `json:"is_previous" gorm:"type:bit;not null;default:false"`
	KRIFrequencyValueID *uint      `json:"kri_frequency_value_id" gorm:"type:int;default:null"`
	CreatedAt           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID         uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID         *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	KRIRegister       *KRIRegister             `gorm:"foreignKey:KRIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Unit              *Unit                    `gorm:"foreignKey:UnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Documents         []KRIValueUpdateDocument `gorm:"foreignKey:KRIValueUpdateID"`
	CreatedByUser     *User                    `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser     *User                    `gorm:"foreignKey:UpdatedByID"`
	KRIFrequencyValue *KRIFrequencyValue       `gorm:"foreignKey:KRIFrequencyValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by KRIValueUpdate to `kri_value_updates`
func (KRIValueUpdate) TableName() string {
	return "kri_value_updates"
}
