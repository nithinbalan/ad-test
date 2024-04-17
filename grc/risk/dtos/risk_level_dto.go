package dtos

type RiskLevelDTO struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Status    string `json:"status"`
}
