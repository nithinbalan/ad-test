package dtos

type AuditFinalReportSectionFormDTO struct {
	TitleEn                   string  `json:"title_en" binding:"required,max=255"`                                    // English title with a max length of 255 characters
	TitleAr                   string  `json:"title_ar" binding:"required,max=255"`                                    // Arabic title with a max length of 255 characters
	DescriptionEn             *string `json:"description_en" binding:"omitempty,max=10000"`                           // Optional English description with a max length of 10,000 characters
	DescriptionAr             *string `json:"description_ar" binding:"omitempty,max=10000"`                           // Optional Arabic description with a max length of 10,000 characters
}
