package dtos

type UserDocumentSubTypeDTO struct {
	ID                      int    `json:"id"`
	UserDocumentTypeID      *uint  `json:"user_document_type_id"`
	TitleEn                 string `json:"title_en"`
	TitleAr                 string `json:"title_ar"`
	UserDocumentTypeTitleEn string `json:"user_document_type_title_en"`
	UserDocumentTypeTitleAr string `json:"user_document_type_title_ar"`
	Status                  string `json:"status"`
}

type UserDocumentSubTypeFormDTO struct {
	ID                 uint   `json:"id"`
	UserDocumentTypeID uint   `json:"user_document_type_id" binding:"required"`
	TitleEn            string `json:"title_en" binding:"required,max=255" gorm:"type:nvarchar(100);not null"`
	TitleAr            string `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(100);not null"`
	Status             string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

type UserDocumentSubTypeListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}

type UserDocumentSubTypeDetailDTO struct {
	ID                 int                 `json:"id"`
	UserDocumentTypeID uint                `json:"user_document_type_id"`
	TitleEn            string              `json:"title_en"`
	TitleAr            string              `json:"title_ar"`
	UserDocumentType   UserDocumentTypeDTO `json:"user_document_type"`
	Status             string              `json:"status"`
}
