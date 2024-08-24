package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditPlanListDTO struct {
	ID                                uint    `json:"id"`
	ReferenceNumber                   string  `json:"reference_number"`
	AuditCategoryID                   uint    `json:"audit_category_id"`
	AuditCategoryTitleEn              string  `json:"audit_category_title_en"`
	AuditCategoryTitleAr              string  `json:"audit_category_title_ar"`
	AuditTeamID                       *uint   `json:"audit_team_id"`
	TeamLeadID                        *uint   `json:"team_lead_id"`
	Title                             string  `json:"title"`
	StartDate                         *string `json:"start_date"`
	TargetDate                        *string `json:"target_date"`
	AuditStatusID                     uint    `json:"audit_status_id"`
	AuditStatusTitleEn                string  `json:"audit_status_title_en"`
	AuditStatusTitleAr                string  `json:"audit_status_title_ar"`
	AuditStatusType                   string  `json:"audit_status_type"`
	AuditStatusColorCode              string  `json:"audit_status_color_code"`
	TeamLeadFirstNameEn               string  `json:"team_lead_first_name_en"`
	TeamLeadFirstNameAr               string  `json:"team_lead_first_name_ar"`
	TeamLeadLastNameEn                string  `json:"team_lead_last_name_en"`
	TeamLeadLastNameAr                string  `json:"team_lead_last_name_ar"`
	TeamLeadInitial                   string  `json:"team_lead_initial"`
	TeamLeadColor                     string  `json:"team_lead_color"`
	TeamLeadEmail                     string  `json:"team_lead_email"`
	TeamLeadProfileImagePath          string  `json:"team_lead_profile_image_path"`
	AuditTeamTitleEn                  string  `json:"audit_team_title_en"`
	AuditTeamTitleAr                  string  `json:"audit_team_title_ar"`
	AuditTeamTeamLeadID               uint    `json:"audit_team_team_lead_id"`
	AuditTeamTeamLeadFirstNameEn      string  `json:"audit_team_team_lead_first_name_en"`
	AuditTeamTeamLeadFirstNameAr      string  `json:"audit_team_team_lead_first_name_ar"`
	AuditTeamTeamLeadLastNameEn       string  `json:"audit_team_team_lead_last_name_en"`
	AuditTeamTeamLeadLastNameAr       string  `json:"audit_team_team_lead_last_name_ar"`
	AuditTeamTeamLeadInitial          string  `json:"audit_team_team_lead_initial"`
	AuditTeamTeamLeadEmail            string  `json:"audit_team_team_lead_email"`
	AuditTeamTeamLeadColor            string  `json:"audit_team_team_lead_color"`
	AuditTeamTeamLeadProfileImagePath string  `json:"audit_team_team_lead_profile_image_path"`
	CreatedUserID                     uint    `json:"created_user_id"`
	CreatedUserFirstNameEn            string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr            string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn             string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr             string  `json:"created_user_last_name_ar"`
	CreatedUserInitial                string  `json:"created_user_initial"`
	CreatedUserColor                  string  `json:"created_user_color"`
	CreatedUserEmail                  string  `json:"created_user_email"`
	CreatedUserProfileImagePath       string  `json:"created_user_profile_image_path"`
	UpdatedUserID                     *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn            string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr            string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn             string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr             string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                string  `json:"updated_user_initial"`
	UpdatedUserColor                  string  `json:"updated_user_color"`
	UpdatedUserEmail                  string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath       string  `json:"updated_user_profile_image_path"`
	CreatedAt                         string  `json:"created_at"`
	UpdatedAt                         *string `json:"updated_at"`
	CreatedByUserID                   uint    `json:"created_by_id"`
	UpdatedByUserID                   *uint   `json:"updated_by_id"`
}

type AuditPlanFormDTO struct {
	ReferenceNumber  *string `json:"reference_number"`
	AuditCategoryID  uint    `json:"audit_category_id" binding:"required"`
	AuditTeamID      *uint   `json:"audit_team_id"`
	TeamLeadID       *uint   `json:"team_lead_id"`
	Title            string  `json:"title" binding:"required,max=255"`
	Description      *string `json:"description"`
	StartDate        *string `json:"start_date"`
	TargetDate       *string `json:"target_date"`
	Purpose          *string `json:"purpose"`
	Objective        *string `json:"objective"`
	UserIds          []*uint `json:"user_ids"`
	AuditableItemIds []*uint `json:"auditable_item_ids"`
}

type AuditPlanDetailDTO struct {
	ID               uint                             `json:"id"`
	ReferenceNumber  string                           `json:"reference_number"`
	AuditCategory    AuditCategoryDTO                 `json:"audit_category"`
	AuditTeam        *AuditTeamDetailDTO              `json:"audit_team"`
	TeamLead         *organizationDTOs.UserMinimalDTO `json:"team_lead"`
	Title            string                           `json:"title"`
	Description      *string                          `json:"description"`
	StartDateString  string                           `json:"start_date"`
	TargetDateString string                           `json:"target_date"`
	Purpose          *string                          `json:"purpose"`
	Objective        *string                          `json:"objective"`
	AuditStatus      AuditStatusDTO                   `json:"audit_status"`
	CreatedAt        time.Time                        `json:"created_at"`
	UpdatedAt        *time.Time                       `json:"updated_at"`
	CreatedByUser    organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser    *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	AuditableItems   []*AuditPlanAuditableItemDTO     `json:"auditable_items"`
	Members          []*AuditPlanAuditMemberDTO       `json:"members"`
}

type AuditPlanListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	StartDate   string
	TargetDate  string
}

type AuditPlanListResponseDTO struct {
	Items      []*AuditPlanListDTO          `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AuditPlanParsedDates struct {
	StartDate  *time.Time
	TargetDate *time.Time
}

type AuditPlanAuditMemberDTO struct {
	ID          *uint                           `json:"id"`
	AuditPlanID uint                            `json:"audit_plan_id"`
	UserID      uint                            `json:"user_id"`
	User        organizationDTOs.UserMinimalDTO `json:"user"`
}

type AuditPlanAuditableItemDTO struct {
	ID              *uint                  `json:"id"`
	AuditPlanID     uint                   `json:"audit_plan_id"`
	AuditableItemID uint                   `json:"auditable_item_id"`
	AuditableItem   AuditableItemDetailDTO `json:"auditable_item"`
}
