package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type PenetrationTestFindingListDTO struct {
	ID                        uint    `json:"id"`
	ReferenceNumber           string  `json:"reference_number"`
	PenetrationTestID         uint    `json:"penetration_test_id"`
	Title                     string  `json:"title"`
	ActionPlanStatusID        uint    `json:"action_plan_status_id"`
	Percentage                *uint   `json:"percentage"`
	ActionPlanStatusTitleEn   string  `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr   string  `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode string  `json:"action_plan_status_color_code"`
	ActionPlanStatusType      string  `json:"action_plan_status_type"`
	Description               *string `json:"description"`
	Recommendation            *string `json:"recommendation"`
	SeverityLevelID           uint    `json:"severity_level_id"`
	SeverityLevelTitleEn      string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr      string  `json:"severity_level_title_ar"`
	SeverityLevelColorCode    string  `json:"severity_level_color_code"`
	OwnerID                   uint    `json:"owner_id"`
	OwnerFirstNameEn          string  `json:"owner_first_name_en"`
	OwnerLastNameEn           string  `json:"owner_last_name_en"`
	OwnerFirstNameAr          string  `json:"owner_first_name_ar"`
	OwnerLastNameAr           string  `json:"owner_last_name_ar"`
	OwnerInitial              string  `json:"owner_initial"`
	OwnerColor                string  `json:"owner_color"`
	OwnerEmail                string  `json:"owner_email"`
	OwnerProfileImagePath     string  `json:"owner_profile_image_path"`
	TargetRemediationDate     *string `json:"target_remediation_date"`
	CreatedAt                 string  `json:"created_at"`
	UpdatedAt                 *string `json:"updated_at"`
}

type PenetrationTestFindingFormDTO struct {
	PenetrationTestID     uint    `json:"penetration_test_id" binding:"required"`
	Title                 string  `json:"title" binding:"required,max=255"`
	SeverityLevelID       uint    `json:"severity_level_id" binding:"required"`
	Description           *string `json:"description"`
	Recommendation        *string `json:"recommendation"`
	OwnerID               uint    `json:"owner_id" binding:"required"`
	TargetRemediationDate *string `json:"target_remediation_date" binding:"omitempty,datetime=2006-01-02"`
}

type PenetrationTestFindingDetailDTO struct {
	ID                    uint                             `json:"id"`
	ReferenceNumber       string                           `json:"reference_number"`
	PenetrationTestID     uint                             `json:"penetration_test_id"`
	Title                 string                           `json:"title"`
	ActionPlanStatusID    uint                             `json:"action_plan_status_id"`
	Percentage            *uint                            `json:"percentage"`
	ActionPlanStatus      masterDTOs.ActionPlanStatusDTO   `json:"action_plan_status"`
	SeverityLevelID       uint                             `json:"severity_level_id"`
	SeverityLevel         masterDTOs.SeverityLevelDTO      `json:"severity_level"`
	Description           *string                          `json:"description"`
	Recommendation        *string                          `json:"recommendation"`
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
	Filters           map[string][]string
	QuickFilter       []string
	DateFrom          string
	DateTo            string
}

type PenetrationTestFindingListResponseDTO struct {
	Items      []*PenetrationTestFindingListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}
