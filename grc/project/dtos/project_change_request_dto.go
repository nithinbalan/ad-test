package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ProjectChangeRequestListDTO struct {
	ID                          uint    `json:"id"`
	ReferenceNumber             string  `json:"reference_number"`
	ProjectID                   uint    `json:"project_id"`
	Title                       string  `json:"title"`
	Description                 *string `json:"description"`
	BusinessJustification       *string `json:"business_justification"`
	AnticipatedBenefits         *string `json:"anticipated_benefits"`
	RisksOfDeferral             *string `json:"risks_of_deferral"`
	RequestedByUserID           uint    `json:"requested_by_user_id"`
	RequestedByUserFirstNameEn  string  `json:"requested_by_user_first_name_en"`
	RequestedByUserFirstNameAr  string  `json:"requested_by_user_first_name_ar"`
	RequestedByUserLastNameEn   string  `json:"requested_by_user_last_name_en"`
	RequestedByUserLastNameAr   string  `json:"requested_by_user_last_name_ar"`
	RequestedByUserInitial      string  `json:"requested_by_user_initial"`
	RequestedByUserColor        string  `json:"requested_by_user_color"`
	RequestedByUserEmail        string  `json:"requested_by_user_email"`
	RequestedByUserImagePath    string  `json:"requested_by_user_profile_image_path"`
	RequestedOn                 *string `json:"requested_on"`
	WorkflowStatusID            *uint   `json:"workflow_status_id"`
	WorkflowStatusTitleEn       *string `json:"workflow_status_title_en"`
	WorkflowStatusTitleAr       *string `json:"workflow_status_title_ar"`
	WorkflowStatusColorCode     *string `json:"workflow_status_color_code"`
	WorkflowStatusType          *string `json:"workflow_status_type"`
	NextReviewLevel             *int    `json:"next_review_level"`
	SubmittedAt                 *string `json:"submitted_at"`
	SubmittedUserID             *uint   `json:"submitted_user_id"`
	SubmittedUserFirstNameEn    *string `json:"submitted_user_first_name_en"`
	SubmittedUserFirstNameAr    *string `json:"submitted_user_first_name_ar"`
	SubmittedUserLastNameEn     *string `json:"submitted_user_last_name_en"`
	SubmittedUserLastNameAr     *string `json:"submitted_user_last_name_ar"`
	ApprovedAt                  *string `json:"approved_at"`
	ApprovedUserID              *uint   `json:"approved_user_id"`
	ApprovedUserFirstNameEn     *string `json:"approved_user_first_name_en"`
	ApprovedUserFirstNameAr     *string `json:"approved_user_first_name_ar"`
	ApprovedUserLastNameEn      *string `json:"approved_user_last_name_en"`
	ApprovedUserLastNameAr      *string `json:"approved_user_last_name_ar"`
	RejectedAt                  *string `json:"rejected_at"`
	RejectedUserID              *uint   `json:"rejected_user_id"`
	RejectedUserFirstNameEn     *string `json:"rejected_user_first_name_en"`
	RejectedUserFirstNameAr     *string `json:"rejected_user_first_name_ar"`
	RejectedUserLastNameEn      *string `json:"rejected_user_last_name_en"`
	RejectedUserLastNameAr      *string `json:"rejected_user_last_name_ar"`
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

type ProjectChangeRequestListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	DateFrom    string
	DateTo      string
}

type ProjectChangeRequestListResponseDTO struct {
	Items      []*ProjectChangeRequestListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO   `json:"pagination"`
}

