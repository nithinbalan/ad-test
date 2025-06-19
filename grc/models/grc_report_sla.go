package models

import "time"

type GRCReportSLA struct {
	ID                       uint       `json:"id" gorm:"primary_key;auto_increment"`
	GRCReportID              uint       `json:"grc_report_id" gorm:"not null"`
	ParentOrganizationUnitID uint       `json:"parent_organization_unit_id" gorm:"not null"`
	OrganizationUnitID       uint       `json:"organization_unit_id" gorm:"not null"`
	ReferenceNumber          string     `json:"reference_number" gorm:"not null"`
	TitleEn                  string     `json:"title_en" gorm:"not null"`
	TitleAr                  *string    `json:"title_ar" gorm:"default:null"`
	ComplianceStatusID       *uint      `json:"compliance_status_id" gorm:"default:null"`
	LifecycleStatusID        uint       `json:"lifecycle_status_id" gorm:"not null"`
	IssueDate                *time.Time `json:"issue_date" gorm:"type:date;default:null"`
	ExpiryDate               *time.Time `json:"expiry_date" gorm:"type:date;default:null"`
	CreatedBYID              uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"not null;type:datetime"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"default:null;type:datetime"`

	// Relationships
	GRCReport              *GRCReport        `json:"grc_report" gorm:"foreignKey:GRCReportID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentOrganizationUnit *OrganizationUnit `json:"parent_organization_unit" gorm:"foreignKey:ParentOrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OrganizationUnit       *OrganizationUnit `json:"organization_unit" gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ComplianceStatus       *ComplianceStatus `json:"compliance_status" gorm:"foreignKey:ComplianceStatusID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LifecycleStatus        *LifecycleStatus  `json:"lifecycle_status" gorm:"foreignKey:LifecycleStatusID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedByUser          *User             `json:"created_by_user" gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UpdatedByUser          *User             `json:"updated_by_user" gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (GRCReportSLA) TableName() string {
	return "grc_report_sla"
}
