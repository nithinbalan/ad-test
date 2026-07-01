package models

type AssetDeviceCategory struct {
	ID       uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn  string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr  string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	IsSystem bool   `json:"is_system" gorm:"column:is_system;not null;default:false"`
	Status   string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null;default:active"`
}

func (AssetDeviceCategory) TableName() string {
	return "asset_device_categories"
}