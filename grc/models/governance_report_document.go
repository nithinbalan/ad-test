package models

import "time"

type GovernanceReportDocument struct {
	ID                     uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GovernanceReportID     uint      `json:"governance_report_id" gorm:"not null"`
	ReferenceNumber        string    `json:"reference_number" gorm:"not null"`
	DocumentTitleEn        string    `json:"document_title_en" gorm:"type:nvarchar(255);not null"`
	DocumentTitleAr        string    `json:"document_title_ar" gorm:"type:nvarchar(255);not null"`
	ResponsibleUserNameEn  string    `json:"responsible_user_name_en" gorm:"type:nvarchar(255);not null"`
	ResponsibleUserNameAr  string    `json:"responsible_user_name_ar" gorm:"type:nvarchar(255);not null"`
	ResponsibleUnitTitleEn string    `json:"responsible_unit_title_en" gorm:"type:nvarchar(255);not null"`
	ResponsibleUnitTitleAr string    `json:"responsible_unit_title_ar" gorm:"type:nvarchar(255);not null"`
	TotalEmployees         int       `json:"total_employees" gorm:"default:0;not null"`
	AcknowledgedCount      int       `json:"acknowledged_count" gorm:"default:0;not null"`
	PendingCount           int       `json:"pending_count" gorm:"default:0;not null"`
	CreatedByID            uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt              time.Time `json:"created_at" gorm:"not null"`

	GovernanceReport *GovernanceReport `gorm:"foreignKey:GovernanceReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBy        *User             `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (GovernanceReportDocument) TableName() string {
	return "governance_report_documents"
}
