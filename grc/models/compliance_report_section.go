package models

type ComplianceReportSection struct {
	ID                 uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceReportID uint    `json:"compliance_report_id" gorm:"not null" binding:"required"`
	Type               string  `json:"type" gorm:"type:nvarchar(50);not null;check:type IN ('executive-summary', 'compliance-status-by-department', 'compliance-status-by-standard', 'action-plan-implementation', 'conclusion-recommendations')"`
	TitleEn            string  `json:"title_en" gorm:"type:nvarchar(500);not null" binding:"required"`
	TitleAr            string  `json:"title_ar" gorm:"type:nvarchar(500);not null" binding:"required"`
	DescriptionEn      *string `json:"description_en" gorm:"type:nvarchar(MAX);default:null"`
	DescriptionAr      *string `json:"description_ar" gorm:"type:nvarchar(MAX);default:null"`
	OrderNum           int     `json:"order_num" gorm:"not null" binding:"required"`

	ComplianceReport ComplianceReport `gorm:"foreignKey:ComplianceReportID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ComplianceReportSection) TableName() string {
	return "compliance_report_sections"
}
