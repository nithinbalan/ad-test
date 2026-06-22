package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type ComplianceAssessmentRequirementListDTO struct {
	ID                                             uint                                                     `json:"id"`
	ClauseID                                       uint                                                     `json:"clause_id"`
	ClauseDescriptionEn                            *string                                                  `json:"clause_description_en"`
	ClauseDescriptionAr                            *string                                                  `json:"clause_description_ar"`
	ComplianceClauseNo                             string                                                   `json:"compliance_clause_no"`
	ComplianceClauseTitleEn                        string                                                   `json:"compliance_clause_title_en"`
	ComplianceClauseTitleAr                        string                                                   `json:"compliance_clause_title_ar"`
	ComplianceRequirementID                        uint                                                     `json:"compliance_requirement_id"`
	ComplianceAssessmentRequirementStatusID        *uint                                                    `json:"compliance_assessment_requirement_status_id"`
	ComplianceRequirementTitleEn                   string                                                   `json:"compliance_requirement_title_en"`
	ComplianceRequirementTitleAr                   string                                                   `json:"compliance_requirement_title_ar"`
	ComplianceRequirementCategoryTitleEn           string                                                   `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr           string                                                   `json:"compliance_requirement_category_title_ar"`
	SeverityLevelID                                uint                                                     `json:"severity_level_id"`
	SeverityLevelTitleEn                           string                                                   `json:"severity_level_title_en"`
	SeverityLevelTitleAr                           string                                                   `json:"severity_level_title_ar"`
	SeverityLevelColorCode                         string                                                   `json:"severity_level_color_code"`
	SeverityLevelType                              string                                                   `json:"severity_level_type"`
	ComplianceAssessmentRequirementStatusTitleEn   string                                                   `json:"compliance_assessment_requirement_status_title_en"`
	ComplianceAssessmentRequirementStatusTitleAr   string                                                   `json:"compliance_assessment_requirement_status_title_ar"`
	ComplianceAssessmentRequirementStatusColorCode string                                                   `json:"compliance_assessment_requirement_status_color_code"`
	ComplianceAssessmentRequirementStatusType      string                                                   `json:"compliance_assessment_requirement_status_type"`
	CreatedAt                                      string                                                   `json:"created_at"`
	UpdatedAt                                      *string                                                  `json:"updated_at"`
	IsAllCompleted                                 bool                                                     `json:"is_all_completed"`
	AverageResidualRiskScore                       *float64                                                 `json:"average_residual_risk_score"`
	AverageResidualRiskRatingID                    *uint                                                    `json:"average_residual_risk_rating_id"`
	AverageResidualRiskRatingTitleEn               *string                                                  `json:"average_residual_risk_rating_title_en"`
	AverageResidualRiskRatingTitleAr               *string                                                  `json:"average_residual_risk_rating_title_ar"`
	AverageResidualRiskRatingColorCode             *string                                                  `json:"average_residual_risk_rating_color_code"`
	AverageResidualRiskRatingType                  *string                                                  `json:"average_residual_risk_rating_type"`
	ResponsibleUnits                               []*ComplianceAssessmentRequirementResponsibleUnitListDTO `json:"responsible_units" gorm:"-"`
}

type ComplianceAssessmentRequirementResponsibleUnitListDTO struct {
	ID                                 uint                              `json:"id"`
	AssessmentRequirementID            uint                              `json:"assessment_requirement_id"`
	ResponsibleUnitID                  uint                              `json:"responsible_unit_id"`
	IsRegulationNotificationAssessment *bool                             `json:"is_regulation_notification_assessment"`
	IsRiskBasedAssessment              *bool                             `json:"is_risk_based_assessment"`
	ResponsibleUnit                    *settingsDTOs.OrganizationUnitDTO `json:"responsible_unit"`
	ResponsibleUserID                  *uint                             `json:"responsible_user_id"`
	ResponsibleUser                    *organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
	IsAcknowledged                     bool                              `json:"is_acknowledged"`
	AcknowledgeUserID                  *uint                             `json:"acknowledge_user_id"`
	AcknowledgeComment                 *string                           `json:"acknowledge_comment"`
	AcknowledgeUser                    *organizationDTOs.UserMinimalDTO  `json:"acknowledge_user"`
}

type ComplianceAssessmentRequirementResponsibleUnitListRowDTO struct {
	ID                                 uint    `json:"id"`
	AssessmentRequirementID            uint    `json:"assessment_requirement_id"`
	ResponsibleUnitID                  uint    `json:"responsible_unit_id"`
	IsRegulationNotificationAssessment *bool   `json:"is_regulation_notification_assessment"`
	IsRiskBasedAssessment              *bool   `json:"is_risk_based_assessment"`
	TitleEn                            string  `json:"title_en"`
	TitleAr                            string  `json:"title_ar"`
	Status                             string  `json:"status"`
	ColorCode                          string  `json:"color_code"`
	ParentID                           *uint   `json:"parent_id"`
	ResponsibleUserID                  *uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn         *string `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn          *string `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr         *string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr          *string `json:"responsible_user_last_name_ar"`
	ResponsibleUserEmail               *string `json:"responsible_user_email"`
	ResponsibleUserInitial             *string `json:"responsible_user_initial"`
	ResponsibleUserColor               *string `json:"responsible_user_color"`
	IsAcknowledged                     bool    `json:"is_acknowledged"`
	AcknowledgeUserID                  *uint   `json:"acknowledge_user_id"`
	AcknowledgeComment                 *string `json:"acknowledge_comment"`
	AcknowledgeUserFirstNameEn         *string `json:"acknowledge_user_first_name_en"`
	AcknowledgeUserLastNameEn          *string `json:"acknowledge_user_last_name_en"`
	AcknowledgeUserFirstNameAr         *string `json:"acknowledge_user_first_name_ar"`
	AcknowledgeUserLastNameAr          *string `json:"acknowledge_user_last_name_ar"`
	AcknowledgeUserEmail               *string `json:"acknowledge_user_email"`
	AcknowledgeUserInitial             *string `json:"acknowledge_user_initial"`
	AcknowledgeUserColor               *string `json:"acknowledge_user_color"`
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

type AcknowledgeRequirementFormDTO struct {
	ResponsibleUnitID uint    `json:"responsible_unit_id" binding:"required"`
	Comment           *string `json:"comment" binding:"omitempty"`
}

type ComplianceAssessmentRequirementUpdateFormDTO struct {
	ResponsibleUnitID          uint    `json:"responsible_unit_id" binding:"required"`
	Score                      int     `json:"score" binding:"omitempty,min=0"`
	Comments                   *string `json:"comments"`
	AssessmentFrameworkValueID uint    `json:"assessment_framework_value_id" binding:"required"`
}

type ComplianceAssessmentRequirementNotApplicableFormDTO struct {
	ResponsibleUnitID uint   `json:"responsible_unit_id" binding:"required"`
	Comments          string `json:"comments" binding:"required"`
	Score             int    `json:"score" binding:"min=0,max=100"`
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
	ID             *uint                            `json:"id,omitempty"`
	Language       string                           `json:"language"`
	URL            string                           `json:"url"`
	Name           string                           `json:"name"`
	Token          string                           `json:"token"`
	OriginalName   string                           `json:"original_name"`
	Type           string                           `json:"type"`
	Size           int                              `json:"size"`
	ApprovedByUser *organizationDTOs.UserMinimalDTO `json:"approved_by"`
	RejectedByUser *organizationDTOs.UserMinimalDTO `json:"rejected_by"`
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

type ComplianceAssessmentRequirementDocumentUpdateFormDTO struct {
	ID     uint   `json:"id" binding:"required"`
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}
