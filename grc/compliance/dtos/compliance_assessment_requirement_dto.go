package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ComplianceAssessmentRequirementListDTO struct {
	ID                                   uint    `json:"id"`
	ClauseID                             uint    `json:"clause_id"`
	IsCompleted                          bool    `json:"is_completed"`
	ComplianceClauseNo                   string  `json:"compliance_clause_no"`
	ComplianceClauseTitleEn              string  `json:"compliance_clause_title_en"`
	ComplianceClauseTitleAr              string  `json:"compliance_clause_title_ar"`
	ComplianceRequirementTitleEn         string  `json:"compliance_requirement_title_en"`
	ComplianceRequirementTitleAr         string  `json:"compliance_requirement_title_ar"`
	ComplianceRequirementCategoryTitleEn string  `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr string  `json:"compliance_requirement_category_title_ar"`
	SeverityLevelTitleEn                 string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr                 string  `json:"severity_level_title_ar"`
	SeverityLevelColorCode               string  `json:"severity_level_color_code"`
	CreatedAt                            string  `json:"created_at"`
	UpdatedAt                            *string `json:"updated_at"`
}

type ComplianceAssessmentRequirementDetailDTO struct {
	ID                    int                              `json:"id"`
	ComplianceAssessment  *ComplianceAssessmentDetailDTO   `json:"compliance_assessment"`
	ComplianceRequirement ComplianceRequirementDetailDTO   `json:"compliance_requirement"`
	IsCompleted           bool                             `json:"is_completed"`
	Score                 *int                             `json:"score"`
	Comments              *string                          `json:"comments"`
	ComplianceStatusID    *uint                            `json:"compliance_status_id"`
	ComplianceStatus      *ComplianceStatusDTO             `json:"compliance_status"`
	CreatedAt             time.Time                        `json:"created_at"`
	UpdatedAt             *time.Time                       `json:"updated_at"`
	CreatedByUser         organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser         *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}

type ComplianceAssessmentRequirementUpdateFormDTO struct {
	Score              int     `json:"score" binding:"required"`
	ComplianceStatusID uint    `json:"compliance_status_id" binding:"required"`
	Comments           *string `json:"comments"`
}

type ComplianceAssessmentRequirementListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Filters     map[string][]string
	QuickFilter []string
}
