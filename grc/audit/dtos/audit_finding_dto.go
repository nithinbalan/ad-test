package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

// DTO for listing audit findings
type AuditFindingListDTO struct {
	ID                          uint    `json:"id"`
	AuditPlanID                 uint    `json:"audit_plan_id"`
	StartDate                   string  `json:"start_date"`
	EndDate                     string  `json:"end_date"`
	AuditPlanTitle              string  `json:"audit_plan_title"`
	AuditCategoryTitleEn        string  `json:"audit_category_title_en"`
	AuditCategoryTitleAr        string  `json:"audit_category_title_ar"`
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

// DTO for creating or updating audit findings
type AuditFindingFormDTO struct {
	AuditPlanID uint   `json:"audit_plan_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
}

// DTO for detailed audit finding view
type AuditFindingDetailDTO struct {
	ID              uint                             `json:"id"`
	AuditPlanID     uint                             `json:"audit_plan_id"`
	AuditPlan       AuditPlanDetailDTO               `json:"audit_plan"`
	StartDateString string                           `json:"start_date"`
	EndDateString   string                           `json:"end_date"`
	CreatedAt       time.Time                        `json:"created_at"`
	UpdatedAt       *time.Time                       `json:"updated_at"`
	CreatedByUser   organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser   *organizationDTOs.UserMinimalDTO `json:"updated_by"`
}

// DTO for query parameters for listing audit findings
type AuditFindingListQueryParamsDTO struct {
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

// DTO for response containing a list of audit findings and pagination info
type AuditFindingListResponseDTO struct {
	Items      []*AuditFindingListDTO       `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

// Struct for parsed date information
type AuditFindingParsedDates struct {
	StartDate time.Time
	EndDate   time.Time
}
