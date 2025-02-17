package dtos

type ComplianceAssessmentStatusDTO struct {
	ID        int    `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	Type      string `json:"type"`
	ColorCode string `json:"color_code"`
	Status    string `json:"status"`
}
