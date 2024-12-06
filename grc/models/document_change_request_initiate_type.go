package models

type DocumentChangeRequestInitiateType struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	Type    string `json:"type" binding:"required,max=45,oneof=modify archive restore" gorm:"type:varchar(45);not null"`
	Status  string `json:"status" binding:"required,max=45,oneof=active inactive" gorm:"type:varchar(45);not null;default:'active'"`
}

// TableName overrides the table name used by DocumentChangeRequestInitiateType to `document_change_request_initiate_types`
func (DocumentChangeRequestInitiateType) TableName() string {
	return "document_change_request_initiate_types"
}
