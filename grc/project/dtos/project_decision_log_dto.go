package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type ProjectDecisionLogListDTO struct {
	ID                          uint    `json:"id"`
	ReferenceNumber             string  `json:"reference_number"`
	ProjectID                   uint    `json:"project_id"`
	Decision                    string  `json:"decision"`
	DecisionNotes               *string `json:"decision_notes"`
	Reason                      *string `json:"reason"`
	DecisionDate                string  `json:"decision_date"`
	RequestedByUserID           uint    `json:"requested_by_user_id"`
	RequestedByUserFirstNameEn  string  `json:"requested_by_user_first_name_en"`
	RequestedByUserFirstNameAr  string  `json:"requested_by_user_first_name_ar"`
	RequestedByUserLastNameEn   string  `json:"requested_by_user_last_name_en"`
	RequestedByUserLastNameAr   string  `json:"requested_by_user_last_name_ar"`
	RequestedByUserInitial      string  `json:"requested_by_user_initial"`
	RequestedByUserColor        string  `json:"requested_by_user_color"`
	RequestedByUserEmail        string  `json:"requested_by_user_email"`
	RequestedByUserImagePath    string  `json:"requested_by_user_profile_image_path"`
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
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   *string `json:"updated_at"`
}

type ProjectDecisionLogDetailDTO struct {
	ID                uint                             `json:"id"`
	ReferenceNumber   string                           `json:"reference_number"`
	ProjectID         uint                             `json:"project_id"`
	Decision          string                           `json:"decision"`
	DecisionNotes     *string                          `json:"decision_notes"`
	Reason            *string                          `json:"reason"`
	DecisionDate      string                           `json:"decision_date"`
	RequestedByUser   *organizationDTOs.UserMinimalDTO `json:"requested_by_user"`
	WorkflowStatus    *masterDTOs.WorkflowStatusDTO    `json:"workflow_status"`
	NextReviewLevel   *int                             `json:"next_review_level"`
	SubmittedByUser   *organizationDTOs.UserMinimalDTO `json:"submitted_by_user"`
	SubmittedAtString *string                          `json:"submitted_at"`
	ApprovedByUser    *organizationDTOs.UserMinimalDTO `json:"approved_by_user"`
	ApprovedAtString  *string                          `json:"approved_at"`
	RejectedByUser    *organizationDTOs.UserMinimalDTO `json:"rejected_by_user"`
	RejectedAtString  *string                          `json:"rejected_at"`
	CreatedByUser     *organizationDTOs.UserMinimalDTO `json:"created_by_user"`
	UpdatedByUser     *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
	CreatedAt         string                           `json:"created_at"`
	UpdatedAt         *string                          `json:"updated_at"`
	Documents         []*ProjectDecisionLogDocumentDTO `json:"documents"`
}

type ProjectDecisionLogDocumentDTO struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	OriginalName *string `json:"original_name"`
	Size         int     `json:"size"`
	Type         string  `json:"type"`
	URL          string  `json:"url"`
	Token        string  `json:"token"`
}

type ProjectDecisionLogFormDTO struct {
	Decision          string                        `json:"decision" binding:"required"`
	DecisionNotes     *string                       `json:"decision_notes"`
	Reason            *string                       `json:"reason"`
	RequestedByUserID uint                          `json:"requested_by_user_id" binding:"required"`
	DecisionDate      string                        `json:"decision_date" binding:"required"`
	Documents         []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type ProjectDecisionLogListQueryParamsDTO struct {
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

type ProjectDecisionLogListResponseDTO struct {
	Data       []*ProjectDecisionLogListDTO `json:"data"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ProjectDecisionLogSummaryDTO struct {
	Total     int64 `json:"total"`
	Submitted int64 `json:"submitted"`
	Approved  int64 `json:"approved"`
	Rejected  int64 `json:"rejected"`
	Pending   int64 `json:"pending"`
}
