package models

import "time"

type OrganizationUnitReportPolicy struct {
	ID                       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationUnitReportID uint      `json:"organization_unit_report_id" gorm:"not null;index"`
	DocumentTitleEn          string    `json:"document_title_en" gorm:"type:nvarchar(255);not null"`
	DocumentTitleAr          string    `json:"document_title_ar" gorm:"type:nvarchar(255)"`
	DocumentCategoryTitleEn  string    `json:"document_category_title_en" gorm:"type:nvarchar(255)"`
	DocumentCategoryTitleAr  string    `json:"document_category_title_ar" gorm:"type:nvarchar(255)"`
	DocumentTypeTitleEn      string    `json:"document_type_title_en" gorm:"type:nvarchar(100)"`
	DocumentTypeTitleAr      string    `json:"document_type_title_ar" gorm:"type:nvarchar(100)"`
	ResponsibleUserNameEn    string    `json:"responsible_user_name_en" gorm:"type:nvarchar(255)"`
	ResponsibleUserNameAr    string    `json:"responsible_user_name_ar" gorm:"type:nvarchar(255)"`
	Acknowledged             bool      `json:"acknowledged" gorm:"not null;default:false"`
	CreatedByID              uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt                time.Time `json:"created_at" gorm:"not null"`

	OrganizationUnitReport *OrganizationUnitReport `gorm:"foreignKey:OrganizationUnitReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (OrganizationUnitReportPolicy) TableName() string {
	return "organization_unit_report_policies"
}
