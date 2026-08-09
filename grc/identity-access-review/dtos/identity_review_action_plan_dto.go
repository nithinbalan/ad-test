package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
)

type IdentityReviewActionPlanFormDTO struct {
	Title             string  `json:"title" binding:"required,max=500"`
	Description       *string `json:"description" binding:"omitempty"`
	SeverityLevelID   *uint   `json:"severity_level_id" binding:"omitempty"`
	TargetDate        *string `json:"target_date" binding:"omitempty,datetime=2006-01-02"`
	ResponsibleUserID uint    `json:"responsible_user_id" binding:"required"`
	ReviewUserID      *uint   `json:"review_user_id" binding:"omitempty"`
	Percentage        uint    `json:"percentage" binding:"min=0,max=100"`
	Comments          *string `json:"comments" binding:"omitempty"`

	AccountDecisionID uint  `json:"account_decision_id" binding:"required"`
	PrivilegedID      *uint `json:"privileged_id" binding:"omitempty"`
}

type IdentityReviewActionPlanListQueryParamsDTO struct {
	Page           int
	PageSize       int
	Keyword        string
	Sort           string
	SortBy         string
	Filters        map[string][]string
	QuickFilter    []string
	TargetDateFrom string
	TargetDateTo   string
	DateFrom       string
	DateTo         string
	Export         bool
}

type IdentityReviewActionPlanListResponseDTO struct {
	Items      []*IdentityReviewActionPlanListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO       `json:"pagination"`
}

type IdentityReviewActionPlanListDTO struct {
	ID                              uint       `json:"id"`
	IdentityReviewUserID            uint       `json:"identity_review_user_id"`
	IdentityReviewUsername          string     `json:"identity_review_username"`
	UserIdentifier                  *string    `json:"user_identifier"`
	Title                           string     `json:"title"`
	Description                     *string    `json:"description"`
	SeverityLevelID                 *uint      `json:"severity_level_id"`
	SeverityLevelTitleEn            *string    `json:"severity_level_title_en"`
	SeverityLevelTitleAr            *string    `json:"severity_level_title_ar"`
	SeverityLevelColorCode          *string    `json:"severity_level_color_code"`
	ResponsibleUserID               uint       `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string     `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string     `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string     `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string     `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string     `json:"responsible_user_initial"`
	ResponsibleUserColor            string     `json:"responsible_user_color"`
	ResponsibleUserEmail            string     `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath *string    `json:"responsible_user_profile_image_path"`
	ReviewUserID                    *uint      `json:"review_user_id"`
	ReviewUserFirstNameEn           *string    `json:"review_user_first_name_en"`
	ReviewUserFirstNameAr           *string    `json:"review_user_first_name_ar"`
	ReviewUserLastNameEn            *string    `json:"review_user_last_name_en"`
	ReviewUserLastNameAr            *string    `json:"review_user_last_name_ar"`
	ReviewUserInitial               *string    `json:"review_user_initial"`
	ReviewUserColor                 *string    `json:"review_user_color"`
	ReviewUserEmail                 *string    `json:"review_user_email"`
	ReviewUserProfileImagePath      *string    `json:"review_user_profile_image_path"`
	TargetDate                      *string    `json:"target_date"`
	Percentage                      uint       `json:"percentage"`
	Comments                        *string    `json:"comments"`
	ActionPlanStatusID              uint       `json:"action_plan_status_id"`
	ActionPlanStatusTitleEn         string     `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr         string     `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode       string     `json:"action_plan_status_color_code"`
	ActionPlanStatusType            string     `json:"action_plan_status_type"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedAt                       *time.Time `json:"updated_at"`
}

type IdentityReviewActionPlanSummaryDTO struct {
	OpenCount   int `json:"open_count"`
	ClosedCount int `json:"closed_count"`
}

type IdentityReviewActionPlanParsedDates struct {
	TargetDate *time.Time
}
