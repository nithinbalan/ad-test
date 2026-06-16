package dtos

type AwarenessCampaignTypeFormDTO struct {
	TitleEn string `json:"title_en" binding:"required,max=100"`
	TitleAr string `json:"title_ar" binding:"required,max=100"`
	Status  string `json:"status" binding:"required,oneof=active inactive"`
}