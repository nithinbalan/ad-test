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
	Category     string                       `json:"category" binding:"required"`
	CategoryAr   string                       `json:"category_ar"`
	Requirements []DocumentRequirementItemDTO `json:"requirements" binding:"required,dive"`
}

type DocumentRequirementItemDTO struct {
	Statement               string `json:"statement" binding:"required"`
	StatementAr             string `json:"statement_ar"`
	Type                    string `json:"type"`
	Source                  string `json:"source"`
	ResponsibleDepartment   string `json:"responsible_department"`
	Evidence                string `json:"evidence"`
	Frequency               string `json:"frequency"`
	ResponsibleDepartmentID uint   `json:"responsible_department_id"`
	FrequencyID             uint   `json:"frequency_id"`
}
