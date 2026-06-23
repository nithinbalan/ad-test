package models

type AwarenessHandlingStatus struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	TitleEn string `gorm:"column:title_en;type:nvarchar(100);not null" json:"title_en"`
	TitleAr string `gorm:"column:title_ar;type:nvarchar(100);not null" json:"title_ar"`
	Type    string `gorm:"column:type;type:varchar(50);not null" json:"type"`
	Status  string `gorm:"column:status;type:varchar(10);default:active;not null" json:"status"`
}

func (AwarenessHandlingStatus) TableName() string {
	return "awareness_handling_statuses"
}
