package dtos

type DocumentRequirementGroupDTO struct {
	ID                   uint                      `json:"id"`
	DocumentID           uint                      `json:"document_id"`
	TitleEn              string                    `json:"title_en"`
	TitleAr              string                    `json:"title_ar"`
	Status               string                    `json:"status"`
	DocumentRequirements []*DocumentRequirementDTO `json:"document_requirements"`
}
type DocumentRequirementGroupFormDTO struct {
	TitleEn string `json:"title_en" binding:"required,max=250"`
	TitleAr string `json:"title_ar" binding:"required,max=250"`
	Status  string `json:"status" binding:"required,oneof=active inactive"`
}
