package dtos

type PenetrationTestTypeDTO struct {
	ID      uint   `json:"id"`
	Code    string `json:"code"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}
