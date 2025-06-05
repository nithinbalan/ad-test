package dtos

type FAQItemDTO struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type DocumentAISummaryDTO struct {
	SummaryEn         *string      `json:"summary_en"`
	SummaryAr         *string      `json:"summary_ar"`
	SampleQuestionsEn []string     `json:"sample_questions_en"`
	SampleQuestionsAr []string     `json:"sample_questions_ar"`
	KeyPointsEn       []string     `json:"key_points_en"`
	KeyPointsAr       []string     `json:"key_points_ar"`
	FAQEn             []FAQItemDTO `json:"faq_en"`
	FAQAr             []FAQItemDTO `json:"faq_ar"`
}
