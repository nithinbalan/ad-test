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
	TitleEn string `json:"title_en" gorm:"column:title_en"` // Month in English
	TitleAr string `json:"title_ar" gorm:"column:title_ar"` // Month in Arabic
	Count   uint   `json:"count" gorm:"column:count"`
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
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserTitleEn          string `json:"responsible_user_title_en"`
	ResponsibleUserTitleAr          string `json:"responsible_user_title_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	Count                           int    `json:"count"`
}

type AuditActionPlanOnTimeVsDelayedDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}
