package models

type UserDocumentSubType struct {
	ID                 uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserDocumentTypeID uint   `gorm:"column:user_document_type_id" json:"user_document_type_id"`
	TitleEn            string `gorm:"type:nvarchar(255);not null" json:"title_en"`
	TitleAr            string `gorm:"type:nvarchar(255);not null" json:"title_ar"`
	Status             string `gorm:"type:varchar(10);default:'active'" json:"status"`

	UserDocumentType *UserDocumentType `gorm:"foreignKey:UserDocumentTypeID" json:"user_document_type,omitempty"`
}

func (UserDocumentSubType) TableName() string {
	return "user_document_sub_types"
}
