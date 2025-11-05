package models

import "time"

type PolicyReportOverview struct {
	ID                          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	PolicyReportID              uint       `json:"policy_report_id" gorm:"not null;index"`
	ReferenceNumber             string     `json:"reference_number" gorm:"type:varchar(50);index"`
	PolicyTitleEn               string     `json:"policy_title_en" gorm:"type:nvarchar(255);not null"`
	PolicyTitleAr               string     `json:"policy_title_ar" gorm:"type:nvarchar(255)"`
	DescriptionEn               string     `json:"description_en" gorm:"type:nvarchar(4000)"`
	DescriptionAr               string     `json:"description_ar" gorm:"type:nvarchar(4000)"`
	ResponsibleUnitTitleEn      string     `json:"responsible_unit_title_en" gorm:"type:nvarchar(255)"`
	ResponsibleUnitTitleAr      string     `json:"responsible_unit_title_ar" gorm:"type:nvarchar(255)"`
	ResponsibleUserFirstNameEn  string     `json:"responsible_user_name_en" gorm:"type:nvarchar(100)"`
	ResponsibleUserFirstNameAr  string     `json:"responsible_user_name_ar" gorm:"type:nvarchar(100)"`
	PublishedDate               *time.Time `json:"published_date" gorm:"type:date"`
	NextReviewDate              *time.Time `json:"next_review_date" gorm:"type:date"`
	Version                     string     `json:"version" gorm:"type:varchar(20)"`
	DocumentStatusTitleEn       string     `json:"document_status_title_en" gorm:"type:nvarchar(100)"`
	DocumentStatusTitleAr       string     `json:"document_status_title_ar" gorm:"type:nvarchar(100)"`
	DocumentStatusColorCode     string     `json:"document_status_color_code"`
	DocumentAccessTypeTitleEn   string     `json:"document_access_type_title_en" gorm:"type:nvarchar(100)"`
	DocumentAccessTypeTitleAr   string     `json:"document_access_type_title_ar" gorm:"type:nvarchar(100)"`
	DocumentAccessTypeColorCode string     `json:"document_access_type_color_code" gorm:"type:nvarchar(100)"`
	CreatedByID                 uint       `json:"created_by_id" gorm:"not null;index"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"not null"`

	// Relationships
	PolicyReport *PolicyReport `gorm:"foreignKey:PolicyReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedBy    *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PolicyReportOverview) TableName() string {
	return "policy_report_overviews"
}
