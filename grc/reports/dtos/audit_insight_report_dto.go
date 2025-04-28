package dtos

type AuditInsightAuditByStatusDTO struct {
	AuditStatusID      uint   `json:"id"`
	AuditStatusTitleEn string `json:"title_en"`
	AuditStatusTitleAr string `json:"title_ar"`
	ColorCode          string `json:"color_code"`
	Count              uint   `json:"count"`
}

type AuditInsightAuditByYearDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}

type AuditInsightAuditByMonthDTO struct {
	TitleEn string `json:"title_en"` // Month in English
	TitleAr string `json:"title_ar"` // Month in Arabic
	Count   uint   `json:"count"`
}

type AuditInsightAuditPlanByCategoryDTO struct {
	AuditCategoryID      uint   `json:"id"`
	AuditCategoryTitleEn string `json:"title_en"`
	AuditCategoryTitleAr string `json:"title_ar"`
	Count                uint   `json:"count"`
}

type AuditInsightAuditPlanByTeamLeadDTO struct {
	TeamLeadID               uint   `json:"id"`
	TeamLeadFirstNameEn      string `json:"first_name_en"`
	TeamLeadFirstNameAr      string `json:"first_name_ar"`
	TeamLeadLastNameEn       string `json:"last_name_en"`
	TeamLeadLastNameAr       string `json:"last_name_ar"`
	TeamLeadTitleEn          string `json:"title_en"`
	TeamLeadTitleAr          string `json:"title_ar"`
	TeamLeadInitial          string `json:"initial"`
	TeamLeadColor            string `json:"color"`
	TeamLeadEmail            string `json:"email"`
	TeamLeadProfileImagePath string `json:"profile_image_path"`
	Count                    uint   `json:"count"`
}

type AuditInsightAuditPlanByAuditableItemDTO struct {
	AuditableItemID      uint   `json:"id"`
	AuditableItemTitleEn string `json:"title_en"`
	AuditableItemTitleAr string `json:"title_ar"`
	Count                uint   `json:"count"`
}

type AuditInsightAuditPlanByMonthDTO struct {
	TitleEn string `json:"title_en"` // Month in English
	TitleAr string `json:"title_ar"` // Month in Arabic
	Count   uint   `json:"count"`
}

type AuditInsightAuditFindingByCategoryDTO struct {
	AuditFindingCategoryID      uint   `json:"id"`
	AuditFindingCategoryTitleEn string `json:"title_en"`
	AuditFindingCategoryTitleAr string `json:"title_ar"`
	Count                       uint   `json:"count"`
}

type AuditInsightAuditFindingByAuditCategoryDTO struct {
	AuditCategoryID      uint   `json:"id"`
	AuditCategoryTitleEn string `json:"title_en"`
	AuditCategoryTitleAr string `json:"title_ar"`
	Count                uint   `json:"count"`
}

type AuditInsightAuditPlanVsAuditDTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   int64  `json:"count"`
}

type AuditInsightAuditFindingByStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type AuditInsightAuditFindingBySeverityLevelDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type AuditInsightAuditFindingByMonthDTO struct {
	TitleEn string `json:"title_en"` // Month in English
	TitleAr string `json:"title_ar"` // Month in Arabic
	Count   uint   `json:"count"`
}

type AuditInsightAuditFindingByYearDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}

type AuditInsightAuditActionPlanByStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type AuditInsightAuditActionPlanByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"id"`
	ResponsibleUserFirstNameEn      string `json:"first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"last_name_en"`
	ResponsibleUserLastNameAr       string `json:"last_name_ar"`
	ResponsibleUserTitleEn          string `json:"title_en"`
	ResponsibleUserTitleAr          string `json:"title_ar"`
	ResponsibleUserInitial          string `json:"initial"`
	ResponsibleUserColor            string `json:"color"`
	ResponsibleUserEmail            string `json:"email"`
	ResponsibleUserProfileImagePath string `json:"profile_image"`
	Count                           int64  `json:"count"`
}

type AuditInsightAuditDraftReportByStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type AuditInsightAuditFinalReportByStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type AuditInsightOverdueActionPlanCountDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   int64  `json:"count"`
}
