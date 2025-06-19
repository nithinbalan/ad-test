package dtos

type AuditPlanByCategoryDTO struct {
	AuditCategoryID      uint   `json:"id"`
	AuditCategoryTitleEn string `json:"title_en"`
	AuditCategoryTitleAr string `json:"title_ar"`
	Count                uint   `json:"count"`
}

type AuditPlanByStatusDTO struct {
	AuditStatusID      uint   `json:"id"`
	AuditStatusTitleEn string `json:"title_en"`
	AuditStatusTitleAr string `json:"title_ar"`
	ColorCode          string `json:"color_code"`
	Count              uint   `json:"count"`
}

type AuditPlanByTeamDTO struct {
	AuditTeamID      uint   `json:"id"`
	AuditTeamTitleEn string `json:"title_en"`
	AuditTeamTitleAr string `json:"title_ar"`
	Count            uint   `json:"count"`
}

type AuditPlanByTeamLeadDTO struct {
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

type AuditPlanByYearDTO struct {
	TitleEn  string `json:"title_en"`
	TitleAr  string `json:"title_ar"`
	Count    uint   `json:"count"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

type AuditByCategoryDTO struct {
	AuditCategoryID      uint   `json:"id"`
	AuditCategoryTitleEn string `json:"title_en"`
	AuditCategoryTitleAr string `json:"title_ar"`
	Count                uint   `json:"count"`
}

type AuditByStatusDTO struct {
	AuditStatusID      uint   `json:"id"`
	AuditStatusTitleEn string `json:"title_en"`
	AuditStatusTitleAr string `json:"title_ar"`
	ColorCode          string `json:"color_code"`
	Count              uint   `json:"count"`
}

type AuditByOpenVsClosedDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	Type      string `json:"type"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type AuditByMonthDTO struct {
	TitleEn  string `json:"title_en"` // Month in English
	TitleAr  string `json:"title_ar"` // Month in Arabic
	Count    uint   `json:"count"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

type AuditByTeamLeadDTO struct {
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

type AuditFindingByCategoryDTO struct {
	AuditFindingCategoryID      uint   `json:"id"`
	AuditFindingCategoryTitleEn string `json:"title_en"`
	AuditFindingCategoryTitleAr string `json:"title_ar"`
	Count                       uint   `json:"count"`
}
type AuditFindingByStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type AuditFindingBySeverityLevelDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type AuditFindingByResponsibleUnitDTO struct {
	ResponsibleUnitID        uint   `json:"id"`
	ResponsibleUnitTitleEn   string `json:"title_en"`
	ResponsibleUnitTitleAr   string `json:"title_ar"`
	ResponsibleUnitColorCode string `json:"color_code"`
	Count                    uint   `json:"count"`
}

type AuditFindingByMonthDTO struct {
	TitleEn  string `json:"title_en"` // Month in English
	TitleAr  string `json:"title_ar"` // Month in Arabic
	Count    uint   `json:"count"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}
type AuditFindingByAuditorDTO struct {
	ID                   uint   `json:"id"`
	UserFirstNameEn      string `json:"user_first_name_en"`
	UserFirstNameAr      string `json:"user_first_name_ar"`
	UserLastNameEn       string `json:"user_last_name_en"`
	UserLastNameAr       string `json:"user_last_name_ar"`
	UserTitleEn          string `json:"user_title_en"`
	UserTitleAr          string `json:"user_title_ar"`
	UserInitial          string `json:"user_initial"`
	UserColor            string `json:"user_color"`
	UserEmail            string `json:"user_email"`
	UserProfileImagePath string `json:"user_profile_image_path"`
	Count                int    `json:"count"`
}

type AuditActionPlanByStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type AuditActionPlanBySeverityLevelDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type AuditActionPlanByResponsibleUnitDTO struct {
	ResponsibleUnitID      uint   `json:"id"`
	ResponsibleUnitTitleEn string `json:"title_en"`
	ResponsibleUnitTitleAr string `json:"title_ar"`
	ColorCode              string `json:"color_code"`
	Count                  uint   `json:"count"`
}

type AuditActionPlanByResponsibleUserDTO struct {
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

type AuditActionPlanOnTimeVsDelayedDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}
