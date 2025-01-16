package models

import (
	"time"
)

type AuditFinalReportSection struct {
	ID                        uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditFinalReportID        uint       `json:"audit_final_report_id" gorm:"not null"`
	TitleEn                   string     `json:"title_en" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	TitleAr                   string     `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	DescriptionEn             *string    `json:"description_en" gorm:"type:nvarchar(MAX);default:null"`
	DescriptionAr             *string    `json:"description_ar" gorm:"type:nvarchar(MAX);default:null"`
	Type                      string     `json:"type" binding:"required,oneof=overview finding-summary findings conclusion custom appendix" gorm:"type:varchar(20);not null;check:type IN ('overview', 'finding-summary', 'findings', 'conclusion', 'custom', 'appendix')"`
	Order                     *float64   `json:"order" gorm:"type:decimal(10,2);default:null"`
	AuditFinalReportSectionID *uint      `json:"audit_final_report_section_id" gorm:"default:null"`
	CreatedAt                 time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                 *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID               uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID               *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	AuditFinalReport *AuditFinalReport        `gorm:"foreignKey:AuditFinalReportID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ParentSection    *AuditFinalReportSection `gorm:"foreignKey:AuditFinalReportSectionID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser    *User                    `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser    *User                    `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by AuditFinalReportSection to `audit_final_report_sections`
func (AuditFinalReportSection) TableName() string {
	return "audit_final_report_sections"
}
