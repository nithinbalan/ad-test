package dtos

type UsersByStatusDTO struct {
	ID                  uint   `json:"id"`
	UserStatusTitleEn   string `json:"user_status_title_en"`
	UserStatusTitleAr   string `json:"user_status_title_ar"`
	UserStatusColorCode string `json:"user_status_color_code"`
	UserStatusType      string `json:"user_status_type"`
	Count               int    `json:"count"`
}
