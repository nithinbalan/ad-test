package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
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
	ID              uint                             `json:"id"`
	ReferenceNumber string                           `json:"reference_number"`
	ProjectID       uint                             `json:"project_id"`
	Decision        string                           `json:"decision"`
	DecisionNotes   *string                          `json:"decision_notes"`
	Reason          *string                          `json:"reason"`
	DecisionDate    string                           `json:"decision_date"`
	RequestedByUser *organizationDTOs.UserMinimalDTO `json:"requested_by_user"`
	CreatedByUser   *organizationDTOs.UserMinimalDTO `json:"created_by_user"`
	UpdatedByUser   *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
	CreatedAt       string                           `json:"created_at"`
	UpdatedAt       *string                          `json:"updated_at"`
	Documents       []*ProjectDecisionLogDocumentDTO `json:"documents"`
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
