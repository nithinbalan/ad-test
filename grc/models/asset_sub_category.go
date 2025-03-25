package models

type AssetSubCategory struct {
	ID              uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AssetCategoryID uint   `json:"asset_category_id" gorm:"not null"  binding:"required"`
	TitleEn         string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr         string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Status          string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active'"`

	// Relations
	AssetCategory *AssetCategory `gorm:"foreignKey:AssetCategoryID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by AssetSubCategory to `asset_categories`
func (AssetSubCategory) TableName() string {
	return "asset_sub_categories"
}
