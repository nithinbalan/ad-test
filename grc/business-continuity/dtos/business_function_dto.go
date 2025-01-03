package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type BusinessFunctionListDTO struct {
	ID                              uint    `json:"id"`
	ReferenceNumber                 string  `json:"reference_number"`
	Title                           string  `json:"title"`
	Description                     *string `json:"description"`
	SeverityLevelID                 *uint   `json:"severity_level_id"`
	ReviewFrequencyID               *uint   `json:"review_frequency_id"`
	LastReviewDate                  *string `json:"last_review_date"`
	NextReviewDate                  *string `json:"next_review_date"`
	BcpStatus                       bool    `json:"bcp_status"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	SeverityLevelTitleEn            string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string  `json:"severity_level_title_ar"`
	SeverityLevelColorCode          string  `json:"severity_level_color_code"`
	ReviewFrequencyTitleEn          *string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr          *string `json:"review_frequency_title_ar"`
	CreatedAt                       string  `json:"created_at"`
	UpdatedAt                       string  `json:"updated_at"`
	CreatedUserID                   uint    `json:"created_user_id"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string  `json:"created_user_last_name_ar"`
	CreatedUserInitial              string  `json:"created_user_initial"`
	CreatedUserColor                string  `json:"created_user_color"`
	CreatedUserEmail                string  `json:"created_user_email"`
	CreatedUserProfileImagePath     string  `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn          string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string  `json:"updated_user_initial"`
	UpdatedUserColor                string  `json:"updated_user_color"`
	UpdatedUserEmail                string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string  `json:"updated_user_profile_image_path"`
}

type BusinessFunctionFormDTO struct {
	Title             string  `json:"title" binding:"required,max=255"`
	Description       *string `json:"description"`
	SeverityLevelID   *uint   `json:"severity_level_id" binding:"required"`
	ResponsibleUserID uint    `json:"responsible_user_id" binding:"required"`
	ProcessIDs        []*uint `json:"process_ids"`
	ReviewFrequencyID *uint   `json:"review_frequency_id"`
	LastReviewDate    *string `json:"last_review_date"`
	NextReviewDate    *string `json:"next_review_date"`
}

type BusinessFunctionListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type BusinessFunctionListResponseDTO struct {
	Items      []*BusinessFunctionListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type BusinessFunctionDetailDTO struct {
	ID                   uint                            `json:"id"`
	ReferenceNumber      string                          `json:"reference_number"`
	Title                string                          `json:"title"`
	Description          string                          `json:"description"`
	SeverityLevelID      *uint                           `json:"severity_level_id" binding:"required"`
	ResponsibleUserID    uint                            `json:"responsible_user_id" binding:"required"`
	ReviewFrequencyID    uint                            `json:"review_frequency_id"`
	BcpStatus            bool                            `json:"bcp_status" gorm:"default:false"`
	NextReviewDateString *string                         `json:"next_review_date"`
	LastReviewDateString *string                         `json:"last_review_date"`
	ResponsibleUser      organizationDTOs.UserMinimalDTO `json:"responsible_user"`
	SeverityLevel        *masterDTOs.SeverityLevelDTO    `json:"severity_level"`
	ReviewFrequency      *masterDTOs.ReviewFrequencyDTO  `json:"review_frequency"`

	Processes     []*BusinessFunctionProcessDTO   `json:"processes"`
	CreatedAt     time.Time                       `json:"created_at"`
	UpdatedAt     *time.Time                      `json:"updated_at"`
	CreatedByUser organizationDTOs.UserMinimalDTO `json:"created_by"`
	UpdatedByUser organizationDTOs.UserMinimalDTO `json:"updated_by"`
}

type BusinessFunctionProcessDTO struct {
	ID                 *uint                             `json:"id"`
	BusinessFunctionID uint                              `json:"business_function_register_id"`
	ProcessID          uint                              `json:"process_id"`
	Process            organizationDTOs.ProcessDetailDTO `json:"process"`
}

type BusinessFunctionParsedDates struct {
	NextReviewDate *time.Time
	LastReviewDate *time.Time
}
