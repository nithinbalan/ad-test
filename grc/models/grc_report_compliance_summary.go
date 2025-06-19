package models

import "time"

type GRCReportComplianceSummary struct {
	ID                     uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GRCReportID            uint       `json:"grc_report_id" gorm:"not null"`
	ReviewUserID           *uint      `json:"review_user_id"`
	ApproveUserID          *uint      `json:"approve_user_id"`
	OrgUnitID              uint       `json:"organization_unit_id" gorm:"not null;column:organization_unit_id"`
	OverallComplianceScore *float64   `json:"overall_compliance_score" gorm:"type:decimal(10,2) default:0.00"`
	CreatedBYID            uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt              time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID            *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt              *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	GRCReport                      *GRCReport                      `gorm:"foreignKey:GRCReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser                  *User                           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser                  *User                           `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ReviewUser                     *User                           `gorm:"foreignKey:ReviewUserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ApproveUser                    *User                           `gorm:"foreignKey:ApproveUserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	OrganizationUnit               *OrganizationUnit               `gorm:"foreignKey:OrgUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	GRCReportComplianceDepartments []GRCReportComplianceDepartment `gorm:"foreignKey:GRCReportComplianceSummaryID;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (GRCReportComplianceSummary) TableName() string {
	return "grc_report_compliance_summary"
}
