package models

type PenetrationTestType struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Code    string `json:"code" binding:"required,max=20" gorm:"type:nvarchar(20);not null"`
	TitleEn string `json:"title_en" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	Status  string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active'"`
}

func (PenetrationTestType) TableName() string {
	return "penetration_test_types"
}
