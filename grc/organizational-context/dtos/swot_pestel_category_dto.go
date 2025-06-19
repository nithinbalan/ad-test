package dtos

type SwotPestelCategoryDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	IsSwot    bool   `json:"is_swot"`
	IsPestel  bool   `json:"is_pestel"`
	Status    string `json:"status"`
}
