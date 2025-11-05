package models

import "time"

type OrganizationUnitReportKeyMetrics struct {
	ID                       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationUnitReportID uint      `json:"organization_unit_report_id" gorm:"not null;index"`
	TotalUsers               int       `json:"total_users" gorm:"not null;default:0"`
	TotalRisks               int       `json:"total_risks" gorm:"not null;default:0"`
	TotalPolicies            int       `json:"total_policies" gorm:"not null;default:0"`
	TotalProcesses           int       `json:"total_processes" gorm:"not null;default:0"`
	TotalCompliances         int       `json:"total_compliances" gorm:"not null;default:0"`
	CreatedByID              uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt                time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	OrganizationUnitReport *OrganizationUnitReport `gorm:"foreignKey:OrganizationUnitReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (OrganizationUnitReportKeyMetrics) TableName() string {
	return "organization_unit_report_key_metrics"
}
