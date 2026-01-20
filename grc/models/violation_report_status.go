package models

type ViolationReportStatus struct {
	ID        uint   `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	ColorCode string `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" binding:"required,oneof=new under_investigation needs_clarification in_progress under_review on_hold closed_rejected closed_resolved closed_irrelevant closed_no_violation_found closed_referred_to_entity" gorm:"type:varchar(100);not null"`
	Status    string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);default:active"`
}

// TableName overrides the table name used by ViolationReportStatus to `violation_report_statuses`
func (ViolationReportStatus) TableName() string {
	return "violation_report_statuses"
}
