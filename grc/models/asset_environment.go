package models

type AssetEnvironment struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" gorm:"type:nvarchar(100);not null"`
	TitleAr string `json:"title_ar" gorm:"type:nvarchar(100);not null"`
	Type    string `json:"type" gorm:"type:varchar(150);not null"`
	Status  string `json:"status" gorm:"type:varchar(10);not null;default:active"`
}

func (AssetEnvironment) TableName() string {
	return "asset_environments"
}