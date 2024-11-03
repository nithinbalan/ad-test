package models

type AssetCategory struct {
	ID          uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AssetTypeID *uint  `json:"asset_type_id" gorm:"default:null"`
	TitleEn     string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr     string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Status      string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active'"`

	// Relations
	AssetType *AssetType `gorm:"foreignKey:AssetTypeID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by AssetCategory to `asset_categories`
func (AssetCategory) TableName() string {
	return "asset_categories"
}
