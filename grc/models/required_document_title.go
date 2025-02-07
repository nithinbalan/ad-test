package models

type RequiredDocumentTitle struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(255);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(255);not null"`
	Status  string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the default table name
func (RequiredDocumentTitle) TableName() string {
	return "required_document_titles"
}
