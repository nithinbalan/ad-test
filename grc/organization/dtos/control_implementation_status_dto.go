package dtos

type ControlImplementationStatusDTO struct {
	ID        int    `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	Status    string `json:"status"`
	ColorCode string `json:"color_code"`
}
