package models

import (
	"time"
)

type RiskRootCauseAnalysis struct {
	ID             uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID uint       `json:"risk_register_id" binding:"required" gorm:"not null"`
	Why            uint       `json:"why" binding:"required" gorm:"not null"`
	Question       string     `json:"question" binding:"required" gorm:"not null"`
	Answer         string     `json:"answer,omitempty" gorm:"default:null"`
	CreatedBy      uint       `json:"created_by" gorm:"not null"`
	UpdatedBy      *uint      `json:"updated_by,omitempty" gorm:"default:null"`
	CreatedAt      time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	RiskRegister  *RiskRegister `gorm:"foreignkey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser *User         `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by RiskRootCauseAnalysis to `risk_root_cause_analysis`
func (RiskRootCauseAnalysis) TableName() string {
	return "risk_root_cause_analysis"
}
