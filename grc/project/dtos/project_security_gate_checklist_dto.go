package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type ProjectSecurityGateChecklistListDTO struct {
	ID                                  uint       `json:"id"`
	SecurityGateID                      uint       `json:"security_gate_id"`
	RequirementID                       uint       `json:"requirement_id"`
	TitleEn                             string     `json:"title_en"`
	TitleAr                             string     `json:"title_ar"`
	IsCritical                          bool       `json:"is_critical"`
	ResponseStatusID                    *uint      `json:"response_status_id"`
	ResponseStatusTitleEn               *string    `json:"response_status_title_en"`
	ResponseStatusTitleAr               *string    `json:"response_status_title_ar"`
	ResponseStatusType                  *string    `json:"response_status_type"`
	ResponseStatusColorCode             *string    `json:"response_status_color_code"`
	ResponseComment                     *string    `json:"response_comment"`
	ResponseIsException                 *bool      `json:"response_is_exception"`
	ResponseExceptionJustification      *string    `json:"response_exception_justification"`
	ResponseCreatedByID                 *uint      `json:"response_created_by_id"`
	ResponseCreatedUserFirstNameEn      *string    `json:"response_created_user_first_name_en"`
	ResponseCreatedUserFirstNameAr      *string    `json:"response_created_user_first_name_ar"`
	ResponseCreatedUserLastNameEn       *string    `json:"response_created_user_last_name_en"`
	ResponseCreatedUserLastNameAr       *string    `json:"response_created_user_last_name_ar"`
	ResponseCreatedUserInitial          *string    `json:"response_created_user_initial"`
	ResponseCreatedUserColor            *string    `json:"response_created_user_color"`
	ResponseCreatedUserEmail            *string    `json:"response_created_user_email"`
	ResponseCreatedUserProfileImagePath *string    `json:"response_created_user_profile_image_path"`
	ResponseCreatedAt                   *time.Time `json:"response_created_at"`
	ResponseUpdatedByID                 *uint      `json:"response_updated_by_id"`
	ResponseUpdatedUserFirstNameEn      *string    `json:"response_updated_user_first_name_en"`
	ResponseUpdatedUserFirstNameAr      *string    `json:"response_updated_user_first_name_ar"`
	ResponseUpdatedUserLastNameEn       *string    `json:"response_updated_user_last_name_en"`
	ResponseUpdatedUserLastNameAr       *string    `json:"response_updated_user_last_name_ar"`
	ResponseUpdatedUserInitial          *string    `json:"response_updated_user_initial"`
	ResponseUpdatedUserColor            *string    `json:"response_updated_user_color"`
	ResponseUpdatedUserEmail            *string    `json:"response_updated_user_email"`
	ResponseUpdatedUserProfileImagePath *string    `json:"response_updated_user_profile_image_path"`
	ResponseUpdatedAt                   *time.Time `json:"response_updated_at"`
}

type ProjectSecurityGateChecklistListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
}

type ProjectSecurityGateChecklistListResponseDTO struct {
	Items      []*ProjectSecurityGateChecklistListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO           `json:"pagination"`
}

type ProjectSecurityGateChecklistDetailDTO struct {
	ID                     uint                                                `json:"id"`
	SecurityGateID         uint                                                `json:"security_gate_id"`
	Requirement            *ProjectSecurityGateChecklistTemplateRequirementDTO `json:"requirement"`
	ResponseStatus         *ProjectSecurityGateChecklistResponseStatusDTO      `json:"response_status"`
	Comment                *string                                             `json:"comment"`
	IsException            bool                                                `json:"is_exception"`
	ExceptionJustification *string                                             `json:"exception_justification"`
	CreatedByUser          *organizationDTOs.UserMinimalDTO                    `json:"created_by"`
	CreatedAt              *time.Time                                          `json:"created_at"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO                    `json:"updated_by"`
	UpdatedAt              *time.Time                                          `json:"updated_at"`
	Documents              []*ProjectSecurityGateChecklistResponseDocumentDTO  `json:"documents"`
	Comments               []*ProjectSecurityGateChecklistCommentDTO           `json:"comments"`
	Controls               []*organizationDTOs.ControlMinimalListDTO           `json:"controls"`
}

type ProjectSecurityGateChecklistResponseDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}

type ProjectSecurityGateChecklistResponseFormDTO struct {
	IsException            bool                          `json:"is_exception"`
	ResponseStatusID       *uint                         `json:"response_status_id" binding:"required_if=IsException false"`
	ExceptionJustification *string                       `json:"exception_justification" binding:"required_if=IsException true"`
	Comment                *string                       `json:"comment" binding:"omitempty"`
	Documents              []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type ProjectSecurityGateChecklistSummaryDTO struct {
	Total                int64   `json:"total"`
	Critical             int64   `json:"critical"`
	Completed            int64   `json:"completed"`
	Outstanding          int64   `json:"outstanding"`
	CompletionPercentage float64 `json:"completion_percentage"`
}
