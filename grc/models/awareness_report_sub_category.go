package models

type AwarenessReportSubCategory struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	TitleEn string `gorm:"column:title_en;type:nvarchar(100);not null" json:"title_en"`
	TitleAr string `gorm:"column:title_ar;type:nvarchar(100);not null" json:"title_ar"`
	Status  string `gorm:"column:status;type:varchar(10);default:active;not null" json:"status"`
}

func (AwarenessReportSubCategory) TableName() string {
	return "awareness_report_sub_categories"
}