package models

import (
	"time"
)

type ComplianceReport struct {
	ID                   uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title                string    `json:"title" gorm:"type:nvarchar(500);not null" binding:"required"`
	ComplianceRegisterID uint      `json:"compliance_register_id" gorm:"not null" binding:"required"`
	StartDate            time.Time `json:"start_date" gorm:"type:date;not null" binding:"required"`
	EndDate              time.Time `json:"end_date" gorm:"type:date;not null" binding:"required"`
	CreatedBYID          uint      `json:"created_by_id" gorm:"not null" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"not null" binding:"required"`

	// Associations
	ComplianceRegister *ComplianceRegister `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Sections []*ComplianceReportSection `gorm:"foreignKey:ComplianceReportID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ComplianceReport) TableName() string {
	return "compliance_reports"
}
