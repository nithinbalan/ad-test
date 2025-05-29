package dtos

type JobDescriptionDTO struct {
	ID                       int    `json:"id"`
	TitleEn                  string `json:"title_en"`
	TitleAr                  string `json:"title_ar"`
	DescriptionEn            string `json:"description_en"`
	DescriptionAr            string `json:"description_ar"`
	RolesAndResponsibilities string `json:"roles_and_responsibilities"`
	Status                   string `json:"status"`
}
