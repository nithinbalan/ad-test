package models

type ViolationReportReason struct {
	ID            uint    `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	TitleEn       string  `json:"title_en" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	TitleAr       string  `json:"title_ar" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	DescriptionEn *string `json:"description_en" gorm:"type:nvarchar(MAX)"`
	DescriptionAr *string `json:"description_ar" gorm:"type:nvarchar(MAX)"`
	Status        string  `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);default:active"`
}

// TableName overrides the table name used by ViolationReportReason to `violation_report_reasons`
func (ViolationReportReason) TableName() string {
	return "violation_report_reasons"
}
