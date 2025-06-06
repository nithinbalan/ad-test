package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type PolicyExceptionFormDTO struct {
	DocumentID                 uint    `json:"document_id" binding:"required"`
	DocumentVersionClauseID    *uint   `json:"document_version_clause_id"`
	DocumentRequirementGroupID *uint   `json:"document_requirement_group_id"`
	ReferenceNumber            string  `json:"reference_number"` // generated in service, not required in form
	TitleEn                    string  `json:"title_en" binding:"required"`
	TitleAr                    string  `json:"title_ar" binding:"required"`
	Requirements               *string `json:"requirements"`
	Justification              *string `json:"justification"`
	StartDate                  string  `json:"start_date" binding:"required"` // consider using time.Time if you parse it before saving
	EndDate                    *string `json:"end_date"`                      // same here if optional
	RequestedByID              uint    `json:"requested_by_id" binding:"required"`
	OrganizationUnitID         uint    `json:"organization_unit_id" binding:"required"`
	WorkflowStatusID           *uint   `json:"workflow_status_id"` // overridden to pending in service, so optional
}

type PolicyExceptionListDTO struct {
	ID                              uint      `json:"id"`
	ReferenceNumber                 string    `json:"reference_number"`
	TitleEn                         string    `json:"title_en"`
	TitleAr                         string    `json:"title_ar"`
	StartDate                       string    `json:"start_date"`
	EndDate                         *string   `json:"end_date"`
	RequestedByID                   uint      `json:"requested_by_id"`
	RequestedByFirstNameEn          string    `json:"requested_by_first_name_en"`
	RequestedByFirstNameAr          string    `json:"requested_by_first_name_ar"`
	RequestedByLastNameEn           string    `json:"requested_by_last_name_en"`
	RequestedByLastNameAr           string    `json:"requested_by_last_name_ar"`
	RequestedByInitial              string    `json:"requested_by_initial"`
	RequestedByEmail                string    `json:"requested_by_email"`
	RequestedByProfileImage         string    `json:"requested_by_profile_image_path"`
	WorkflowStatusTitleEn           string    `json:"workflow_status_title_en"`
	WorkflowStatusTitleAr           string    `json:"workflow_status_title_ar"`
	WorkflowStatusColor             string    `json:"workflow_status_color"`
	DocumentRequirementGroupTitleEn string    `json:"document_requirement_group_title_en"`
	DocumentRequirementGroupTitleAr string    `json:"document_requirement_group_title_ar"`
	DocumentVersionTitleEn          string    `json:"document_version_title_en"`
	DocumentVersionTitleAr          string    `json:"document_version_title_ar"`
	CreatedAt                       time.Time `json:"created_at"`
	CreatedByID                     uint      `json:"created_by_id"`
	CreatedByFirstNameEn            string    `json:"created_by_first_name_en"`
	CreatedByFirstNameAr            string    `json:"created_by_first_name_ar"`
	CreatedByLastNameEn             string    `json:"created_by_last_name_en"`
	CreatedByLastNameAr             string    `json:"created_by_last_name_ar"`
	CreatedByInitial                string    `json:"created_by_initial"`
	CreatedByColor                  string    `json:"created_by_color"`
	CreatedByEmail                  string    `json:"created_by_email"`
	CreatedByProfileImagePath       string    `json:"created_by_profile_image_path"`
	UpdatedByID                     *uint     `json:"updated_by_id"`
	UpdatedByFirstNameEn            string    `json:"updated_by_first_name_en"`
	UpdatedByFirstNameAr            string    `json:"updated_by_first_name_ar"`
	UpdatedByLastNameEn             string    `json:"updated_by_last_name_en"`
	UpdatedByLastNameAr             string    `json:"updated_by_last_name_ar"`
	UpdatedByInitial                string    `json:"updated_by_initial"`
	UpdatedByColor                  string    `json:"updated_by_color"`
	UpdatedByEmail                  string    `json:"updated_by_email"`
	UpdatedByProfileImagePath       string    `json:"updated_by_profile_image_path"`
	UpdatedAt                       *string   `json:"updated_at"`
}

type PolicyExceptionListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	DateFrom    string
	DateTo      string
	Export      bool
}

type PolicyExceptionListResponseDTO struct {
	Items      []*PolicyExceptionListDTO    `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type PolicyExceptionDetailDTO struct {
	ID                      uint       `json:"id"`
	DocumentID              uint       `json:"document_id"`
	DocumentVersionClauseID *uint      `json:"document_version_clause_id"`
	ReferenceNumber         string     `json:"reference_number"`
	TitleEn                 string     `json:"title_en"`
	TitleAr                 string     `json:"title_ar"`
	Requirements            *string    `json:"requirements"`
	Justification           *string    `json:"justification"`
	StartDateString         string     `json:"start_date"`
	EndDateString           string     `json:"end_date"`
	RequestedByID           uint       `json:"requested_by_id"`
	RequestedByName         string     `json:"requested_by_name"`
	OrganizationUnitID      uint       `json:"organization_unit_id"`
	OrganizationUnitName    string     `json:"organization_unit_name"`
	WorkflowStatusID        *uint      `json:"workflow_status_id"`
	WorkflowStatusTitleEn   string     `json:"workflow_status_title_en"`
	WorkflowStatusColor     string     `json:"workflow_status_color"`
	CreatedAt               time.Time  `json:"created_at"`
	CreatedByID             uint       `json:"created_by_id"`
	CreatedByName           string     `json:"created_by_name"`
	UpdatedAt               *time.Time `json:"updated_at"`
	UpdatedByID             *uint      `json:"updated_by_id"`
	UpdatedByName           *string    `json:"updated_by_name"`
}

type PolicyExceptionParsedDates struct {
	StartDate time.Time
	EndDate   time.Time
}

type MapPolicyExceptionRisksDTO struct {
	RiskIDs []uint `json:"risk_ids" binding:"required"`
}
