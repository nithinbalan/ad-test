package models

import (
	"time"
)

type KRIFrequencyValue struct {
	ID                    uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KRIRegisterID         uint      `json:"kri_register_id" binding:"required" gorm:"not null"`
	Year                  int       `json:"year" binding:"required" gorm:"type:int;not null"`
	MonthEn               string    `json:"month_en" binding:"required" gorm:"type:varchar(1000);not null"`
	MonthAr               string    `json:"month_ar" binding:"required" gorm:"type:varchar(1000);not null"`
	Type                  string    `json:"type" binding:"required" gorm:"type:varchar(1000);not null"`
	TargetValue           float64   `json:"target_value" binding:"required" gorm:"type:decimal(18,2);not null"`
	ActualValue           *float64  `json:"actual_value" gorm:"type:decimal(18,2);null"`
	AchievementValue      *float64  `json:"achievement_value" gorm:"type:decimal(18,2);null"`
	RiskToleranceStatusID *uint     `json:"risk_tolerance_status_id" gorm:"default:null"`
	IsAchieved            bool      `json:"is_achieved" gorm:"default:false"`
	IsMissed              bool      `json:"is_missed" gorm:"default:false"`
	CreatedBYID           *uint     `json:"created_by_id" gorm:"default:null"`
	CreatedAt             time.Time `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	// Relationships
	KRIRegister         *KRIRegister            `gorm:"foreignKey:KRIRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBy           *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskToleranceStatus *KRIRiskToleranceStatus `gorm:"foreignKey:RiskToleranceStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by KRIFrequencyValue to `kri_frequency_values`
func (KRIFrequencyValue) TableName() string {
	return "kri_frequency_values"
}