type ProjectChangeRequestDetailDTO struct {
	ID                    uint                               `json:"id"`
	ReferenceNumber       string                             `json:"reference_number"`
	Project               ProjectDetailDTO                   `json:"project"`
	Title                 string                             `json:"title"`
	Description           *string                            `json:"description"`
	BusinessJustification *string                            `json:"business_justification"`
	AnticipatedBenefits   *string                            `json:"anticipated_benefits"`
	RisksOfDeferral       *string                            `json:"risks_of_deferral"`
	RequestedByUser       organizationDTOs.UserMinimalDTO    `json:"requested_by_user"`
	RequestedOnString     *string                            `json:"requested_on"`
	WorkflowStatus        *masterDTOs.WorkflowStatusDTO      `json:"workflow_status"`
	NextReviewLevel       *int                               `json:"next_review_level"`
	SubmittedByUser       *organizationDTOs.UserMinimalDTO   `json:"submitted_by_user"`
	SubmittedAtString     *string                            `json:"submitted_at"`
	ApprovedByUser        *organizationDTOs.UserMinimalDTO   `json:"approved_by_user"`
	ApprovedAtString      *string                            `json:"approved_at"`
	RejectedByUser        *organizationDTOs.UserMinimalDTO   `json:"rejected_by_user"`
	RejectedAtString      *string                            `json:"rejected_at"`
	CreatedByUser         organizationDTOs.UserMinimalDTO    `json:"created_by"`
	UpdatedByUser         *organizationDTOs.UserMinimalDTO   `json:"updated_by"`
	CreatedAt             time.Time                          `json:"created_at"`
	UpdatedAt             *time.Time                         `json:"updated_at"`
	ImpactAssessments     []*ProjectImpactAssessmentDTO      `json:"impact_assessments"`
	Documents             []*ProjectChangeRequestDocumentDTO `json:"documents"`
}

type ProjectChangeRequestDocumentDTO struct {
	ID           uint    `json:"id"`
	URL          string  `json:"url"`
	Name         string  `json:"file_name"`
	OriginalName *string `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
	Token        string  `json:"token"`
}

type ProjectImpactAssessmentDTO struct {
	ID                            uint                            `json:"id"`
	ProjectChangeRequestID        uint                            `json:"project_change_request_id"`
	ProjectImpactAssessmentTypeID uint                            `json:"project_impact_assessment_type_id"`
	ProjectImpactAssessmentType   *ProjectImpactAssessmentTypeDTO `json:"project_impact_assessment_type"`
	CurrentDateString             *string                         `json:"current_date"`
	RevisedDateString             *string                         `json:"revised_date"`
	CurrentBudget                 *float64                        `json:"current_budget"`
	RevisedBudget                 *float64                        `json:"revised_budget"`
	CurrentResourceCount          *int                            `json:"current_resource_count"`
	RevisedResourceCount          *int                            `json:"revised_resource_count"`
	TechnicalComplexity           *string                         `json:"technical_complexity"`
}

type ProjectImpactAssessmentFormDTO struct {
	ProjectImpactAssessmentTypeID uint     `json:"project_impact_assessment_type_id" binding:"required"`
	CurrentDate                   *string  `json:"current_date"`
	RevisedDate                   *string  `json:"revised_date"`
	CurrentBudget                 *float64 `json:"current_budget"`
	RevisedBudget                 *float64 `json:"revised_budget"`
	CurrentResourceCount          *int     `json:"current_resource_count"`
	RevisedResourceCount          *int     `json:"revised_resource_count"`
	TechnicalComplexity           *string  `json:"technical_complexity"`
}

type ProjectChangeRequestFormDTO struct {
	Title                 string                           `json:"title" binding:"required"`
	Description           *string                          `json:"description"`
	BusinessJustification *string                          `json:"business_justification"`
	AnticipatedBenefits   *string                          `json:"anticipated_benefits"`
	RisksOfDeferral       *string                          `json:"risks_of_deferral"`
	RequestedByUserID     uint                             `json:"requested_by_user_id" binding:"required"`
	RequestedOn           string                           `json:"requested_on" binding:"required"`
	ImpactAssessments     []ProjectImpactAssessmentFormDTO `json:"impact_assessments"`
	Documents             []*commonDTOs.DocumentFormDTO    `json:"documents" binding:"omitempty"`
}

type ProjectChangeRequestSummaryDTO struct {
	Total     int64 `json:"total"`
	Submitted int64 `json:"submitted"`
	Approved  int64 `json:"approved"`
	Rejected  int64 `json:"rejected"`
	Pending   int64 `json:"pending"`
}
