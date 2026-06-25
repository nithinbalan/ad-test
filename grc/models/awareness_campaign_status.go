package models

type AwarenessCampaignStatus struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" gorm:"type:nvarchar(100);not null"`
	TitleAr string `json:"title_ar" gorm:"type:nvarchar(100);not null"`
	Type    string `json:"type" binding:"required,oneof=planned submitted cancelled completed not_started" gorm:"type:varchar(50);not null"`
	Status  string `json:"status" gorm:"type:varchar(10);not null"`
}

func (AwarenessCampaignStatus) TableName() string {
	return "awareness_campaign_statuses"
}
