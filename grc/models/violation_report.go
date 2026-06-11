package models

import (
	"time"
)

type ViolationReport struct {
	ID                           uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	ReferenceNumber              string     `json:"reference_number" binding:"required,max=45" gorm:"size:45;unique;not null"`
	ReportedAt                   time.Time  `json:"reported_at" binding:"required" gorm:"type:date;not null"`
	ViolationReportChannelID     uint       `json:"violation_report_channel_id" binding:"required"`
	IDOfReporter                 *uint      `json:"id_of_reporter"`
	NameOfReporter               *string    `json:"name_of_reporter" gorm:"type:nvarchar(max)"`
	PhoneOfReporter              *string    `json:"phone_of_reporter" gorm:"type:nvarchar(max)"`
	EmailOfReporter              *string    `json:"email_of_reporter" gorm:"type:nvarchar(max)"`
	TitleEn                      string     `json:"title_en" binding:"required" gorm:"type:nvarchar(max);not null"`
	TitleAr                      string     `json:"title_ar" binding:"required" gorm:"type:nvarchar(max);not null"`
	ReportedEntity               *string    `json:"reported_entity" gorm:"type:nvarchar(max)"`
	DescriptionEn                *string    `json:"description_en" gorm:"type:nvarchar(max)"`
	DescriptionAr                *string    `json:"description_ar" gorm:"type:nvarchar(max)"`
	EvidenceDescription          *string    `gorm:"type:nvarchar(max)" json:"evidence_description"`
	ViolationReportCategoryID    *uint      `json:"violation_report_category_id"`
	ViolationReportSubCategoryID *uint      `json:"violation_report_sub_category_id"`
	ViolationReportReasonID      uint       `json:"violation_report_reason_id" binding:"required"`
	SeverityLevelID              *uint      `json:"severity_level_id"`
	ViolationReportStatusID      uint       `json:"violation_report_status_id" binding:"required"`
	CreatedBYID                  uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID                  *uint      `json:"updated_by_id"`
	CreatedAt                    time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                    *time.Time `gorm:"type:datetime"`

	ViolationReportChannel     *ViolationReportChannel     `gorm:"foreignKey:ViolationReportChannelID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	ViolationReportCategory    *ViolationReportCategory    `gorm:"foreignKey:ViolationReportCategoryID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	ViolationReportSubCategory *ViolationReportSubCategory `gorm:"foreignKey:ViolationReportSubCategoryID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	ViolationReportReason      *ViolationReportReason      `gorm:"foreignKey:ViolationReportReasonID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	ViolationReportStatus      *ViolationReportStatus      `gorm:"foreignKey:ViolationReportStatusID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	SeverityLevel              *SeverityLevel              `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	Documents                  []*ViolationReportDocument  `gorm:"foreignKey:ViolationReportID"`
	CreatedByUser              *User                       `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	UpdatedByUser              *User                       `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by ViolationReport to `violation_reports`
func (ViolationReport) TableName() string {
	return "violation_reports"
}
