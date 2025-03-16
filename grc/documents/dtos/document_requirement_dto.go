package dtos

type DocumentRequirementDTO struct {
	ID                         uint   `json:"id"`
	DocumentRequirementGroupID uint   `json:"document_requirement_group_id"`
	TitleEn                    string `json:"title_en"`
	TitleAr                    string `json:"title_ar"`
}
type DocumentRequirementFormDTO struct {
	TitleEn string `json:"title_en" binding:"required,max=250"`
	TitleAr string `json:"title_ar" binding:"required,max=250"`
}
