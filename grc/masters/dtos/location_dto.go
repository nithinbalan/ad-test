package dtos

type LocationDTO struct {
	ID      uint    `json:"id"`
	TitleEn string  `json:"title_en"`
	TitleAr string  `json:"title_ar"`
	Address *string `json:"address"`
	Status  string  `json:"status"`
}
