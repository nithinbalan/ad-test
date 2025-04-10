package dtos

type ComplianceStatusDTO struct {
	ID        int    `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	Type      string `json:"type"`
	ColorCode string `json:"color_code"`
	Status    string `json:"status"`
	ScoreFrom *int   `json:"score_from"`
	ScoreTo   *int   `json:"score_to"`
}
