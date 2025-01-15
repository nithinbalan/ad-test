package dtos

type AuditPlanByCategoryDTO struct {
	AuditPlanCategoryID      uint   `json:"id" gorm:"column:audit_category_id"`
	AuditPlanCategoryTitleEn string `json:"title_en" gorm:"column:audit_category_title_en"`
	AuditPlanCategoryTitleAr string `json:"title_ar" gorm:"column:audit_category_title_ar"`
	Count                    uint   `json:"count" gorm:"column:count"`
}

type AuditPlanByStatusDTO struct {
	AuditStatusID      uint   `json:"id" gorm:"column:audit_status_id"`
	AuditStatusTitleEn string `json:"title_en" gorm:"column:audit_status_title_en"`
	AuditStatusTitleAr string `json:"title_ar" gorm:"column:audit_status_title_ar"`
	ColorCode          string `json:"color_code" gorm:"column:color_code"`
	Count              uint   `json:"count" gorm:"column:count"`
}

type AuditPlanByTeamDTO struct {
	AuditTeamID      uint   `json:"id" gorm:"column:audit_team_id"`
	AuditTeamTitleEn string `json:"title_en" gorm:"column:audit_team_title_en"`
	AuditTeamTitleAr string `json:"title_ar" gorm:"column:audit_team_title_ar"`
	Count            uint   `json:"count" gorm:"column:count"`
}

type AuditPlanByTeamLeadDTO struct {
	TeamLeadID               uint   `json:"id" gorm:"column:team_lead_id"`
	TeamLeadFirstNameEn      string `json:"first_name_en" gorm:"column:team_lead_first_name_en"`
	TeamLeadFirstNameAr      string `json:"first_name_ar" gorm:"column:team_lead_first_name_ar"`
	TeamLeadLastNameEn       string `json:"last_name_en" gorm:"column:team_lead_last_name_en"`
	TeamLeadLastNameAr       string `json:"last_name_ar" gorm:"column:team_lead_last_name_ar"`
	TeamLeadTitleEn          string `json:"title_en" gorm:"column:team_lead_lead_title_en"`
	TeamLeadTitleAr          string `json:"title_ar" gorm:"column:team_lead_lead_title_ar"`
	TeamLeadInitial          string `json:"initial" gorm:"column:team_lead_initial"`
	TeamLeadColor            string `json:"color" gorm:"column:team_lead_color"`
	TeamLeadEmail            string `json:"email" gorm:"column:team_lead_email"`
	TeamLeadProfileImagePath string `json:"profile_image_path" gorm:"column:team_lead_profile_image_path"`
	Count                    uint   `json:"count" gorm:"column:count"`
}

type AuditPlanByYearDTO struct {
	TitleEn string `json:"title_en" gorm:"column:title_en"`
	TitleAr string `json:"title_ar" gorm:"column:title_ar"`
	Count   uint   `json:"count" gorm:"column:count"`
}

type AuditByCategoryDTO struct {
	AuditCategoryID      uint   `json:"id" gorm:"column:audit_category_id"`
	AuditCategoryTitleEn string `json:"title_en" gorm:"column:audit_category_title_en"`
	AuditCategoryTitleAr string `json:"title_ar" gorm:"column:audit_category_title_ar"`
	Count                uint   `json:"count" gorm:"column:count"`
}

type AuditByStatusDTO struct {
	AuditStatusID      uint   `json:"id" gorm:"column:audit_status_id"`
	AuditStatusTitleEn string `json:"title_en" gorm:"column:audit_status_title_en"`
	AuditStatusTitleAr string `json:"title_ar" gorm:"column:audit_status_title_ar"`
	ColorCode          string `json:"color_code" gorm:"column:color_code"`
	Count              uint   `json:"count" gorm:"column:count"`
}

type AuditByOpenVsClosedDTO struct {
	ID        uint   `json:"id" gorm:"column:id"`
	TitleEn   string `json:"title_en" gorm:"column:title_en"`
	TitleAr   string `json:"title_ar" gorm:"column:title_ar"`
	Type      string `json:"type" gorm:"column:type"`
	ColorCode string `json:"color_code" gorm:"column:color_code"`
	Count     uint   `json:"count" gorm:"column:count"`
}

type AuditByMonthDTO struct {
	TitleEn string `json:"title_en" gorm:"column:title_en"` // Month in English
	TitleAr string `json:"title_ar" gorm:"column:title_ar"` // Month in Arabic
	Count   uint   `json:"count" gorm:"column:count"`
}

type AuditByTeamLeadDTO struct {
	TeamLeadID               uint   `json:"id" gorm:"column:team_lead_id"`
	TeamLeadFirstNameEn      string `json:"first_name_en" gorm:"column:team_lead_first_name_en"`
	TeamLeadFirstNameAr      string `json:"first_name_ar" gorm:"column:team_lead_first_name_ar"`
	TeamLeadLastNameEn       string `json:"last_name_en" gorm:"column:team_lead_last_name_en"`
	TeamLeadLastNameAr       string `json:"last_name_ar" gorm:"column:team_lead_last_name_ar"`
	TeamLeadTitleEn          string `json:"title_en" gorm:"column:team_lead_lead_title_en"`
	TeamLeadTitleAr          string `json:"title_ar" gorm:"column:team_lead_lead_title_ar"`
	TeamLeadInitial          string `json:"initial" gorm:"column:team_lead_initial"`
	TeamLeadColor            string `json:"color" gorm:"column:team_lead_color"`
	TeamLeadEmail            string `json:"email" gorm:"column:team_lead_email"`
	TeamLeadProfileImagePath string `json:"profile_image_path" gorm:"column:team_lead_profile_image_path"`
	Count                    uint   `json:"count" gorm:"column:count"`
}

type AuditFindingByCategoryDTO struct {
	AuditFindingCategoryID      uint   `json:"id" gorm:"column:id"`
	AuditFindingCategoryTitleEn string `json:"title_en" gorm:"column:title_en"`
	AuditFindingCategoryTitleAr string `json:"title_ar" gorm:"column:title_ar"`
	Count                       uint   `json:"count" gorm:"column:count"`
}
