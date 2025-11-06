package models

import "time"

type OrganizationUnitReportUser struct {
	ID                       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationUnitReportID uint      `json:"organization_unit_report_id" gorm:"not null;index"`
	UserNameEn               string    `json:"user_name_en" gorm:"type:nvarchar(255);not null"`
	UserNameAr               string    `json:"user_name_ar" gorm:"type:nvarchar(255)"`
	DesignationTitleEn       string    `json:"designation_title_en" gorm:"type:nvarchar(255)"`
	DesignationTitleAr       string    `json:"designation_title_ar" gorm:"type:nvarchar(255)"`
	RoleTitleEn              string    `json:"role_title_en" gorm:"type:nvarchar(255);not null"`
	RoleTitleAr              string    `json:"role_title_ar" gorm:"type:nvarchar(255)"`
	CreatedByID              uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt                time.Time `json:"created_at" gorm:"not null"`

	// Relatonships
	OrganizationUnitReport *OrganizationUnitReport `gorm:"foreignKey:OrganizationUnitReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (OrganizationUnitReportUser) TableName() string {
	return "organization_unit_report_users"
}
