package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type StakeholderEngagementListDTO struct {
	ID                            int     `json:"id"`
	StakeholderID                 uint    `json:"stakeholder_id"`
	StakeholderTitleEn            string  `json:"stakeholder_title_en"`
	StakeholderTitleAr            string  `json:"stakeholder_title_ar"`
	Title                         string  `json:"title"`
	ReviewFrequencyTitleEn        *string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr        *string `json:"review_frequency_title_ar"`
	ReviewFrequencyType           *string `json:"review_frequency_type"`
	EngagementStrategy            *string `json:"engagement_strategy"`
	Description                   *string `json:"description"`
	AssignedOwnerID               uint    `json:"assigned_owner_id"`
	AssignedOwnerFirstNameEn      string  `json:"assigned_owner_first_name_en"`
	AssignedOwnerFirstNameAr      string  `json:"assigned_owner_first_name_ar"`
	AssignedOwnerLastNameEn       string  `json:"assigned_owner_last_name_en"`
	AssignedOwnerLastNameAr       string  `json:"assigned_owner_last_name_ar"`
	AssignedOwnerInitial          string  `json:"assigned_owner_initial"`
	AssignedOwnerColor            string  `json:"assigned_owner_color"`
	AssignedOwnerEmail            string  `json:"assigned_owner_email"`
	AssignedOwnerProfileImagePath string  `json:"assigned_owner_profile_image_path"`
	CreatedUserID                 uint    `json:"created_user_id"`
	CreatedUserFirstNameEn        string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr        string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn         string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr         string  `json:"created_user_last_name_ar"`
	CreatedUserInitial            string  `json:"created_user_initial"`
	CreatedUserColor              string  `json:"created_user_color"`
	CreatedUserEmail              string  `json:"created_user_email"`
	CreatedUserProfileImagePath   string  `json:"created_user_profile_image_path"`
	UpdatedUserID                 *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn        *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr        *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn         *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr         *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial            *string `json:"updated_user_initial"`
	UpdatedUserColor              *string `json:"updated_user_color"`
	UpdatedUserEmail              *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath   *string `json:"updated_user_profile_image_path"`
}

type StakeholderEngagementFormDTO struct {
	Title              string  `json:"title" binding:"required"`
	ReviewFrequencyID  *uint   `json:"review_frequency_id" binding:"omitempty"`
	EngagementStrategy *string `json:"engagement_strategy"`
	Description        *string `json:"description"`
	AssignedOwnerID    *uint   `json:"assigned_owner_id"`
}

type StakeholderEngagementListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type StakeholderEngagementDetailDTO struct {
	ID                 uint                                  `json:"id"`
	Title              string                                `json:"title"`
	ReviewFrequencyID  *uint                                 `json:"review_frequency_id"`
	EngagementStrategy *string                               `json:"engagement_strategy"`
	Description        *string                               `json:"description"`
	Stakeholder        StakeholderDetailDTO                  `json:"stakeholder"`
	ReviewFrequency    *masterDTOs.ReviewFrequencyDTO        `json:"review_frequency"`
	AssignedOwnerID    *uint                                 `json:"assigned_owner_id"`
	AssignedOwner      *organizationDTOs.UserMinimalDTO      `json:"assigned_owner"`
	CreatedAt          time.Time                             `json:"created_at"`
	UpdatedAt          *time.Time                            `json:"updated_at"`
	CreatedByUser      organizationDTOs.UserMinimalDTO       `json:"created_by"`
	UpdatedByUser      *organizationDTOs.UserMinimalDTO      `json:"updated_by"`
	Evaluations        []*StakeholderEngagementEvaluationDTO `json:"evaluations"`
}

type StakeholderEngagementEvaluationDTO struct {
	ID              uint                             `json:"id"`
	Comment         string                           `json:"comment"`
	Score           *int                             `json:"score"`
	EvaluatedByID   uint                             `json:"evaluated_by_id"`
	EvaluatedByUser organizationDTOs.UserMinimalDTO  `json:"evaluated_by_user"`
	EvaluatedAt     *time.Time                       `json:"evaluated_at" binding:"omitempty,datetime=2006-01-02"`
	CreatedAt       time.Time                        `json:"created_at"`
	UpdatedAt       *time.Time                       `json:"updated_at"`
	CreatedByUser   organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser   *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}
