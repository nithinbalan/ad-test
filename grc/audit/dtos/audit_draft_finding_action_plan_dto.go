package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditDraftFindingActionPlanListDTO struct {
	ID                              uint   `json:"id"`
	AuditDraftFindingID             uint   `json:"audit_draft_finding_id"`
	AuditDraftFindingTitle          string `json:"audit_draft_finding_title"`
	AuditDraftFindingDescription    string `json:"audit_draft_finding_description"`
	Title                           string `json:"title"`
	Description                     string `json:"description"`
	ResponsibleUserID               uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	TargetDate                      string `json:"target_date"`
	Percentage                      uint   `json:"percentage"`
	Comments                        string `json:"comments"`
	CreatedUserID                   uint   `json:"created_user_id"`
	CreatedUserFirstNameEn          string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string `json:"created_user_last_name_ar"`
	CreatedUserInitial              string `json:"created_user_initial"`
	CreatedUserColor                string `json:"created_user_color"`
	CreatedUserEmail                string `json:"created_user_email"`
	CreatedUserProfileImagePath     string `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint  `json:"updated_user_id"`
	UpdatedUserFirstNameEn          string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string `json:"updated_user_initial"`
	UpdatedUserColor                string `json:"updated_user_color"`
	UpdatedUserEmail                string `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string `json:"updated_user_profile_image_path"`
}

type AuditDraftFindingActionPlanListQueryParamsDTO struct {
	Sort           string
	SortBy         string
	Keyword        string
	Status         string
	Page           int
	PageSize       int
	Filters        map[string][]string
	QuickFilter    []string
	TargetDateFrom string
	TargetDateTo   string
}

type AuditDraftFindingActionPlanListResponseDTO struct {
	Items      []*AuditDraftFindingActionPlanListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO          `json:"pagination"`
}

type AuditDraftFindingActionPlanFormDTO struct {
	Title             string  `json:"title" binding:"required"`
	Description       *string `json:"description" binding:"omitempty"`
	ResponsibleUserID uint    `json:"responsible_user_id" binding:"required"`
	TargetDate        *string `json:"target_date" binding:"omitempty,datetime=2006-01-02"`
}

type AuditDraftFindingActionPlanParsedDates struct {
	TargetDate *time.Time
}

type AuditDraftFindingActionPlanDetailDTO struct {
	ID                  int                             `json:"id"`
	AuditDraftFindingID uint                            `json:"audit_draft_finding_id"`
	Title               string                          `json:"title"`
	Description         string                          `json:"description"`
	AuditDraftFinding   AuditDraftFindingDetailDTO      `json:"audit_draft_finding"`
	ResponsibleUser     organizationDTOs.UserMinimalDTO `json:"responsible_user"`
	TargetDateString    *string                         `json:"target_date"`
	CreatedAt           time.Time                       `json:"created_at"`
	UpdatedAt           *time.Time                      `json:"updated_at"`
	CreatedByUser       organizationDTOs.UserMinimalDTO `json:"created_by_user"`
	UpdatedByUser       organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}
