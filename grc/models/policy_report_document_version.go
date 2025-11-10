package models

import "time"

type PolicyReportDocumentVersion struct {
	ID                  uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PolicyReportID      uint      `json:"policy_report_id" gorm:"not null;index"`
	DocumentVersion     float64   `json:"document_version" gorm:"type:decimal(4,2);"`
	PublishedAt         *string   `json:"published_at" gorm:"type:varchar(50)"`
	PublishedUserNameEn *string   `json:"published_user_name_en" gorm:"type:nvarchar(255)"`
	PublishedUserNameAr *string   `json:"published_user_name_ar" gorm:"type:nvarchar(255)"`
	CreatedByID         uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt           time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	PolicyReport  *PolicyReport `gorm:"foreignKey:PolicyReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PolicyReportDocumentVersion) TableName() string {
	return "policy_report_document_versions"
}
