package dtos

type BusinessContinuityDurationUnitDTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Minutes int    `json:"minutes"`
	Type    string `json:"type"`
	Status  string `json:"status"`
}
