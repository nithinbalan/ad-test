package models

import "time"

type OrganizationUnitReportComplianceRequirement struct {
	ID                                   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationUnitReportID             uint   `json:"organization_unit_report_id" gorm:"not null;index"`
	TitleEn                              string `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr                              string `json:"title_ar" gorm:"type:nvarchar(255)"`
	ComplianceRequirementCategoryTitleEn string `json:"compliance_requirement_category_title_en" gorm:"type:nvarchar(255)"`
	ComplianceRequirementCategoryTitleAr string `json:"compliance_requirement_category_title_ar" gorm:"type:nvarchar(255)"`
	ComplianceStatusTitleEn              string `json:"compliance_status_title_en" gorm:"type:nvarchar(150)"`
	ComplianceStatusTitleAr              string `json:"compliance_status_title_ar" gorm:"type:nvarchar(150)"`
	ComplianceStatusColorCode            string `json:"compliance_status_color_code" gorm:"type:varchar(9)"`
	SeverityLevelTitleEn                 string `json:"severity_level_title_en" gorm:"type:nvarchar(100)"`
	SeverityLevelTitleAr                 string `json:"severity_level_title_ar" gorm:"type:nvarchar(100)"`
	SeverityLevelColorCode               string `json:"severity_level_color_code" gorm:"type:varchar(9)"`

	CreatedByID uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	OrganizationUnitReport *OrganizationUnitReport `gorm:"foreignKey:OrganizationUnitReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (OrganizationUnitReportComplianceRequirement) TableName() string {
	return "organization_unit_report_compliance_requirements"
}
