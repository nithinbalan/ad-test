package models

import (
	"time"
)

type KPIRegister struct {
	ID                uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Title             string  `json:"title" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	Description       *string `json:"description" gorm:"type:nvarchar(max)"`
	KPICategoryID     uint    `json:"kri_category_id"`
	ResponsibleUserID *uint   `json:"responsible_user_id"`
	KPIStatusID       uint    `json:"kpi_status_id"`

	ReviewFrequencyID *uint    `json:"review_frequency_id"`
	CalculationTypeID *uint    `json:"calculation_type_id"`
	Formula           *string  `json:"formula"`
	UnitID            *uint    `json:"unit_id"`
	Source            *string  `json:"source"`
	ActualValue       *float64 `json:"actual_value"`

	TargetValue     *float64 `json:"target_value"`
	LowThreshold    *int     `json:"low_threshold"`
	MediumThreshold *int     `json:"medium_threshold"`
	Weightage       *float64 `json:"weightage"`

	StartDate *time.Time `json:"start_date" gorm:"type:date;default:null"`
	EndDate   *time.Time `json:"end_date" gorm:"type:date;default:null"`

	Remarks *string `json:"remarks" gorm:"type:nvarchar(max)"`

	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	KPICategory     *KPICategory     `gorm:"foreignKey:KPICategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser *User            `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	KPIStatus       *KPIStatus       `gorm:"foreignKey:KPIStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewFrequency *ReviewFrequency `gorm:"foreignKey:ReviewFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CalculationType *CalculationType `gorm:"foreignKey:CalculationTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User            `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser   *User            `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by KPIRegister to `kpi_register`
func (KPIRegister) TableName() string {
	return "kpi_register"
}
