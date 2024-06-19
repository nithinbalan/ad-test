package dtos

type ControlEfficiencyStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	Status    string `json:"status"`
	ColorCode string `json:"color_code"`
	ScoreFrom uint   `json:"score_from"`
	ScoreTo   uint   `json:"score_to"`
}
