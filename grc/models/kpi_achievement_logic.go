package models

type KPIAchievementLogic struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Type    string `json:"type" binding:"required,max=25,oneof=higher lower" gorm:"type:varchar(25);not null"`
	Status  string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

func (KPIAchievementLogic) TableName() string {
	return "kpi_achievement_logics"
}
