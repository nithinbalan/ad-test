package models

import "time"

type OrganizationUnitReportUserRole struct {
	ID                       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationUnitReportID uint      `json:"organization_unit_report_id" gorm:"not null;index"`
	RoleTitleEn              string    `json:"role_title_en" gorm:"type:nvarchar(255);not null"`
	RoleTitleAr              string    `json:"role_title_ar" gorm:"type:nvarchar(255)"`
	Count                    int       `json:"count" gorm:"not null;default:0"`
	CreatedByID              uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt                time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	OrganizationUnitReport *OrganizationUnitReport `gorm:"foreignKey:OrganizationUnitReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (OrganizationUnitReportUserRole) TableName() string {
	return "organization_unit_report_user_roles"
}
