package models

type DocumentChangeRequestType struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn string `json:"title_en" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	Type    string `json:"type" binding:"required,oneof=minor major archive" gorm:"type:varchar(45);not null;check:type IN ('minor', 'major', 'archive')"`
	Status  string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(45);default:'active';not null;check:status IN ('active', 'inactive')"`
}

// TableName overrides the table name used by GORM for DocumentChangeRequestType
func (DocumentChangeRequestType) TableName() string {
	return "document_change_request_types"
}
