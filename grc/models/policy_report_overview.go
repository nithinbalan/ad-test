package models

import "time"

type PolicyReportOverview struct {
	ID                          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PolicyReportID              uint      `json:"policy_report_id" gorm:"not null;index"`
	ReferenceNumber             string    `json:"reference_number" gorm:"type:varchar(50);index"`
	PolicyTitleEn               string    `json:"policy_title_en" gorm:"type:nvarchar(255);not null"`
	PolicyTitleAr               string    `json:"policy_title_ar" gorm:"type:nvarchar(255)"`
	DescriptionEn               string    `json:"description_en" gorm:"type:nvarchar(max)"`
	DescriptionAr               string    `json:"description_ar" gorm:"type:nvarchar(max)"`
	ResponsibleUnitTitleEn      string    `json:"responsible_unit_title_en" gorm:"type:nvarchar(255)"`
	ResponsibleUnitTitleAr      string    `json:"responsible_unit_title_ar" gorm:"type:nvarchar(255)"`
	ResponsibleUserNameEn       string    `json:"responsible_user_name_en" gorm:"type:nvarchar(255)"`
	ResponsibleUserNameAr       string    `json:"responsible_user_name_ar" gorm:"type:nvarchar(255)"`
	PublishedDate               *string   `json:"published_date" gorm:"type:varchar(50)"`
	NextReviewDate              *string   `json:"next_review_date" gorm:"type:varchar(50)"`
	DocumentVersion             float64   `json:"document_version" gorm:"type:decimal(4,2);"`
	DocumentStatusTitleEn       string    `json:"document_status_title_en" gorm:"type:nvarchar(100)"`
	DocumentStatusTitleAr       string    `json:"document_status_title_ar" gorm:"type:nvarchar(100)"`
	DocumentStatusColorCode     string    `json:"document_status_color_code"`
	DocumentAccessTypeTitleEn   string    `json:"document_access_type_title_en" gorm:"type:nvarchar(100)"`
	DocumentAccessTypeTitleAr   string    `json:"document_access_type_title_ar" gorm:"type:nvarchar(100)"`
	DocumentAccessTypeColorCode string    `json:"document_access_type_color_code" gorm:"type:varchar(9)"`
	CreatedByID                 uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt                   time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	PolicyReport  *PolicyReport `gorm:"foreignKey:PolicyReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PolicyReportOverview) TableName() string {
	return "policy_report_overviews"
}
