package models

import "time"

type OrganizationUnitReportRiskMitigation struct {
	ID                            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationUnitReportID      uint      `json:"organization_unit_report_id" gorm:"not null;index;"`
	RiskActionPlanStatusTitleEn   string    `json:"risk_action_plan_status_title_en" gorm:"type:nvarchar(100);not null;"`
	RiskActionPlanStatusTitleAr   string    `json:"risk_action_plan_status_title_ar" gorm:"type:nvarchar(100)"`
	RiskActionPlanStatusColorCode string    `json:"risk_action_plan_status_color_code" gorm:"type:varchar(9)"`
	Count                         int       `json:"count" gorm:"not null;default:0"`
	CreatedByID                   uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt                     time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	OrganizationUnitReport *OrganizationUnitReport `gorm:"foreignKey:OrganizationUnitReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (OrganizationUnitReportRiskMitigation) TableName() string {
	return "organization_unit_report_risk_mitigations"
}
