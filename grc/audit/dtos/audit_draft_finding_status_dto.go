package dtos

type AuditDraftFindingStatusDTO struct {
	ID          uint   `json:"id"`
	TitleEn     string `json:"title_en"`
	TitleAr     string `json:"title_ar"`
	ColorCode   string `json:"color_code"`
	Type        string `json:"type"`
	IsFinalized bool   `json:"is_finalized"`
}
