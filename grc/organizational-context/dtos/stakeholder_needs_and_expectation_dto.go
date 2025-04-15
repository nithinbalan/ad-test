package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type StakeholderNeedsAndExpectationListDTO struct {
	ID                          int     `json:"id"`
	StakeholderID               uint    `json:"stakeholder_id"`
	StakeholderTitleEn          string  `json:"stakeholder_title_en"`
	StakeholderTitleAr          string  `json:"stakeholder_title_ar"`
	Title                       string  `json:"title"`
	SeverityLevelTitleEn        *string `json:"severity_level_title_en"`
	SeverityLevelTitleAr        *string `json:"severity_level_title_ar"`
	SeverityLevelColorCode      *string `json:"severity_level_color_code"`
	SeverityLevelType           *string `json:"severity_level_type"`
	SuccessCriteria             *string `json:"success_criteria"`
	Challenges                  *string `json:"challenges"`
	Description                 *string `json:"description"`
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

type StakeholderNeedsAndExpectationFormDTO struct {
	Title           string  `json:"title" binding:"required"`
	SeverityLevelID *uint   `json:"severity_level_id" binding:"omitempty"`
	SuccessCriteria *string `json:"success_criteria"`
	Challenges      *string `json:"challenges"`
	Description     *string `json:"description"`
}

type StakeholderNeedsAndExpectationListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type StakeholderNeedsAndExpectationDetailDTO struct {
	ID              uint                                           `json:"id"`
	Title           string                                         `json:"title"`
	SeverityLevelID *uint                                          `json:"severity_level_id"`
	SuccessCriteria *string                                        `json:"success_criteria"`
	Challenges      *string                                        `json:"challenges"`
	Description     *string                                        `json:"description"`
	Stakeholder     StakeholderDetailDTO                           `json:"stakeholder"`
	SeverityLevel   *masterDTOs.SeverityLevelDTO                   `json:"severity_level"`
	CreatedAt       time.Time                                      `json:"created_at"`
	UpdatedAt       *time.Time                                     `json:"updated_at"`
	CreatedByUser   organizationDTOs.UserMinimalDTO                `json:"created_by"`
	UpdatedByUser   *organizationDTOs.UserMinimalDTO               `json:"updated_by"`
	Evaluations     []*StakeholderNeedsAndExpectationEvaluationDTO `json:"evaluations"`
}

type StakeholderNeedsAndExpectationEvaluationDTO struct {
	ID               uint                             `json:"id"`
	EvaluationStatus EvaluationStatusDTO              `json:"evaluatiom_status"`
	Comment          string                           `json:"comment"`
	Score            *int                             `json:"score"`
	EvaluatedByID    uint                             `json:"evaluated_by_id"`
	EvaluatedAt      *string                          `json:"evaluated_at" binding:"omitempty,datetime=2006-01-02"`
	CreatedAt        time.Time                        `json:"created_at"`
	UpdatedAt        *time.Time                       `json:"updated_at"`
	CreatedByUser    organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser    *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}
