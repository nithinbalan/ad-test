package models

import "time"

type ViolationReportComment struct {
	ID                      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	ViolationReportID       uint   `json:"violation_report_id" gorm:"not null"`
	ViolationReportStatusID *uint  `json:"violation_report_status_id" gorm:"not null"`
	Comment                 string `json:"comment" gorm:"type:nvarchar(MAX);not null"`

	// Audit Fields
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	CreatedByID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	UpdatedByID *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	ViolationReport       *ViolationReport       `gorm:"foreignKey:ViolationReportID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ViolationReportStatus *ViolationReportStatus `gorm:"foreignKey:ViolationReportStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser         *User                  `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the default table name for GORM
func (ViolationReportComment) TableName() string {
	return "violation_report_comments"
}
