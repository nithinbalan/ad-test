package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type BcpCommunicationPlanListDTO struct {
	ID                              uint    `json:"id"`
	Title                           string  `json:"title"`
	BcpPhaseTitleEn                 string  `json:"bcp_phase_title_en"`
	BcpPhaseTitleAr                 string  `json:"bcp_phase_title_ar"`
	BcpProcessTitleEn               string  `json:"bcp_process_title_en"`
	BcpProcessTitleAr               string  `json:"bcp_process_title_ar"`
	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	Frequency                       *string `json:"frequency"`
	Recipient                       *string `json:"recipient"`
	Channel                         *string `json:"channel"`
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
	UpdatedUserFirstNameEn          *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              *string `json:"updated_user_initial"`
	UpdatedUserColor                *string `json:"updated_user_color"`
	UpdatedUserEmail                *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath     *string `json:"updated_user_profile_image_path"`
}

type BcpCommunicationPlanListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type BcpCommunicationPlanListResponseDTO struct {
	Items      []*BcpCommunicationPlanListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO   `json:"pagination"`
}

type BcpCommunicationPlanFormDTO struct {
	BcpPhaseID        uint    `json:"bcp_phase_id" binding:"required"`
	BcpProcessID      uint    `json:"bcp_process_id" binding:"required"`
	Title             string  `json:"title" binding:"required"`
	ResponsibleUserID uint    `json:"responsible_user_id" binding:"required"`
	Frequency         *string `json:"frequency"`
	Channel           *string `json:"channel"`
	Recipient         *string `json:"recipient"`
}

type BcpCommunicationPlanDetailDTO struct {
	ID                     int                              `json:"id"`
	Title                  string                           `json:"title"`
	BusinessContinuityPlan BusinessContinuityPlanDetailDTO  `json:"business_continuity_plan"`
	BcpPhaseID             uint                             `json:"bcp_phase_id"`
	BcpPhase               BcpPhaseDTO                      `json:"bcp_phase"`
	BcpProcessID           uint                             `json:"bcp_process_id"`
	BcpProcess             BcpProcessDTO                    `json:"bcp_process"`
	ResponsibleUserID      uint                             `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ResponsibleUser        organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
	Frequency              *string                          `json:"frequency"`
	Recipient              *string                          `json:"recipient"`
	Channel                *string                          `json:"channel"`
	CreatedAt              time.Time                        `json:"created_at"`
	UpdatedAt              *time.Time                       `json:"updated_at"`
	CreatedByUser          organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}
