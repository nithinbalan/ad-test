package models

type AuditFinalReportStatus struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Type      string `json:"type" binding:"required,oneof=draft in-review submit revert reject published completed approve" gorm:"type:varchar(20);not null;check:type IN ('draft', 'in-review', 'submit', 'revert', 'reject', 'published', 'completed','approve')"`
	ColorCode string `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Status    string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active'"`
}

// TableName overrides the table name used by AuditFinalReportStatus to `audit_final_report_statuses`
func (AuditFinalReportStatus) TableName() string {
	return "audit_final_report_statuses"
}
