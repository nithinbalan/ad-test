package models

import (
	"time"
)

type ViolationReportAssociatedDocument struct {
	ID                uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ViolationReportID uint      `json:"violation_report_id" binding:"required"`
	DocumentID        uint      `json:"document_id" binding:"required"`
	CreatedByID       uint      `json:"created_by_id" binding:"required"`
	CreatedAt         time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ViolationReport *ViolationReport `gorm:"foreignKey:ViolationReportID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document        *Document        `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User            `gorm:"foreignKey:CreatedByID"`
}

func (ViolationReportAssociatedDocument) TableName() string {
	return "violation_report_associated_documents"
}
