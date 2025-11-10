package models

import "time"

type PolicyReportDocumentAttachment struct {
	ID                  uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PolicyReportID      uint      `json:"policy_report_id" gorm:"not null;index"`
	DocumentTitleEn     string    `json:"document_title_en" gorm:"type:nvarchar(255);not null"`
	DocumentTitleAr     string    `json:"document_title_ar" gorm:"type:nvarchar(255);not null"`
	DocumentTypeTitleEn string    `json:"document_type_title_en" gorm:"type:nvarchar(100);not null"`
	DocumentTypeTitleAr string    `json:"document_type_title_ar" gorm:"type:nvarchar(100);not null"`
	LinkedDate          string    `json:"linked_date" gorm:"type:varchar(50)"`
	CreatedByID         uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt           time.Time `json:"created_at" gorm:"not null"`

	// Relations
	PolicyReport  *PolicyReport `gorm:"foreignKey:PolicyReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PolicyReportDocumentAttachment) TableName() string {
	return "policy_report_document_attachments"
}
