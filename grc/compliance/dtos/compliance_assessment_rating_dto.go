package dtos

type ComplianceAssessmentRatingDTO struct {
	ID        uint     `json:"id"`
	TitleEn   string   `json:"title_en"`
	TitleAr   string   `json:"title_ar"`
	ColorCode string   `json:"color_code"`
	Type      string   `json:"type"`
	ScoreFrom *float64 `json:"score_from"`
	ScoreTo   *float64 `json:"score_to"`
	Status    string   `json:"status"`
}
