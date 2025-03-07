package dtos

type SwotPestelGroupDTO struct {
	ID                        int    `json:"id"`
	SwotPestelQuestionGroupID uint   `json:"swot_pestel_question_group_id"`
	TitleEn                   string `json:"title_en"`
	TitleAr                   string `json:"title_ar"`
	Order                     int    `json:"order"`
	Status                    string `json:"status"`
	IsSwot                    bool   `json:"is_swot"`
	IsPestel                  bool   `json:"is_pestel"`
}
