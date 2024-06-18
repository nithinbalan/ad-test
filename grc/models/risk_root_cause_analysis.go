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
	CreatedByID    uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID    *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	CreatedAt      time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	RiskRegister  *RiskRegister `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser *User         `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by RiskRootCauseAnalysis to `risk_root_cause_analysis`
func (RiskRootCauseAnalysis) TableName() string {
	return "risk_root_cause_analysis"
}
