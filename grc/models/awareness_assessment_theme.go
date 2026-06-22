package models

type AwarenessAssessmentTheme struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr string `json:"title_ar" gorm:"type:nvarchar(255);not null"`
	Status  string `json:"status" gorm:"type:varchar(10);not null"`
}

func (AwarenessAssessmentTheme) TableName() string {
	return "awareness_assessment_themes"
}