package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditListDTO struct {
	ID                          uint    `json:"id" gorm:"column:id"`
	AuditPlanID                 uint    `json:"audit_plan_id" gorm:"column:audit_plan_id"`
	StartDate                   string  `json:"start_date" gorm:"column:start_date"`
	EndDate                     string  `json:"end_date" gorm:"column:end_date"`
	AuditCategoryTitleEn        string  `json:"audit_category_title_en"`
	AuditCategoryTitleAr        string  `json:"audit_category_title_ar"`
	AuditStatusTitleEn          string  `json:"audit_status_title_en"`
	AuditStatusTitleAr          string  `json:"audit_status_title_ar"`
	AuditStatusType             string  `json:"audit_status_type"`
	AuditStatusColorCode        string  `json:"audit_status_color_code"`
	AuditPlanTitle              string  `json:"audit_plan_title"`
	AuditPlanDescription        string  `json:"audit_plan_description"`
	AuditPlanPurpose            string  `json:"audit_plan_purpose"`
	AuditPlanObjective          string  `json:"audit_plan_objective"`
	TeamLeadID                  *uint   `json:"team_lead_id"`
	TeamLeadFirstNameEn         string  `json:"team_lead_first_name_en"`
	TeamLeadFirstNameAr         string  `json:"team_lead_first_name_ar"`
	TeamLeadLastNameEn          string  `json:"team_lead_last_name_en"`
	TeamLeadLastNameAr          string  `json:"team_lead_last_name_ar"`
	TeamLeadEmail               string  `json:"team_lead_email"`
	TeamLeadInitial             string  `json:"team_lead_initial"`
	TeamLeadColor               string  `json:"team_lead_color"`
	TeamLeadProfileImagePath    string  `json:"team_lead_profile_image_path"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   *string `json:"updated_at"`
	CreatedUserID               uint    `json:"created_user_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
}

type AuditFormDTO struct {
	AuditPlanID uint   `json:"audit_plan_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
}

type AuditDetailDTO struct {
	ID                 uint                             `json:"id"`
	AuditPlanID        uint                             `json:"audit_plan_id"`
	AuditPlan          AuditPlanDetailDTO               `json:"audit_plan"`
	StartDateString    string                           `json:"start_date"`
	EndDateString      string                           `json:"end_date"`
	AuditChecklists    []*AuditChecklistDetailDTO       `json:"audit_checklists"`
	AuditDraftFindings []*AuditDraftFindingDetailDTO    `json:"audit_draft_findings"`
	AuditStatus        AuditStatusDTO                   `json:"audit_status"`
	CreatedAt          time.Time                        `json:"created_at"`
	UpdatedAt          *time.Time                       `json:"updated_at"`
	CreatedByUser      organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser      *organizationDTOs.UserMinimalDTO `json:"updated_by"`
}

type AuditListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	StartDate   string
	EndDate     string
}

type AuditListResponseDTO struct {
	Items      []*AuditListDTO              `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AuditParsedDates struct {
	StartDate time.Time
	EndDate   time.Time
}

type AuditAuditableItemListDTO struct {
	ID                              uint   `json:"id"`
	AuditPlanID                     uint   `json:"audit_plan_id"`
	AuditableItemID                 uint   `json:"auditable_item_id"`
	CreatedByID                     uint   `json:"created_by_id"`
	CreatedAt                       string `json:"created_at"`
	ProcessID                       *uint  `json:"process_id"`
	ControlID                       *uint  `json:"control_id"`
	DocumentID                      *uint  `json:"document_id"`
	IncidentRegisterID              *uint  `json:"incident_register_id"`
	ResponsibleUnitID               uint   `json:"responsible_unit_id"`
	ComplianceRequirementID         *uint  `json:"compliance_requirement_id"`
	RiskRegisterID                  *uint  `json:"risk_register_id"`
	ProcessTitleEn                  string `json:"process_title_en"`
	ProcessTitleAr                  string `json:"process_title_ar"`
	ControlTitleEn                  string `json:"control_title_en"`
	ControlTitleAr                  string `json:"control_title_ar"`
	DocumentTitleEn                 string `json:"document_title_en"`
	DocumentTitleAr                 string `json:"document_title_ar"`
	IncidentRegisterTitle           string `json:"incident_register_title"`
	ComplianceRequirementTitle      string `json:"compliance_requirement_title"`
	ResponsibleUnitTitleEn          string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	RiskRegisterTitle               string `json:"risk_register_title"`
	UserID                          uint   `json:"user_id"`
	UserFirstNameEn                 string `json:"user_first_name_en"`
	UserFirstNameAr                 string `json:"user_first_name_ar"`
	UserLastNameEn                  string `json:"user_last_name_en"`
	UserLastNameAr                  string `json:"user_last_name_ar"`
	UserEmail                       string `json:"user_email"`
	UserInitial                     string `json:"user_initial"`
	UserColor                       string `json:"user_color"`
	UserProfileImagePath            string `json:"user_profile_image_path"`
	ReferenceNumber                 string `json:"reference_number"`
	AuditPlanTitle                  string `json:"audit_plan_title"`
	AuditableItemTitleEn            string `json:"auditable_item_title_en"`
	AuditableItemTitleAr            string `json:"auditable_item_title_ar"`
	AuditableItemType               string `json:"auditable_item_type"`
	AuditableItemComment            string `json:"auditable_item_comment"`
	AuditCategoryTitleEn            string `json:"audit_category_title_en"`
	AuditCategoryTitleAr            string `json:"audit_category_title_ar"`
	AuditTeamTitleEn                string `json:"audit_team_title_en"`
	AuditTeamTitleAr                string `json:"audit_team_title_ar"`
	TeamLeadID                      uint   `json:"team_lead_id"`
	TeamLeadFirstNameEn             string `json:"team_lead_first_name_en"`
	TeamLeadFirstNameAr             string `json:"team_lead_first_name_ar"`
	TeamLeadLastNameEn              string `json:"team_lead_last_name_en"`
	TeamLeadLastNameAr              string `json:"team_lead_last_name_ar"`
	TeamLeadInitial                 string `json:"team_lead_initial"`
	TeamLeadColor                   string `json:"team_lead_color"`
	TeamLeadEmail                   string `json:"team_lead_email"`
	TeamLeadProfileImagePath        string `json:"team_lead_profile_image_path"`
	AuditStatusTitleEn              string `json:"audit_status_title_en"`
	AuditStatusTitleAr              string `json:"audit_status_title_ar"`
	AuditStatusType                 string `json:"audit_status_type"`
	AuditStatusColorCode            string `json:"audit_status_color_code"`
	CreatedUserID                   uint   `json:"created_user_id"`
	CreatedUserFirstNameEn          string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string `json:"created_user_last_name_ar"`
	CreatedUserInitial              string `json:"created_user_initial"`
	CreatedUserColor                string `json:"created_user_color"`
	CreatedUserEmail                string `json:"created_user_email"`
	CreatedUserProfileImagePath     string `json:"created_user_profile_image_path"`
	AuditorRiskRatingTitleEn        string `json:"auditor_risk_rating_title_en"`
	AuditorRiskRatingTitleAr        string `json:"auditor_risk_rating_title_ar"`
	AuditorRiskRatingColorCode      string `json:"auditor_risk_rating_color_code"`
	AuditorRiskRatingType           string `json:"auditor_risk_rating_type"`
}

type AuditAuditableItemListResponseDTO struct {
	Items      []*AuditAuditableItemListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AuditAuditableItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}
