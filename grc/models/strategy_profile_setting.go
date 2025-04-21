package models

type StrategyProfileSetting struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn  string `json:"title_en" gorm:"size:500;not null" binding:"required"`
	TitleAr  string `json:"title_ar" gorm:"size:500;not null" binding:"required"`
	IsEnable bool   `json:"is_enable"`
}

func (StrategyProfileSetting) TableName() string {
	return "strategy_profile_settings"
}
