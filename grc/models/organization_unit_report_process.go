package models

import "time"

type OrganizationUnitReportProcess struct {
	ID                       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationUnitReportID uint      `json:"organization_unit_report_id" gorm:"not null;index"`
	TitleEn                  string    `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr                  string    `json:"title_ar" gorm:"type:nvarchar(255)"`
	ProcessGroupTitleEn      string    `json:"process_group_title_en" gorm:"type:nvarchar(255)"`
	ProcessGroupTitleAr      string    `json:"process_group_title_ar" gorm:"type:nvarchar(255)"`
	ProcessOwnerNameEn       string    `json:"process_owner_name_en" gorm:"type:nvarchar(100)"`
	ProcessOwnerNameAR       string    `json:"process_owner_name_ar" gorm:"type:nvarchar(100)"`
	ResponsibleUnitTitleEn   string    `json:"responsible_unit_title_en" gorm:"type:nvarchar(255)"`
	ResponsibleUnitTitleAr   string    `json:"responsible_unit_title_ar" gorm:"type:nvarchar(255)"`
	RiskRatingTitleEn        string    `json:"risk_rating_title_en" gorm:"type:nvarchar(100)"`
	RiskRatingTitleAr        string    `json:"risk_rating_title_ar" gorm:"type:nvarchar(100)"`
	LifecycleStatusTitleEn   string    `json:"lifecycle_status_title_en" gorm:"type:nvarchar(100)"`
	LifecycleStatusTitleAr   string    `json:"lifecycle_status_title_ar" gorm:"type:nvarchar(100)"`
	CreatedByID              uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt                time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	OrganizationUnitReport *OrganizationUnitReport `gorm:"foreignKey:OrganizationUnitReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (OrganizationUnitReportProcess) TableName() string {
	return "organization_unit_report_processes"
}
