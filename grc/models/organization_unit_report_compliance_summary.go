package models

import "time"

type OrganizationUnitReportComplianceSummary struct {
	ID                        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationUnitReportID  uint      `json:"organization_unit_report_id" gorm:"not null;index;"`
	ComplianceStatusTitleEn   string    `json:"compliance_status_title_en" gorm:"type:nvarchar(150);not null;"`
	ComplianceStatusTitleAr   string    `json:"compliance_status_title_ar" gorm:"type:nvarchar(150)"`
	ComplianceStatusColorCode string    `json:"compliance_status_color_code" gorm:"type:varchar(9)"`
	Count                     int       `json:"count" gorm:"not null;default:0"`
	CreatedByID               uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt                 time.Time `json:"created_at" gorm:"not null"`

	OrganizationUnitReport *OrganizationUnitReport `gorm:"foreignKey:OrganizationUnitReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (OrganizationUnitReportComplianceSummary) TableName() string {
	return "organization_unit_report_compliance_summary"
}
