package models

import "time"

type PolicyReportAcknowledgementOrganizationUnit struct {
	ID                      uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PolicyReportID          uint      `json:"policy_report_id" gorm:"not null;index"`
	OrganizationUnitTitleEn string    `json:"organization_unit_title_en" gorm:"type:nvarchar(255)"`
	OrganizationUnitTitleAr string    `json:"organization_unit_title_ar" gorm:"type:nvarchar(255)"`
	TotalUsers              int       `json:"total_users" gorm:"not null;default:0"`
	AcknowledgedUsers       int       `json:"acknowledged_users" gorm:"not null;default:0"`
	NotAcknowledgedUsers    int       `json:"not_acknowledged_users" gorm:"not null;default:0"`
	AcknowledgedPercent     float32   `json:"acknowledged_percent" gorm:"type:numeric(5,2);not null;default:0"`
	CreatedByID             uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt               time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	PolicyReport  *PolicyReport `gorm:"foreignKey:PolicyReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PolicyReportAcknowledgementOrganizationUnit) TableName() string {
	return "policy_report_acknowledgement_organization_units"
}
