package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type IdentityReviewUserFormDTO struct {
	Username      string `json:"username" binding:"required,max=255"`
	DepartmentID  uint   `json:"department_id" binding:"required"`
	PrivilegeID   *uint  `json:"privilege_id" binding:"omitempty"`
	AccountTypeID *uint  `json:"account_type_id" binding:"omitempty"`
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
	ID                   uint       `json:"id"`
	Username             string     `json:"username"`
	DepartmentID         uint       `json:"department_id"`
	DepartmentTitleEn    string     `json:"department_title_en"`
	DepartmentTitleAr    string     `json:"department_title_ar"`
	PrivilegeID          *uint      `json:"privilege_id"`
	PrivilegeTitleEn     *string    `json:"privilege_title_en"`
	PrivilegeTitleAr     *string    `json:"privilege_title_ar"`
	AccountTypeID        *uint      `json:"account_type_id"`
	AccountTypeTitleEn   *string    `json:"account_type_title_en"`
	AccountTypeTitleAr   *string    `json:"account_type_title_ar"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            *time.Time `json:"updated_at"`
	CreatedByID          uint       `json:"created_by_id"`
	CreatedByFirstNameEn string     `json:"created_by_first_name_en"`
	CreatedByLastNameEn  string     `json:"created_by_last_name_en"`
	CreatedByInitial     string     `json:"created_by_initial"`
	UpdatedByID          *uint      `json:"updated_by_id"`
	UpdatedByFirstNameEn *string    `json:"updated_by_first_name_en"`
	UpdatedByLastNameEn  *string    `json:"updated_by_last_name_en"`
	UpdatedByInitial     *string    `json:"updated_by_initial"`
}

type IdentityReviewUserListResponseDTO struct {
	Items      []*IdentityReviewUserListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
