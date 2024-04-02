package dtos

type ComplianceStatusDTO struct {
	ID        int    `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Status    string `json:"status"`
}
