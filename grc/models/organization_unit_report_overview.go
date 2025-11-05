package models

import "time"

type OrganizationUnitReportOverview struct {
	ID                         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationUnitReportID   uint      `json:"organization_unit_report_id" gorm:"not null;index"`
	OrganizationUnitTitleEn    string    `json:"organization_unit_title_en" gorm:"type:nvarchar(255);not null"`
	OrganizationUnitTitleAr    string    `json:"organization_unit_title_ar" gorm:"type:nvarchar(255)"`
	OrganizationUnitHeadNameEn *string   `json:"organization_unit_head_name_en" gorm:"type:nvarchar(255);"`
	OrganizationUnitHeadNameAr *string   `json:"organization_unit_head_name_ar" gorm:"type:nvarchar(255)"`
	OrganizationUnitHeadEmail  *string   `json:"organization_unit_head_email" gorm:"type:varchar(255);"`
	TotalEmployees             int       `json:"total_employees" gorm:"not null;default:0"`
	CreatedByID                uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt                  time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	OrganizationUnitReport *OrganizationUnitReport `gorm:"foreignKey:OrganizationUnitReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (OrganizationUnitReportOverview) TableName() string {
	return "organization_unit_report_overviews"
}
