package models

import "time"

type OrganizationUnitReportPolicySummary struct {
	ID                       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationUnitReportID uint      `json:"organization_unit_report_id" gorm:"not null;index"`
	TotalPolicyCount         int       `json:"total_policy_count" gorm:"not null;default:0"`
	ActivePolicyCount        int       `json:"active_policy_count" gorm:"not null;default:0"`
	UnderReviewPolicyCount   int       `json:"under_review_policy_count" gorm:"not null;default:0"`
	PolicyHealthPercent      float32   `json:"policy_health_percent" gorm:"type:numeric(5,2);not null;default:0"`
	PolicyAwarenessPercent   float32   `json:"policy_awareness_percent" gorm:"type:numeric(5,2);not null;default:0"`
	CreatedByID              uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt                time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	OrganizationUnitReport *OrganizationUnitReport `gorm:"foreignKey:OrganizationUnitReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (OrganizationUnitReportPolicySummary) TableName() string {
	return "organization_unit_report_policy_summary"
}
