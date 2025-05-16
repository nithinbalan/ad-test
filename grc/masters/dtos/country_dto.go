package dtos

type CountryDTO struct {
	ID        uint    `json:"id"`
	TitleEn   string  `json:"title_en"`
	TitleAr   string  `json:"title_ar"`
	Code      string  `json:"code"`
	PhoneCode *string `json:"phone_code"`
	Status    string  `json:"status"`
}
