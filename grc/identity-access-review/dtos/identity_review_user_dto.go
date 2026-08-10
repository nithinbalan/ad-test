package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type IdentityReviewUserFormDTO struct {
	IdentityReviewCampaignID uint   `json:"identity_review_campaign_id" binding:"required"`
	Username                 string `json:"username" binding:"required,max=255"`
	UserIdentifier           string `json:"user_identifier" binding:"required,max=255"`
	LastLogin                string `json:"last_login" binding:"required,max=255"`
	CurrentRole              string `json:"current_role" binding:"required"`
	DepartmentID             uint   `json:"department_id" binding:"required"`
	PrivilegeID              *uint  `json:"privilege_id" binding:"omitempty"`
	AccountTypeID            *uint  `json:"account_type_id" binding:"omitempty"`
}

type IdentityReviewUserAccountDecisionUpdateFormDTO struct {
	AccountDecisionID uint `json:"account_decision_id" binding:"required"`
}

type IdentityReviewUserPrivilegeUpdateFormDTO struct {
	PrivilegeID *uint `json:"privilege_id" binding:"omitempty"`
}

type IdentityReviewUserListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type IdentityReviewUserListDTO struct {
	ID                         uint       `json:"id"`
	IdentityReviewCampaignID   uint       `json:"identity_review_campaign_id"`
	IdentityReviewCampaignName string     `json:"identity_review_campaign_name"`
	Username                   string     `json:"username"`
	UserIdentifier             *string    `json:"user_identifier"`
	LastLogin                  *string    `json:"last_login"`
	CurrentRole                *string    `json:"current_role"`
	DepartmentID               uint       `json:"department_id"`
	DepartmentTitleEn          string     `json:"department_title_en"`
	DepartmentTitleAr          string     `json:"department_title_ar"`
	PrivilegeID                *uint      `json:"privilege_id"`
	PrivilegeTitleEn           *string    `json:"privilege_title_en"`
	PrivilegeTitleAr           *string    `json:"privilege_title_ar"`
	AccountTypeID              *uint      `json:"account_type_id"`
	AccountTypeTitleEn         *string    `json:"account_type_title_en"`
	AccountTypeTitleAr         *string    `json:"account_type_title_ar"`
	AccountDecisionID          *uint      `json:"account_decision_id"`
	AccountDecisionTitleEn     *string    `json:"account_decision_title_en"`
	AccountDecisionTitleAr     *string    `json:"account_decision_title_ar"`
	AccountDecisionColorCode   *string    `json:"account_decision_color_code"`
	AccountDecisionType        *string    `json:"account_decision_type"`
	CreatedAt                  time.Time  `json:"created_at"`
	UpdatedAt                  *time.Time `json:"updated_at"`
	CreatedByID                uint       `json:"created_by_id"`
	CreatedByFirstNameEn       string     `json:"created_by_first_name_en"`
	CreatedByLastNameEn        string     `json:"created_by_last_name_en"`
	CreatedByInitial           string     `json:"created_by_initial"`
	UpdatedByID                *uint      `json:"updated_by_id"`
	UpdatedByFirstNameEn       *string    `json:"updated_by_first_name_en"`
	UpdatedByLastNameEn        *string    `json:"updated_by_last_name_en"`
	UpdatedByInitial           *string    `json:"updated_by_initial"`
}

type IdentityReviewUserListResponseDTO struct {
	Items      []*IdentityReviewUserListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
