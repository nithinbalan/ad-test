package models

import (
	"time"
)

type ViolationReportRootCauseAnalysis struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ViolationReportID uint       `json:"violation_report_id" binding:"required" gorm:"not null"`
	Why               uint       `json:"why" binding:"required" gorm:"not null"`
	Question          string     `json:"question" binding:"required" gorm:"not null"`
	Answer            string     `json:"answer,omitempty" gorm:"default:null"`
	CreatedBYID       uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID       *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	ViolationReport *ViolationReport `gorm:"foreignKey:ViolationReportID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User            `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser   *User            `gorm:"foreignKey:UpdatedBYID"`
}

func (ViolationReportRootCauseAnalysis) TableName() string {
	return "violation_report_root_cause_analyses"
}
