package models

type AIModule struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn string `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr string `json:"title_ar,omitempty" gorm:"type:nvarchar(255)"`
	Order   int    `json:"order" gorm:"default:0"`

	Features []AIModuleFeature `json:"features,omitempty" gorm:"foreignKey:AIModuleID"`
}

func (AIModule) TableName() string {
	return "ai_modules"
}
