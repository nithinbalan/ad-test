package models

import (
	"time"
)

type ViolationReportControl struct {
	ID                uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ViolationReportID uint      `json:"violation_report_id" binding:"required"`
	ControlID         uint      `json:"control_id" binding:"required"`
	CreatedByID       uint      `json:"created_by_id" binding:"required"`
	CreatedAt         time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ViolationReport *ViolationReport `gorm:"foreignKey:ViolationReportID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control         *Control         `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User            `gorm:"foreignKey:CreatedByID"`
}

func (ViolationReportControl) TableName() string {
	return "violation_report_controls"
}
