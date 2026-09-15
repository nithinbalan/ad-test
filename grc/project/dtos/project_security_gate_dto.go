package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type ProjectSecurityGateListDTO struct {
	ID                            uint       `json:"id"`
	ProjectTitleEn                string     `json:"project_title_en"`
	ProjectTitleAr                string     `json:"project_title_ar"`
	StageID                       uint       `json:"stage_id"`
	StageTitleEn                  string     `json:"stage_title_en"`
	StageTitleAr                  string     `json:"stage_title_ar"`
	StageType                     string     `json:"stage_type"`
	StageColorCode                string     `json:"stage_color_code"`
	BusinessOwnerID               uint       `json:"business_owner_id"`
	BusinessOwnerFirstNameEn      string     `json:"business_owner_first_name_en"`
	BusinessOwnerFirstNameAr      string     `json:"business_owner_first_name_ar"`
	BusinessOwnerLastNameEn       string     `json:"business_owner_last_name_en"`
	BusinessOwnerLastNameAr       string     `json:"business_owner_last_name_ar"`
	BusinessOwnerInitial          string     `json:"business_owner_initial"`
	BusinessOwnerColor            string     `json:"business_owner_color"`
	BusinessOwnerEmail            string     `json:"business_owner_email"`
	BusinessOwnerProfileImagePath string     `json:"business_owner_profile_image_path"`
	ReadinessStatusID             uint       `json:"readiness_status_id"`
	ReadinessStatusTitleEn        string     `json:"readiness_status_title_en"`
	ReadinessStatusTitleAr        string     `json:"readiness_status_title_ar"`
	ReadinessStatusType           string     `json:"readiness_status_type"`
	ReadinessStatusColorCode      string     `json:"readiness_status_color_code"`
	Notes                         *string    `json:"notes"`
	ChecklistCount                int64      `json:"checklist_count"`
	CreatedByID                   uint       `json:"created_by_id"`
	CreatedUserFirstNameEn        string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr        string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn         string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr         string     `json:"created_user_last_name_ar"`
	CreatedUserInitial            string     `json:"created_user_initial"`
	CreatedUserColor              string     `json:"created_user_color"`
	CreatedUserEmail              string     `json:"created_user_email"`
	CreatedUserProfileImagePath   string     `json:"created_user_profile_image_path"`
	CreatedAt                     time.Time  `json:"created_at"`
	UpdatedByID                   *uint      `json:"updated_by_id"`
	UpdatedUserFirstNameEn        *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr        *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn         *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr         *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial            *string    `json:"updated_user_initial"`
	UpdatedUserColor              *string    `json:"updated_user_color"`
	UpdatedUserEmail              *string    `json:"updated_user_email"`
	UpdatedUserProfileImagePath   *string    `json:"updated_user_profile_image_path"`
	UpdatedAt                     *time.Time `json:"updated_at"`
}

type ProjectSecurityGateDetailDTO struct {
	ID              uint                                    `json:"id"`
	ProjectTitleEn  string                                  `json:"project_title_en"`
	ProjectTitleAr  string                                  `json:"project_title_ar"`
	Stage           *ProjectSecurityGateStageDTO            `json:"stage"`
	TemplateVersion *ProjectSecurityGateTemplateVersionDTO  `json:"template_version"`
	BusinessOwner   *organizationDTOs.UserMinimalDTO        `json:"business_owner"`
	ReadinessStatus *ProjectSecurityGateReadinessStatusDTO  `json:"readiness_status"`
	Notes           *string                                 `json:"notes"`
	CreatedByUser   *organizationDTOs.UserMinimalDTO        `json:"created_by"`
	CreatedAt       time.Time                               `json:"created_at"`
	UpdatedByUser   *organizationDTOs.UserMinimalDTO        `json:"updated_by"`
	UpdatedAt       *time.Time                              `json:"updated_at"`
	Classifications []*ProjectSecurityGateClassificationDTO `json:"classifications"`
}

type ProjectSecurityGateTemplateVersionDTO struct {
	ID      uint   `json:"id"`
	Version int    `json:"version"`
	Status  string `json:"status"`
}

type ProjectSecurityGateListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Export   bool

	// Optional filters, zero means unfiltered
	ReadinessStatusID uint
	StageID           uint
	ClassificationID  uint
}

type ProjectSecurityGateListResponseDTO struct {
	Items      []*ProjectSecurityGateListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO  `json:"pagination"`
}

type ProjectSecurityGateSummaryDTO struct {
	Total  int64 `json:"total"`
	Draft  int64 `json:"draft"`
	GoLive int64 `json:"go_live"`
}

type ProjectSecurityGateFormDTO struct {
	ProjectTitleEn    string  `json:"project_title_en" binding:"required,max=100"`
	ProjectTitleAr    string  `json:"project_title_ar" binding:"required,max=100"`
	StageID           uint    `json:"stage_id" binding:"required"`
	BusinessOwnerID   uint    `json:"business_owner_id" binding:"required"`
	ClassificationIDs []uint  `json:"classification_ids" binding:"required,min=1"`
	Notes             *string `json:"notes"`
}
