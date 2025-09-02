package dtos

type PolicyFolderListDTO struct {
	ID                  uint   `json:"id"`
	TitleEn             string `json:"title_en"`
	TitleAr             string `json:"title_ar"`
	Status              string `json:"status"`
	CreatedAt           string `json:"created_at"`
	PolicyRegisterCount int    `json:"policy_register_count"`
}

type PolicyFolderDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

type PolicyFolderFormDTO struct {
	TitleEn          string `json:"title_en" binding:"required"`
	TitleAr          string `json:"title_ar" binding:"required"`
	ParentDocumentID *uint  `json:"parent_document_id"`
}
