package models

import (
	"time"

	masterModel "jethur/grc/masters/models"
	organizationModel "jethur/grc/organization/models"
	settingsModel "jethur/grc/settings/models"
)

type ComplianceRegister struct {
	ID                      uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceCategoryID    uint      `json:"compliance_category_id" binding:"required" gorm:"not null"`
	ComplianceStatusID      uint      `json:"compliance_status_id" binding:"required" gorm:"not null"`
	ResponsibleUnitID       uint      `json:"responsible_unit_id" binding:"required" gorm:"not null"`
	ResponsibleUserID       uint      `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ReviewFrequencyID       uint      `json:"review_frequency_id" binding:"required" gorm:"not null"`
	ReviewResponsibleUserID uint      `json:"review_responsible_user_id" binding:"omitempty" gorm:"default:null"`
	SeverityLevelID         uint      `json:"severity_level_id" binding:"omitempty" gorm:"default:null"`
	ReferenceNumber         string    `json:"reference_number" binding:"omitempty,max=45" gorm:"size:45"`
	StandardCode            string    `json:"standard_code" binding:"omitempty,max=45" gorm:"size:45"`
	TitleEn                 string    `json:"title_en" binding:"required,max=255" gorm:"type:varchar(255);not null"`
	TitleAr                 string    `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	DescriptionEn           string    `json:"description_en" gorm:"type:varchar(max)"`
	DescriptionAr           string    `json:"description_ar" gorm:"type:nvarchar(max)"`
	Version                 string    `json:"version" binding:"omitempty,max=45" gorm:"size:45"`
	Deadline                string    `json:"deadline" binding:"omitempty" gorm:"type:date"`
	IssueDate               string    `json:"issue_date" binding:"omitempty" gorm:"type:date"`
	ExpiryDate              string    `json:"expiry_date" binding:"omitempty" gorm:"type:date"`
	NextReviewDate          string    `json:"next_review_date" binding:"omitempty" gorm:"type:date"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt               time.Time `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	CreatedBy               uint      `json:"created_by" gorm:"not null"`
	UpdatedBy               uint      `json:"updated_by" gorm:"default:null"`
	LifecycleStatusID       uint      `json:"lifecycle_status_id" gorm:"not null"`

	// Relationships
	ComplianceCategory    *ComplianceCategory             `gorm:"foreignkey:ComplianceCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceStatus      *ComplianceStatus               `gorm:"foreignkey:ComplianceStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit       *settingsModel.OrganizationUnit `gorm:"foreignkey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser       *organizationModel.User         `gorm:"foreignkey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewFrequency       *masterModel.ReviewFrequency    `gorm:"foreignkey:ReviewFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewResponsibleUser *organizationModel.User         `gorm:"foreignkey:ReviewResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel         *masterModel.SeverityLevel      `gorm:"foreignkey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	LifecycleStatus       *masterModel.LifecycleStatus    `gorm:"foreignkey:LifecycleStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *organizationModel.User         `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser         *organizationModel.User         `gorm:"foreignKey:UpdatedBy"`
}

type ComplianceRegisterWithAllFields struct {
	*ComplianceRegister
	ComplianceCategoryTitleEn        string `json:"compliance_category_title_en"`
	ComplianceCategoryTitleAr        string `json:"compliance_category_title_ar"`
	ComplianceStatusTitleEn          string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr          string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode        string `json:"compliance_status_color_code"`
	ResponsibleUnitTitleEn           string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr           string `json:"responsible_unit_title_ar"`
	ResponsibleUserFirstNameEn       string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr       string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn        string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr        string `json:"responsible_user_last_name_ar"`
	ReviewFrequencyTitleEn           string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr           string `json:"review_frequency_title_ar"`
	ReviewResponsibleUserFirstNameEn string `json:"review_responsible_user_first_name_en"`
	ReviewResponsibleUserFirstNameAr string `json:"review_responsible_user_first_name_ar"`
	ReviewResponsibleUserLastNameEn  string `json:"review_responsible_user_last_name_en"`
	ReviewResponsibleUserLastNameAr  string `json:"review_responsible_user_last_name_ar"`
	SeverityLevelTitleEn             string `json:"severity_level_title_en"`
	SeverityLevelTitleAr             string `json:"severity_level_title_ar"`
	SeverityLevelColorCode           string `json:"severity_level_color_code"`
	LifecycleStatusTitleEn           string `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr           string `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode         string `json:"lifecycle_status_color_code"`
	CreatedUserFirstNameEn           string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr           string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn            string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr            string `json:"created_user_last_name_ar"`
}

// TableName overrides the table name used by ComplianceRegister to `compliance_register`
func (ComplianceRegister) TableName() string {
	return "compliance_register"
}

func (ComplianceRegisterWithAllFields) TableName() string {
	// Return an empty string to indicate that ComplianceRegisterWithAllFields
	// does not have its own table name
	return ""
}
