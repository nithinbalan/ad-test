package dtos

type SwotGroupDTO struct {
	ID                  int    `json:"id"`
	SwotQuestionGroupID uint   `json:"swot_question_group_id"`
	TitleEn             string `json:"title_en"`
	TitleAr             string `json:"title_ar"`
	Order               int    `json:"order"`
	Status              string `json:"status"`
}
