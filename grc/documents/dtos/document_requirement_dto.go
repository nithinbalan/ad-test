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

type DocumentRequirementGroupWithRequirementsDTO struct {
	RequirementGroupTitleEn string                       `json:"requirement_group_title_en" binding:"required"`
	RequirementGroupTitleAr string                       `json:"requirement_group_title_ar" binding:"required"`
	Requirements            []DocumentRequirementItemDTO `json:"requirements" binding:"required,dive"`
}

type DocumentRequirementItemDTO struct {
	RequirementTitleEn      string `json:"requirement_title_en" binding:"required"`
	RequirementTitleAr      string `json:"requirement_title_ar" binding:"required"`
	Type                    string `json:"type"`
	Source                  string `json:"source"`
	ResponsibleDepartment   string `json:"responsible_department"`
	Evidence                string `json:"evidence"`
	Frequency               string `json:"frequency"`
	ResponsibleDepartmentID uint   `json:"responsible_department_id"`
	FrequencyID             uint   `json:"frequency_id"`
}
