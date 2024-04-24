package dtos

type NotificationListDTO struct {
	ID         uint   `json:"id"`
	ModulesID  uint   `json:"modules_id"`
	Identifier string `json:"identifier"`
	TitleEn    string `json:"title_en"`
	TitleAr    string `json:"title_ar"`
}

type NotificationListResponseDTO struct {
	ID            uint                   `json:"id"`
	TitleEn       string                 `json:"title_en"`
	TitleAr       string                 `json:"title_ar"`
	Notifications *[]NotificationListDTO `json:"notifications"`
}
