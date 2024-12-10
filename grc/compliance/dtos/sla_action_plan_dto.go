package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type SLAActionPlanListDTO struct {
	ID                              uint    `json:"id"`
	SLACheckQuestionTitleEn         string  `json:"sla_check_question_title_en"`
	SLACheckQuestionTitleAr         string  `json:"sla_check_question_title_ar"`
	Title                           string  `json:"title"`
	Description                     string  `json:"description"`
	SeverityLevelTitleEn            string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string  `json:"severity_level_title_ar"`
	SeverityLevelColorCode          string  `json:"severity_level_color_code"`
	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	ReviewUserFirstNameEn           string  `json:"review_user_first_name_en"`
	ReviewUserID                    *uint   `json:"review_user_id"`
	ReviewUserFirstNameAr           *string `json:"review_user_first_name_ar"`
	ReviewUserLastNameEn            *string `json:"review_user_last_name_en"`
	ReviewUserLastNameAr            *string `json:"review_user_last_name_ar"`
	ReviewUserInitial               *string `json:"review_user_initial"`
	ReviewUserColor                 *string `json:"review_user_color"`
	ReviewUserEmail                 *string `json:"review_user_email"`
	ReviewUserProfileImagePath      *string `json:"review_user_profile_image_path"`
	TargetDate                      string  `json:"target_date"`
	Percentage                      uint    `json:"percentage"`
	Comments                        string  `json:"comments"`
	ActionPlanStatusTitleEn         string  `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr         string  `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode       string  `json:"action_plan_status_color_code"`
	CreatedUserID                   uint    `json:"created_user_id"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string  `json:"created_user_last_name_ar"`
	CreatedUserInitial              string  `json:"created_user_initial"`
	CreatedUserColor                string  `json:"created_user_color"`
	CreatedUserEmail                string  `json:"created_user_email"`
	CreatedUserProfileImagePath     string  `json:"created_user_profile_image_path"`
	UpdatedUserID                   uint    `json:"updated_user_id"`
	UpdatedUserFirstNameEn          string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string  `json:"updated_user_initial"`
	UpdatedUserColor                string  `json:"updated_user_color"`
	UpdatedUserEmail                string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string  `json:"updated_user_profile_image_path"`
}

type SLAActionPlanListQueryParamsDTO struct {
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

type SLAActionPlanListResponseDTO struct {
	Items      []*SLAActionPlanListDTO      `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type SLAActionPlanFormDTO struct {
	SLACheckQuestionID *uint   `json:"sla_check_question_id" binding:"omitempty"`
	Title              string  `json:"title" binding:"required"`
	Description        *string `json:"description" binding:"omitempty"`
	SeverityLevelID    *uint   `json:"severity_level_id" binding:"omitempty"`
	ResponsibleUserID  uint    `json:"responsible_user_id" binding:"required"`
	TargetDate         *string `json:"target_date" binding:"omitempty,datetime=2006-01-02"`
	Percentage         *uint   `json:"percentage" binding:"omitempty"`
	Comments           *string `json:"comments" binding:"omitempty"`
	ReviewUserID       *uint   `json:"review_user_id"`
	ActionPlanStatusID uint    `json:"action_plan_status_id" binding:"required"`
}

type SLAActionPlanParsedDates struct {
	TargetDate *time.Time
}

type SLAActionPlanDetailDTO struct {
	ID                      int                              `json:"id"`
	SlaRegister             *SLAActionPlanDetailDTO          `json:"sla_register"`
	ReferenceNumber         string                           `json:"reference_number"`
	SLACheckQuestion        *SLACheckQuestionDetailDTO       `json:"sla_check_question"`
	Title                   string                           `json:"title"`
	Description             string                           `json:"description"`
	SeverityLevel           *masterDTOs.SeverityLevelDTO     `json:"severity_level"`
	ResponsibleUserID       uint                             `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ResponsibleUser         organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
	ReviewUser              organizationDTOs.UserMinimalDTO  `json:"review_user"`
	ActionPlanStatus        masterDTOs.ActionPlanStatusDTO   `json:"action_plan_status"`
	Percentage              *uint                            `json:"percentage"`
	Comments                *string                          `json:"comments"`
	TargetDateString        *string                          `json:"target_date"`
	ActionPlanStatusUpdates []*SLAActionPlanStatusUpdateDTO  `json:"action_plan_status_updates"`
	CreatedAt               time.Time                        `json:"created_at"`
	UpdatedAt               *time.Time                       `json:"updated_at"`
	CreatedByUser           organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser           *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}

type SLAActionPlanStatusUpdateDTO struct {
	ID               uint                                    `json:"id"`
	ActionPlanStatus masterDTOs.ActionPlanStatusDTO          `json:"action_plan_status"`
	Comment          string                                  `json:"comment"`
	Percentage       *uint                                   `json:"percentage"`
	CreatedAt        time.Time                               `json:"created_at"`
	UpdatedAt        *time.Time                              `json:"updated_at"`
	CreatedByUser    organizationDTOs.UserMinimalDTO         `json:"created_by_user"`
	UpdatedByUser    *organizationDTOs.UserMinimalDTO        `json:"updated_by_user"`
	Documents        []*SLAActionPlanStatusUpdateDocumentDTO `json:"documents"`
}

type SLAActionPlanStatusUpdateDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}

type SLAActionPlanStatusCountResponseDTO struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     int    `json:"count"`
}
