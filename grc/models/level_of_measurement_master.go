package models

type LevelOfMeasurement struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Type    string `json:"type"`
	Status  string `json:"status"`
}

func (LevelOfMeasurement) TableName() string {
	return "level_of_measurements"
}
