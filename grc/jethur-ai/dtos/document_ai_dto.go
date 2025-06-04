package dtos

type FAQEntry struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type DocumentAISummaryDTO struct {
	Summary         *string    `json:"summary"`
	SampleQuestions []string   `json:"sample_questions"`
	KeyPoints       []string   `json:"key_points"`
	FAQ             []FAQEntry `json:"faq"`
}
