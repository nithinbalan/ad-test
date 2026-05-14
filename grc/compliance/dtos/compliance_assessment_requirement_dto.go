package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ComplianceAssessmentRequirementListDTO struct {
	ID                                             uint     `json:"id"`
	ClauseID                                       uint     `json:"clause_id"`
	ClauseDescriptionEn                            *string  `json:"clause_description_en"`
	ClauseDescriptionAr                            *string  `json:"clause_description_ar"`
	ComplianceClauseNo                             string   `json:"compliance_clause_no"`
	ComplianceClauseTitleEn                        string   `json:"compliance_clause_title_en"`
	ComplianceClauseTitleAr                        string   `json:"compliance_clause_title_ar"`
	ComplianceRequirementID                        uint     `json:"compliance_requirement_id"`
	ComplianceAssessmentRequirementStatusID        *uint    `json:"compliance_assessment_requirement_status_id"`
	ComplianceRequirementTitleEn                   string   `json:"compliance_requirement_title_en"`
	ComplianceRequirementTitleAr                   string   `json:"compliance_requirement_title_ar"`
	ComplianceRequirementCategoryTitleEn           string   `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr           string   `json:"compliance_requirement_category_title_ar"`
	SeverityLevelID                                uint     `json:"severity_level_id"`
	SeverityLevelTitleEn                           string   `json:"severity_level_title_en"`
	SeverityLevelTitleAr                           string   `json:"severity_level_title_ar"`
	SeverityLevelColorCode                         string   `json:"severity_level_color_code"`
	SeverityLevelType                              string   `json:"severity_level_type"`
	ComplianceAssessmentRequirementStatusTitleEn   string   `json:"compliance_assessment_requirement_status_title_en"`
	ComplianceAssessmentRequirementStatusTitleAr   string   `json:"compliance_assessment_requirement_status_title_ar"`
	ComplianceAssessmentRequirementStatusColorCode string   `json:"compliance_assessment_requirement_status_color_code"`
	ComplianceAssessmentRequirementStatusType      string   `json:"compliance_assessment_requirement_status_type"`
	CreatedAt                                      string   `json:"created_at"`
	UpdatedAt                                      *string  `json:"updated_at"`
	IsAllCompleted                                 bool     `json:"is_all_completed"`
	AverageResidualRiskScore                       *float64 `json:"average_residual_risk_score"`
	AverageResidualRiskRatingID                    *uint    `json:"average_residual_risk_rating_id"`
	AverageResidualRiskRatingTitleEn               *string  `json:"average_residual_risk_rating_title_en"`
	AverageResidualRiskRatingTitleAr               *string  `json:"average_residual_risk_rating_title_ar"`
	AverageResidualRiskRatingColorCode             *string  `json:"average_residual_risk_rating_color_code"`
	AverageResidualRiskRatingType                  *string  `json:"average_residual_risk_rating_type"`
}

type ComplianceAssessmentRequirementListResponseDTO struct {
	Items      []*ComplianceAssessmentRequirementListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO              `json:"pagination"`
}

type ComplianceAssessmentRequirementDetailDTO struct {
	ID                                    int                                                  `json:"id"`
	ComplianceAssessment                  *ComplianceAssessmentDetailDTO                       `json:"compliance_assessment"`
	ComplianceRequirement                 ComplianceRequirementDetailDTO                       `json:"compliance_requirement"`
	Score                                 *float64                                             `json:"score"`
	ComplianceStatusID                    *uint                                                `json:"compliance_status_id"`
	ComplianceStatus                      *ComplianceStatusDTO                                 `json:"compliance_status"`
	ComplianceAssessmentRequirementStatus *ComplianceAssessmentRequirementStatusDTO            `json:"compliance_assessment_requirement_status"`
	CreatedAt                             time.Time                                            `json:"created_at"`
	UpdatedAt                             *time.Time                                           `json:"updated_at"`
	CreatedByUser                         organizationDTOs.UserMinimalDTO                      `json:"created_by_user"`
	UpdatedByUser                         *organizationDTOs.UserMinimalDTO                     `json:"updated_by_user"`
	ResponsibleUnits                      []*ComplianceAssessmentRequirementResponsibleUnitDTO `json:"responsible_units"`
	ResponsibleUnit                       *ComplianceAssessmentRequirementResponsibleUnitDTO   `json:"responsible_unit"`
	CompletedByUser                       *organizationDTOs.UserMinimalDTO                     `json:"completed_by_user"`
	IsAllCompleted                        bool                                                 `json:"is_all_completed"`
	AverageResidualRiskScore              *float64                                             `json:"average_residual_risk_score"`
	AverageResidualRiskRating             *ComplianceAssessmentRatingDTO                       `json:"average_residual_risk_rating"`
}

type ComplianceAssessmentRequirementUpdateFormDTO struct {
	ResponsibleUnitID          uint    `json:"responsible_unit_id" binding:"required"`
	Score                      int     `json:"score" binding:"omitempty,min=0"`
	Comments                   *string `json:"comments"`
	AssessmentFrameworkValueID uint    `json:"assessment_framework_value_id" binding:"required"`
}

type ComplianceAssessmentRequirementDocumentTitleDTO struct {
	ID                                *uint                                         `json:"id"`
	ComplianceAssessmentRequirementID uint                                          `json:"compliance_assessment_requirement_id"`
	RequiredDocumentTitleID           *uint                                         `json:"required_document_title_id"`
	RequiredDocumentTitle             masterDTOs.RequiredDocumentTitleDTO           `json:"required_document_title"`
	Documents                         []*ComplianceAssessmentRequirementDocumentDTO `json:"documents" gorm:"-"`
}

type ComplianceAssessmentRequirementDocumentFormDTO struct {
	ID           *uint  `json:"id"`
	Language     string `json:"language"`
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
	Size         int    `json:"size"`
}

type ComplianceAssessmentRequirementDocumentDTO struct {
	ID           *uint  `json:"id,omitempty"`
	Language     string `json:"language"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
	Size         int    `json:"size"`
}

type ComplianceAssessmentRequirementListQueryParamsDTO struct {
	Sort                               string
	SortBy                             string
	Keyword                            string
	Filters                            map[string][]string
	QuickFilter                        []string
	Page                               int
	PageSize                           int
	ClauseID                           uint
	IsRegulationNotificationAssessment *bool
	IsRiskBasedAssessment              *bool
}

type ComplianceAssessmentRequirementClausesListDTO struct {
	ID                    uint                                             `json:"id"`
	ClauseNo              string                                           `json:"clause_no"`
	TitleEn               string                                           `json:"title_en"`
	TitleAr               *string                                          `json:"title_ar"`
	ChildClauses          []*ComplianceAssessmentRequirementClausesListDTO `json:"child_clauses"`
	TotalRequirements     int                                              `json:"total_requirements"`
	CompletedRequirements int                                              `json:"completed_requirements"`
}
type ClauseRequirementCountDTO struct {
	ID                    uint `json:"id"`
	TotalRequirements     int  `json:"total_requirements"`
	CompletedRequirements int  `json:"completed_requirements"`
}
