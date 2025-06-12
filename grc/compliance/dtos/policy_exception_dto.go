package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	mastersDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type PolicyExceptionFormDTO struct {
	DocumentID                 uint    `json:"document_id" binding:"required"`
	DocumentVersionClauseID    *uint   `json:"document_version_clause_id"`
	DocumentRequirementGroupID *uint   `json:"document_requirement_group_id"`
	DocumentRequirementID      *uint   `json:"document_requirement_id"`
	TitleEn                    string  `json:"title_en" binding:"required"`
	TitleAr                    string  `json:"title_ar" binding:"required"`
	Requirements               *string `json:"requirements"`
	Justification              *string `json:"justification"`
	StartDate                  string  `json:"start_date" binding:"required"`
	EndDate                    string  `json:"end_date" binding:"required"`
	RequestedByID              uint    `json:"requested_by_id" binding:"required"`
	OrganizationUnitID         uint    `json:"organization_unit_id" binding:"required"`
}

type PolicyExceptionListDTO struct {
	ID                              uint      `json:"id"`
	DocumentID                      uint      `json:"document_id"`
	DocumentVersionClauseID         uint      `json:"document_version_clause_id"`
	DocumentRequirementGroupID      uint      `json:"document_requirement_group_id"`
	DocumentRequirementID           uint      `json:"document_requirement_id"`
	ReferenceNumber                 string    `json:"reference_number"`
	DocumentTitleEn                 string    `json:"document_title_en"`
	DocumentTitleAr                 string    `json:"document_title_ar"`
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
	DocumentRequirementTitleEn      string    `json:"document_requirement_title_en"`
	DocumentRequirementTitleAr      string    `json:"document_requirement_title_ar"`
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
	ID              uint       `json:"id"`
	ReferenceNumber string     `json:"reference_number"`
	Requirements    *string    `json:"requirements"`
	Justification   *string    `json:"justification"`
	StartDateString string     `json:"start_date"`
	EndDateString   string     `json:"end_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`

	// Nested DTOs
	Document                 *documentDTOs.DocumentDetailsDTO          `json:"document"`
	Risks                    []PolicyExceptionRiskDTO                  `json:"risks"`
	DocumentVersion          *documentDTOs.DocumentVersionDTO          `json:"document_version"`
	DocumentVersionClause    *documentDTOs.DocumentClauseDetailsDTO    `json:"document_version_clause"`
	DocumentRequirementGroup *documentDTOs.DocumentRequirementGroupDTO `json:"document_requirement_group"`
	DocumentRequirement      *documentDTOs.DocumentRequirementDTO      `json:"document_requirement"`
	RequestedBy              organizationDTOs.UserMinimalDTO           `json:"requested_by"`
	OrganizationUnit         settingsDTOs.OrganizationUnitDTO          `json:"organization_unit"`
	WorkflowStatus           *mastersDTOs.WorkflowStatusDTO            `json:"workflow_status"`
	CreatedByUser            organizationDTOs.UserMinimalDTO           `json:"created_by"`
	UpdatedByUser            *organizationDTOs.UserMinimalDTO          `json:"updated_by"`
}

type PolicyExceptionParsedDates struct {
	StartDate time.Time
	EndDate   time.Time
}

type MapPolicyExceptionRisksDTO struct {
	RiskIDs []uint `json:"risk_ids" binding:"required"`
}

type RiskMinimalDTO struct {
	ID              uint   `json:"id"`
	ReferenceNumber string `json:"reference_number"`
	Title           string `json:"title"`
}

type PolicyExceptionRiskDTO struct {
	ID                uint                             `json:"id"`
	PolicyExceptionID uint                             `json:"policy_exception_id"`
	RiskRegisterID    uint                             `json:"risk_register_id"`
	RiskRegister      RiskMinimalDTO                   `json:"risk"`
	CreatedAt         time.Time                        `json:"created_at"`
	CreatedBy         *organizationDTOs.UserMinimalDTO `json:"created_by"` // nested user details
}
