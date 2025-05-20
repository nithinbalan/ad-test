package models

import "time"

type GRCReportGovernanceDocument struct {
	ID                      uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GRCReportID             uint       `json:"grc_report_id" gorm:"not null"`
	ReferenceNumber         string     `json:"reference_number" gorm:"type:nvarchar(50);not null"`
	DocumentTitleEn         string     `json:"document_title_en" gorm:"type:nvarchar(255);not null"`
	DocumentTitleAr         string     `json:"document_title_ar" gorm:"type:nvarchar(255);not null"`
	ResponsibleUserNameEn   string     `json:"responsible_user_name_en" gorm:"type:nvarchar(255);not null"`
	ResponsibleUserNameAr   string     `json:"responsible_user_name_ar" gorm:"type:nvarchar(255);not null"`
	ResponsibleUnitTitleEn  string     `json:"responsible_unit_title_en" gorm:"type:nvarchar(255);not null"`
	ResponsibleUnitTitleAr  string     `json:"responsible_unit_title_ar" gorm:"type:nvarchar(255);not null"`
	TotalEmployees          int        `json:"total_employees" gorm:"not null;default:0"`
	AcknowledgedCount       int        `json:"acknowledged_count" gorm:"not null;default:0"`
	PendingCount            int        `json:"pending_count" gorm:"not null;default:0"`
	DocumentStatusTitleEn   string     `json:"document_status_title_en" gorm:"type:nvarchar(255);not null"`
	DocumentStatusTitleAr   string     `json:"document_status_title_ar" gorm:"type:nvarchar(255);not null"`
	DocumentStatusColorCode string     `json:"document_status_color_code" gorm:"type:nvarchar(50);not null"`
	AcknowledgeDueDate      *time.Time `json:"acknowledge_due_date" gorm:"type:date;default:null"`
	CreatedBYID             uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt               time.Time  `json:"created_at" gorm:"type:datetime;not null"`

	GRCReport     *GRCReport `gorm:"foreignKey:GRCReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User      `gorm:"foreignKey:CreatedBYID"`
}

func (GRCReportGovernanceDocument) TableName() string {
	return "grc_report_governance_documents"
}
