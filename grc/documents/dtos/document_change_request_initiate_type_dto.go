package dtos

type DocumentChangeRequestInitiateTypeDTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Type    string `json:"type"`
	Status  string `json:"status"`
}
