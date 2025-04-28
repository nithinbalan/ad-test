package dtos

type SwotPestelQuestionGroupListDTO struct {
	ID                  uint           `json:"id"`
	TitleEn             string         `json:"title_en"`
	TitleAr             string         `json:"title_ar"`
	Order               int            `json:"order"`
	Status              string         `json:"status"`
	IsSwot              bool           `json:"is_swot"`
	IsPestel            bool           `json:"is_pestel"`
	SwotPestelQuestions []*QuestionDTO `json:"swot_pestel_questions" gorm:"-"`
}

type QuestionDTO struct {
	ID       int    `json:"id"`
	TitleEn  string `json:"title_en"`
	TitleAr  string `json:"title_ar"`
	Order    int    `json:"order"`
	Status   string `json:"status"`
	IsSwot   bool   `json:"is_swot"`
	IsPestel bool   `json:"is_pestel"`
}

type SwotPestelQuestionGroupListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Status   string
	Filters  map[string][]string
	IsSwot   string
	IsPestel string
}

type SwotPestelQuestionGroupDTO struct {
	ID                  uint                     `json:"id"`
	TitleEn             string                   `json:"title_en"`
	TitleAr             string                   `json:"title_ar"`
	Order               int                      `json:"order"`
	Status              string                   `json:"status"`
	IsSwot              bool                     `json:"is_swot"`
	IsPestel            bool                     `json:"is_pestel"`
	SwotPestelQuestions []*SwotPestelQuestionDTO `json:"swot_pestel_questions"`
}
