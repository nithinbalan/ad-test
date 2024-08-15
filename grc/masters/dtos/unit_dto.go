package dtos

type UnitDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	CodeEn  string `json:"code_en"`
	CodeAr  string `json:"code_ar"`
	Status  string `json:"status"`
}
