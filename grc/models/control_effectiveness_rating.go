package models

type ControlEffectivenessRating struct {
	ID               uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn          string  `json:"title_en" gorm:"column:title_en;type:nvarchar(255);not null"`
	TitleAr          string  `json:"title_ar" gorm:"column:title_ar;type:nvarchar(255);not null"`
	Type             string  `json:"type" binding:"required,max=100,oneof=control_missing Ineffective improvement_needed partially_effective very_effective" gorm:"type:varchar(100);not null"`
	MinScore         float64 `json:"min_score" gorm:"column:min_score;type:decimal(5,2);not null"`
	MaxScore         float64 `json:"max_score" gorm:"column:max_score;type:decimal(5,2);not null"`
	ReductionPercent float64 `json:"reduction_percent" gorm:"column:reduction_percent;type:decimal(5,1);not null"`
	Status           string  `json:"status" gorm:"column:status;type:varchar(20);not null"`
}

func (ControlEffectivenessRating) TableName() string {
	return "control_effectiveness_ratings"
}
