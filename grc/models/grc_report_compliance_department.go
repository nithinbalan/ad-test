package models

import "time"

type GRCReportComplianceDepartment struct {
	ID                           uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GRCReportID                  uint       `json:"grc_report_id" gorm:"not null"`
	GRCReportComplianceSummaryID uint       `json:"grc_report_compliance_summary_id" gorm:"not null"`
	OrgUnitID                    uint       `json:"organization_unit_id" gorm:"not null;column:organization_unit_id"`
	ComplianceScore              float64    `json:"compliance_score" gorm:"type:decimal(10,2) default:0.00"`
	ComplianceStatusID           uint       `json:"compliance_status_id" gorm:"not null"`
	CreatedBYID                  uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt                    time.Time  `json:"created_at" gorm:"not null"`
	UpdatedBYID                  *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt                    *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	GRCReportComplianceSummary *GRCReportComplianceSummary `json:"grc_report_compliance_summary" gorm:"foreignKey:GRCReportComplianceSummaryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	OrganizationUnit           *OrganizationUnit           `json:"organization_unit" gorm:"foreignKey:OrgUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser              *User                       `json:"created_by" gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser              *User                       `json:"updated_by" gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ComplianceStatus           *ComplianceStatus           `json:"compliance_status" gorm:"foreignKey:ComplianceStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (GRCReportComplianceDepartment) TableName() string {
	return "grc_report_compliance_departments"
}
