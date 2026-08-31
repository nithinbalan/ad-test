package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type PenetrationTestFindingListDTO struct {
	ID                     uint    `json:"id"`
	ReferenceNumber        string  `json:"reference_number"`
	PenetrationTestID      uint    `json:"penetration_test_id"`
	Title                  string  `json:"title"`
	Status                 string  `json:"status"`
	SeverityLevelID        uint    `json:"severity_level_id"`
	SeverityLevelTitleEn   string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr   string  `json:"severity_level_title_ar"`
	SeverityLevelColorCode string  `json:"severity_level_color_code"`
	OwnerID                uint    `json:"owner_id"`
	OwnerFirstNameEn       string  `json:"owner_first_name_en"`
	OwnerLastNameEn        string  `json:"owner_last_name_en"`
	OwnerFirstNameAr       string  `json:"owner_first_name_ar"`
	OwnerLastNameAr        string  `json:"owner_last_name_ar"`
	TargetRemediationDate  *string `json:"target_remediation_date"`
	CreatedAt              string  `json:"created_at"`
	UpdatedAt              *string `json:"updated_at"`
}

type PenetrationTestFindingFormDTO struct {
	PenetrationTestID     uint    `json:"penetration_test_id" binding:"required"`
	Title                 string  `json:"title" binding:"required,max=255"`
	Status                string  `json:"status" binding:"omitempty,oneof='Open' 'In Remediation' 'Retest Requested' 'Closed'"`
	SeverityLevelID       uint    `json:"severity_level_id" binding:"required"`
	Description           string  `json:"description" binding:"required"`
	Recommendation        string  `json:"recommendation" binding:"required"`
	OwnerID               uint    `json:"owner_id" binding:"required"`
	TargetRemediationDate *string `json:"target_remediation_date" binding:"omitempty,datetime=2006-01-02"`
}

type PenetrationTestFindingDetailDTO struct {
	ID                    uint                             `json:"id"`
	ReferenceNumber       string                           `json:"reference_number"`
	PenetrationTestID     uint                             `json:"penetration_test_id"`
	Title                 string                           `json:"title"`
	Status                string                           `json:"status"`
	SeverityLevelID       uint                             `json:"severity_level_id"`
	SeverityLevel         masterDTOs.SeverityLevelDTO      `json:"severity_level"`
	Description           string                           `json:"description"`
	Recommendation        string                           `json:"recommendation"`
	OwnerID               uint                             `json:"owner_id"`
	Owner                 organizationDTOs.UserMinimalDTO  `json:"owner"`
	TargetRemediationDate *time.Time                       `json:"target_remediation_date"`
	CreatedAt             time.Time                        `json:"created_at"`
	UpdatedAt             *time.Time                       `json:"updated_at"`
	CreatedByUser         organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser         *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}

type PenetrationTestFindingListQueryParamsDTO struct {
	Sort              string
	SortBy            string
	Keyword           string
	Page              int
	PageSize          int
	PenetrationTestID uint
}

type PenetrationTestFindingListResponseDTO struct {
	Items      []*PenetrationTestFindingListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}
