package models

type RagSetting struct {
	ID                   uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ThresholdDays        int    `json:"threshold_days" gorm:"not null"`
	OnTrackColorCode     string `json:"on_track_color_code" gorm:"type:varchar(25);not null"`
	ApproachingColorCode string `json:"approaching_color_code" gorm:"type:varchar(25);not null"`
	OverdueColorCode     string `json:"overdue_color_code" gorm:"type:varchar(25);not null"`
}

func (RagSetting) TableName() string {
	return "rag_settings"
}
