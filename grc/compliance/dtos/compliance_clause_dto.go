package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
)

type ComplianceClauseDTO struct {
	ID       int    `json:"id"`
	ClauseNo string `json:"clause_no"`
	Title    string `json:"title"`
}

// ComplianceClauseListDTO represents a compliance clause in list responses.
type ComplianceClauseListDTO struct {
	ID                   uint                       `json:"id"`
	ComplianceRegisterID uint                       `json:"compliance_register_id"`
	ParentClauseID       *uint                      `json:"parent_clause_id"`
	ClauseNo             string                     `json:"clause_no"`
	TitleEn              string                     `json:"title_en"`
	TitleAr              *string                    `json:"title_ar"`
	DescriptionEn        *string                    `json:"description_en"`
	DescriptionAr        *string                    `json:"description_ar"`
	ScopeEn              *string                    `json:"scope_en"`
	ScopeAr              *string                    `json:"scope_ar"`
	ObjectiveEn          *string                    `json:"objective_en"`
	ObjectiveAr          *string                    `json:"objective_ar"`
	CreatedByID          *uint                      `json:"created_by_id,omitempty"`
	CreatedAt            time.Time                  `json:"created_at"`
	UpdatedByID          *uint                      `json:"updated_by_id,omitempty"`
	UpdatedAt            *time.Time                 `json:"updated_at,omitempty"`
	ChildClauses         []*ComplianceClauseListDTO `json:"child_clauses"`
}

// ComplianceClauseListQueryParamsDTO represents filtering and pagination for compliance clauses.
type ComplianceClauseListQueryParamsDTO struct {
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

// ComplianceClauseListResponseDTO represents a paginated response of compliance clauses.
type ComplianceClauseListResponseDTO struct {
	Items      []*ComplianceClauseListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

// ComplianceClauseFormDTO represents the request body for creating or updating a compliance clause.
type ComplianceClauseFormDTO struct {
	ParentClauseID   *uint   `json:"parent_clause_id,omitempty"`
	PreviousClauseID *uint   `json:"previous_clause_id,omitempty"`
	ClauseNo         string  `json:"clause_no,omitempty"` // Auto-generated if not provided
	TitleEn          string  `json:"title_en" binding:"required"`
	TitleAr          *string `json:"title_ar,omitempty"`
	DescriptionEn    *string `json:"description_en,omitempty"`
	DescriptionAr    *string `json:"description_ar,omitempty"`
	ScopeEn          *string `json:"scope_en,omitempty"`
	ScopeAr          *string `json:"scope_ar,omitempty"`
	ObjectiveEn      *string `json:"objective_en,omitempty"`
	ObjectiveAr      *string `json:"objective_ar,omitempty"`
}

// ComplianceClauseDetailDTO represents detailed information for a compliance clause.
type ComplianceClauseDetailDTO struct {
	ID                     uint                              `json:"id"`
	ComplianceRegisterID   uint                              `json:"compliance_register_id"`
	ParentClauseID         *uint                             `json:"parent_clause_id"`
	ClauseNo               string                            `json:"clause_no"`
	TitleEn                string                            `json:"title_en"`
	TitleAr                *string                           `json:"title_ar"`
	DescriptionEn          *string                           `json:"description_en"`
	DescriptionAr          *string                           `json:"description_ar"`
	ScopeEn                *string                           `json:"scope_en"`
	ScopeAr                *string                           `json:"scope_ar"`
	ObjectiveEn            *string                           `json:"objective_en"`
	ObjectiveAr            *string                           `json:"objective_ar"`
	ComplianceRequirements []*ComplianceRequirementDetailDTO `json:"compliance_requirements"`
}
