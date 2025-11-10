package models

import "time"

type OrganizationUnitReportComplianceSummary struct {
	ID                       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationUnitReportID uint      `json:"organization_unit_report_id" gorm:"not null;index;"`
	TotalRequirementCount    int       `json:"total_requirement_count" gorm:"not null;default:0"`
	ComplianceStatuses       *string   `json:"compliance_statuses" gorm:"type:nvarchar(max)"`
	CreatedByID              uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt                time.Time `json:"created_at" gorm:"not null"`

	OrganizationUnitReport *OrganizationUnitReport `gorm:"foreignKey:OrganizationUnitReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (OrganizationUnitReportComplianceSummary) TableName() string {
	return "organization_unit_report_compliance_summary"
}
