package models

import "time"

type OrganizationUnitReportRiskSummary struct {
	ID                       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationUnitReportID uint      `json:"organization_unit_report_id" gorm:"not null;index"`
	TotalRiskCount           int       `json:"total_risk_count" gorm:"not null;default:0"`
	OpenRiskCount            int       `json:"open_risk_count" gorm:"not null;default:0"`
	ClosedRiskCount          int       `json:"closed_risk_count" gorm:"not null;default:0"`
	CreatedByID              uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt                time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	OrganizationUnitReport *OrganizationUnitReport `gorm:"foreignKey:OrganizationUnitReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (OrganizationUnitReportRiskSummary) TableName() string {
	return "organization_unit_report_risk_summary"
}
