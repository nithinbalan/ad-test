package models

import "time"

type PolicyReportDocumentAccess struct {
	ID                      uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PolicyReportID          uint      `json:"policy_report_id" gorm:"not null;index"`
	UserNameEn              string    `json:"user_name_en" gorm:"type:nvarchar(255)"`
	UserNameAr              string    `json:"user_name_ar" gorm:"type:nvarchar(255)"`
	DesignationTitleEn      string    `json:"designation_title_en" gorm:"type:nvarchar(255)"`
	DesignationTitleAr      string    `json:"designation_title_ar" gorm:"type:nvarchar(255)"`
	OrganizationUnitTitleEn string    `json:"organization_unit_title_en" gorm:"type:nvarchar(255)"`
	OrganizationUnitTitleAr string    `json:"organization_unit_title_ar" gorm:"type:nvarchar(255)"`
	CreatedByID             uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt               time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	PolicyReport  *PolicyReport `gorm:"foreignKey:PolicyReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PolicyReportDocumentAccess) TableName() string {
	return "policy_report_document_accesses"
}
