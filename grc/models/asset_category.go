package models

type AssetCategory struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Status  string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active'"`

	// Relations
}

// TableName overrides the table name used by AssetCategory to `asset_categories`
func (AssetCategory) TableName() string {
	return "asset_categories"
}
