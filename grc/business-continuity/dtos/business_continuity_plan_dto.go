package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type BusinessContinuityPlanListDTO struct {
	ID                              uint    `json:"id"`
	ReferenceNumber                 string  `json:"reference_number"`
	BusinessFunctionID              uint    `json:"business_function_id"`
	BusinessFunctionTitle           string  `json:"business_function_title"`
	Version                         *int    `json:"version"`
	IsLatestVersion                 *bool   `json:"is_latest_version"`
	Purpose                         *string `json:"purpose"`
	Scope                           *string `json:"scope"`
	Objectives                      *string `json:"objectives"`
	Assumptions                     *string `json:"assumptions"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	BcpStatusTitleEn                string  `json:"bcp_status_title_en"`
	BcpStatusTitleAr                string  `json:"bcp_status_title_ar"`
	BcpStatusColorCode              string  `json:"bcp_status_color_code"`
	BcpStatusType                   string  `json:"bcp_status_type"`
	Rto                             *int    `json:"rto"`
	RtoUnitId                       *uint   `json:"rto_unit_id"`
	RtoUnitTitleEn                  *string `json:"rto_unit_title_en"`
	RtoUnitTitleAr                  *string `json:"rto_unit_title_ar"`
	Mtpd                            *int    `json:"mtpd"`
	MtpdUnitId                      *uint   `json:"mtpd_unit_id"`
	MtpdUnitTitleEn                 *string `json:"mtpd_unit_title_en"`
	MtpdUnitTitleAr                 *string `json:"mtpd_unit_title_ar"`
	Rpo                             *int    `json:"rpo"`
	RpoUnitId                       *uint   `json:"rpo_unit_id"`
	RpoUnitTitleEn                  *string `json:"rpo_unit_title_en"`
	RpoUnitTitleAr                  *string `json:"rpo_unit_title_ar"`
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
	UpdatedUserFirstNameEn          *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              *string `json:"updated_user_initial"`
	UpdatedUserColor                *string `json:"updated_user_color"`
	UpdatedUserEmail                *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath     *string `json:"updated_user_profile_image_path"`
}

type BusinessContinuityPlanFormDTO struct {
	BusinessFunctionID uint                          `json:"business_function_id" binding:"required"`
	Version            *int                          `json:"version"`
	Purpose            *string                       `json:"purpose"`
	Scope              *string                       `json:"scope"`
	Objectives         *string                       `json:"objectives"`
	Assumptions        *string                       `json:"assumptions"`
	ResponsibleUserID  *uint                         `json:"responsible_user_id" binding:"required"`
	Documents          []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type BusinessContinuityPlanListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type BusinessContinuityPlanListResponseDTO struct {
	Items      []*BusinessContinuityPlanListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}

type BusinessContinuityPlanDetailDTO struct {
	ID                 uint                             `json:"id"`
	ReferenceNumber    string                           `json:"reference_number"`
	BusinessFunctionID uint                             `json:"business_function_id"`
	BusinessFunction   BusinessFunctionDetailDTO        `json:"business_function"`
	Version            int                              `json:"version"`
	IsLatestVersion    *bool                            `json:"is_latest_version"`
	Purpose            *string                          `json:"purpose"`
	Scope              *string                          `json:"scope"`
	Objectives         *string                          `json:"objectives"`
	Assumptions        *string                          `json:"assumptions"`
	ResponsibleUserID  *uint                            `json:"responsible_user_id"`
	ResponsibleUser    *organizationDTOs.UserMinimalDTO `json:"responsible_user"`
	ResponsibleUnit    settingsDTOs.OrganizationUnitDTO `json:"responsible_unit"`

	BcpStatus             BcpStatusDTO                     `json:"bcp_status"`
	Rto                   *int                             `json:"rto"`
	RtoUnitId             *uint                            `json:"rto_unit_id"`
	RtoUnit               *masterDTOs.TimeUnitDTO          `json:"rto_unit"`
	Mtpd                  *int                             `json:"mtpd"`
	MtpdUnitId            *uint                            `json:"mtpd_unit_id"`
	MtpdUnit              *masterDTOs.TimeUnitDTO          `json:"mtpd_unit"`
	Rpo                   *int                             `json:"rpo"`
	RpoUnitId             *uint                            `json:"rpo_unit_id"`
	RpoUnit               *masterDTOs.TimeUnitDTO          `json:"rpo_unit"`
	Processes             []*BcpProcessDTO                 `json:"processes"`
	Documents             []*BcpDocumentDTO                `json:"documents" gorm:"-"`
	BcpRecoveryPlans      []*BcpRecoveryPlanDetailDTO      `json:"bcp_recovery_plans"`
	BcpCommunicationPlans []*BcpCommunicationPlanDetailDTO `json:"bcp_communication_plans"`
	CreatedAt             time.Time                        `json:"created_at"`
	UpdatedAt             *time.Time                       `json:"updated_at"`
	CreatedByUser         organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser         *organizationDTOs.UserMinimalDTO `json:"updated_by"`
}

type BcpProcessDTO struct {
	ID                       *uint                             `json:"id"`
	BusinessContinuityPlanID uint                              `json:"business_continuity_plan_id"`
	ProcessID                uint                              `json:"process_id"`
	Rto                      *int                              `json:"rto"`
	RtoUnitId                *uint                             `json:"rto_unit_id"`
	RtoUnit                  *masterDTOs.TimeUnitDTO           `json:"rto_unit"`
	Mtpd                     *int                              `json:"mtpd"`
	MtpdUnitId               *uint                             `json:"mtpd_unit_id"`
	MtpdUnit                 *masterDTOs.TimeUnitDTO           `json:"mtpd_unit"`
	Rpo                      *int                              `json:"rpo"`
	RpoUnitId                *uint                             `json:"rpo_unit_id"`
	RpoUnit                  *masterDTOs.TimeUnitDTO           `json:"rpo_unit"`
	Process                  organizationDTOs.ProcessDetailDTO `json:"process"`
}

type BcpDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}

type BcpTestSummaryDTO struct {
	TestScore        int `json:"test_score"`
	ExecutedTests    int `json:"executed_tests"`
	SuccessfulTests  int `json:"successful_tests"`
	FailedTests      int `json:"failed_tests"`
	PartiallySuccess int `json:"partially_success"`
}
