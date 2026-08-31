package models

type PenetrationTestEnvironmentReadiness struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	Status  string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active'"`
}

func (PenetrationTestEnvironmentReadiness) TableName() string {
	return "penetration_test_environment_readinesses"
}
