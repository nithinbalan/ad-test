package models

type AssetDocumentCategory struct {
	ID            uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn       string `json:"title_en" binding:"required,max=250" gorm:"type:nvarchar(250);not null"`
	TitleAr       string `json:"title_ar" binding:"required,max=250" gorm:"type:nvarchar(250);not null"`
	DescriptionEn string `json:"description_en" gorm:"type:text"`
	DescriptionAr string `json:"description_ar" gorm:"type:text"`
	Status        string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active'"`

	// Relations
}

// TableName overrides the table name used by AssetDocumentCategory to `asset_document_categories`
func (AssetDocumentCategory) TableName() string {
	return "asset_document_categories"
}
