package models

type DistributionStatus struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Type    string `json:"type"`
	Status  string `json:"status"`
}

func (DistributionStatus) TableName() string {
	return "distribution_statuses"
}
