package dtos

type AuditWorkflowStatusDTO struct {
	ID          uint   `json:"id"`
	TitleEn     string `json:"title_en"`
	TitleAr     string `json:"title_ar"`
	ColorCode   string `json:"color_code"`
	Type        string `json:"type"`
}
