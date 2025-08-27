package dtos

type ReminderContextDTO struct {
	ID      uint   `json:"id"`
	Type    string `json:"type"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}
