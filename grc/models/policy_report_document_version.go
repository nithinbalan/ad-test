package models

import "time"

type PolicyReportDocumentVersion struct {
	ID                  uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	PolicyReportID      uint       `json:"policy_report_id" gorm:"not null;index"`
	Version             int        `json:"version" gorm:"not null;"`
	PublishedAt         *time.Time `json:"published_at" gorm:"type:date"`
	PublishedUserNameEn string     `json:"published_user_name_en" gorm:"type:nvarchar(255)"`
	PublishedUserNameAr string     `json:"published_user_name_ar" gorm:"type:nvarchar(255)"`
	PublishedByEmail    string     `json:"published_by_email" gorm:"type:varchar(255)"`
	CreatedByID         uint       `json:"created_by_id" gorm:"not null;index"`
	CreatedAt           time.Time  `json:"created_at" gorm:"not null"`

	// Relationships
	PolicyReport  *PolicyReport `gorm:"foreignKey:PolicyReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PolicyReportDocumentVersion) TableName() string {
	return "policy_report_document_versions"
}
