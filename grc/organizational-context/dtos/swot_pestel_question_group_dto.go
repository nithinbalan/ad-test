package dtos

type SwotPestelQuestionGroupDTO struct {
	ID       int    `json:"id"`
	TitleEn  string `json:"title_en"`
	TitleAr  string `json:"title_ar"`
	Order    int    `json:"order"`
	Status   string `json:"status"`
	IsSwot   bool   `json:"is_swot"`
	IsPestel bool   `json:"is_pestel"`
}
