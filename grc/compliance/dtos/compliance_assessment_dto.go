package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type ComplianceAssessmentListDTO struct {
	ID                          uint       `json:"id"`
	ComplianceRegisterID        uint       `json:"compliance_register_id"`
	ComplianceRegisterTitleEn   string     `json:"compliance_register_title_en"`
	ComplianceRegisterTitleAr   string     `json:"compliance_register_title_ar"`
	Description                 *string    `json:"description"`
	OrganizationUnitID          *uint      `json:"organization_unit_id"`
	OrganizationUnitTitleEn     string     `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr     string     `json:"organization_unit_title_ar"`
	ReviewUserID                uint       `json:"review_user_id"`
	ReviewUserFirstNameEn       string     `json:"review_user_first_name_en"`
	ReviewUserFirstNameAr       string     `json:"review_user_first_name_ar"`
	ReviewUserLastNameEn        string     `json:"review_user_last_name_en"`
	ReviewUserLastNameAr        string     `json:"review_user_last_name_ar"`
	ReviewUserInitial           string     `json:"review_user_initial"`
	ReviewUserColor             string     `json:"review_user_color"`
	ReviewUserEmail             string     `json:"review_user_email"`
	ReviewUserProfileImagePath  string     `json:"review_user_profile_image_path"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   *time.Time `json:"updated_at"`
	CreatedUserID               uint       `json:"created_user_id"`
	CreatedUserFirstNameEn      string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string     `json:"created_user_last_name_ar"`
	CreatedUserInitial          string     `json:"created_user_initial"`
	CreatedUserColor            string     `json:"created_user_color"`
	CreatedUserEmail            string     `json:"created_user_email"`
	CreatedUserProfileImagePath string     `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string    `json:"updated_user_initial"`
	UpdatedUserColor            *string    `json:"updated_user_color"`
	UpdatedUserEmail            *string    `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string    `json:"updated_user_profile_image_path"`
}

type ComplianceAssessmentListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type ComplianceAssessmentListResponseDTO struct {
	Items      []*ComplianceAssessmentListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO   `json:"pagination"`
}

type ComplianceAssessmentFormDTO struct {
	ComplianceRegisterID uint    `json:"compliance_register_id" binding:"required"`
	OrganizationUnitIDs  []*uint `json:"organization_unit_ids" binding:"required"`
	ReviewUserID         uint    `json:"review_user_id" binding:"required"`
	Description          *string `json:"description"`
}

type ComplianceAssessmentUpdateFormDTO struct {
	ComplianceRegisterID uint    `json:"compliance_register_id" binding:"required"`
	OrganizationUnitIDs  uint    `json:"organization_unit_ids" binding:"required"`
	ReviewUserID         uint    `json:"review_user_id" binding:"required"`
	Description          *string `json:"description"`
}

type ComplianceAssessmentDetailDTO struct {
	ID                     int                                   `json:"id"`
	ComplianceRegister     *ComplianceRegisterDetailDTO          `json:"compliance_register"`
	OrganizationUnitID     uint                                  `json:"organization_unit_id"`
	ReviewuserID           uint                                  `json:"review_user_id"`
	OrganizationUnit       settingsDTOs.OrganizationUnitDTO      `json:"organization_unit"`
	ReviewUser             *organizationDTOs.UserMinimalDTO      `json:"review_user"`
	Description            *string                               `json:"description"`
	CreatedAt              time.Time                             `json:"created_at"`
	UpdatedAt              *time.Time                            `json:"updated_at"`
	CreatedByUser          organizationDTOs.UserMinimalDTO       `json:"created_by_user"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO      `json:"updated_by_user"`
	AssessmentRequirements []*ComplianceAssessmentRequirementDTO `json:"compliance_requirements"`
}

type ComplianceAssessmentRequirementDTO struct {
	ID                      *uint                                              `json:"id"`
	ComplianceAssessmentID  uint                                               `json:"compliance_assessment_id"`
	ComplianceRequirementID uint                                               `json:"compliance_requirement_id"`
	ComplianceRequirement   ComplianceRequirementDetailDTO                     `json:"compliance_requirement"`
	IsCompleted             bool                                               `json:"is_completed"`
	Score                   *int                                               `json:"score"`
	Comments                *string                                            `json:"comments"`
	ComplianceStatusID      *uint                                              `json:"compliance_status_id"`
	ComplianceStatus        *ComplianceStatusDTO                               `json:"compliance_status"`
	RequiredDocumentTitles  []*ComplianceAssessmentRequirementDocumentTitleDTO `json:"required_document_titles"`
}

type AssessmentCompletionSummaryDTO struct {
	Total        int `json:"total"`
	Completed    int `json:"completed"`
	NotCompleted int `json:"not_completed"`
}

type ComplianceStatusSummaryDTO struct {
	TotalRequirements int            `json:"total_requirements"`
	StatusCounts      map[string]int `json:"status_counts"`
	PendingCount      int            `json:"pending_count"`
}

type ComplianceScoreSummaryDTO struct {
	TotalRequirements int     `json:"total_requirements_assessed"`
	AverageScore      float64 `json:"average_score"`
}
