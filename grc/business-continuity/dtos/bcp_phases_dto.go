package dtos

type BcpPhaseDTO struct {
	ID       int    `json:"id"`
	TitleEn  string `json:"title_en"`
	TitleAr  string `json:"title_ar"`
	TimeFram string `json:"time_frame"`
	Status   string `json:"status"`
}
