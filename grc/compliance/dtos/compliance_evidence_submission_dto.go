package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTO "jethur/grc/settings/dtos"
	"time"
)

type ComplianceEvidenceSubmissionListDTO struct {
	ID                                             uint    `json:"id"`
	IsCompleted                                    bool    `json:"is_completed"`
	IsAcknowledged                                 bool    `json:"is_acknowledged"`
	ComplianceAssessmentID                         uint    `json:"compliance_assessment_id"`
	ClauseID                                       uint    `json:"clause_id"`
	ComplianceClauseNo                             string  `json:"compliance_clause_no"`
	ComplianceClauseTitleEn                        string  `json:"compliance_clause_title_en"`
	ComplianceClauseTitleAr                        string  `json:"compliance_clause_title_ar"`
	ComplianceRequirementID                        *uint   `json:"compliance_requirement_id"`
	ComplianceRequirementTitleEn                   *string `json:"compliance_requirement_title_en"`
	ComplianceRequirementTitleAr                   *string `json:"compliance_requirement_title_ar"`
	ComplianceRegisterTitleEn                      string  `json:"compliance_register_title_en"`
	ComplianceRegisterTitleAr                      string  `json:"compliance_register_title_ar"`
	ComplianceAssessmentRequirementStatusID        *uint   `json:"compliance_assessment_requirement_status_id"`
	ComplianceAssessmentRequirementStatusTitleEn   *string `json:"compliance_assessment_requirement_status_title_en"`
	ComplianceAssessmentRequirementStatusTitleAr   *string `json:"compliance_assessment_requirement_status_title_ar"`
	ComplianceAssessmentRequirementStatusColorCode *string `json:"compliance_assessment_requirement_status_color_code"`
	ComplianceAssessmentRequirementStatusType      *string `json:"compliance_assessment_requirement_status_type"`
	ResponsibleUnitID                              uint    `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn                         string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr                         string  `json:"responsible_unit_title_ar"`
	ResponsibleUserID                              *uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn                     *string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr                     *string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn                      *string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr                      *string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial                         *string `json:"responsible_user_initial"`
	ResponsibleUserColor                           *string `json:"responsible_user_color"`
	ResponsibleUserEmail                           *string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath                *string `json:"responsible_user_profile_image_path"`
	CompletedUserID                                *uint   `json:"completed_user_id"`
	CompletedUserFirstNameEn                       *string `json:"completed_user_first_name_en"`
	CompletedUserFirstNameAr                       *string `json:"completed_user_first_name_ar"`
	CompletedUserLastNameEn                        *string `json:"completed_user_last_name_en"`
	CompletedUserLastNameAr                        *string `json:"completed_user_last_name_ar"`
	CompletedUserInitial                           *string `json:"completed_user_initial"`
	CompletedUserColor                             *string `json:"completed_user_color"`
	CompletedUserEmail                             *string `json:"completed_user_email"`
	CompletedUserProfileImagePath                  *string `json:"completed_user_profile_image_path"`
	IsRiskBasedAssessment                          *bool   `json:"is_risk_based_assessment"`
	IsRegulationNotificationAssessment             *bool   `json:"is_regulation_notification_assessment"`
	CreatedAt                                      string  `json:"created_at"`
	UpdatedAt                                      *string `json:"updated_at"`
	CompletedAt                                    *string `json:"completed_at"`
}

type ComplianceEvidenceSubmissionListResponseDTO struct {
	Items      []*ComplianceEvidenceSubmissionListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO           `json:"pagination"`
}

type ComplianceEvidenceSubmissionListQueryParamsDTO struct {
	Sort          string
	SortBy        string
	Keyword       string
	Filters       map[string][]string
	QuickFilter   []string
	Page          int
	PageSize      int
	IsSuperAdmin  bool
	IsPassiveUser bool
	Export        bool
}

type ComplianceEvidenceSubmissionRequirementDetailDTO struct {
	ID                                      uint                                               `json:"id"`
	ComplianceAssessmentRequirementID       uint                                               `json:"compliance_assessment_requirement_id"`
	ResponsibleUnitID                       uint                                               `json:"responsible_unit_id"`
	ResponsibleUnit                         *settingsDTO.OrganizationUnitDTO                   `json:"responsible_unit"`
	ResponsibleUserID                       *uint                                              `json:"responsible_user_id"`
	ResponsibleUser                         *organizationDTOs.UserMinimalDTO                   `json:"responsible_user"`
	IsCompleted                             bool                                               `json:"is_completed"`
	IsAcknowledged                          bool                                               `json:"is_acknowledged"`
	IsNotApplicable                         bool                                               `json:"is_not_applicable"`
	IsRiskBasedAssessment                   *bool                                              `json:"is_risk_based_assessment"`
	IsRegulationNotificationAssessment      *bool                                              `json:"is_regulation_notification_assessment"`
	Score                                   *int                                               `json:"score"`
	Comments                                *string                                            `json:"comments"`
	AcknowledgeUserID                       *uint                                              `json:"acknowledge_user_id"`
	AcknowledgeComment                      *string                                            `json:"acknowledge_comment"`
	AcknowledgeUser                         *organizationDTOs.UserMinimalDTO                   `json:"acknowledge_user"`
	ComplianceAssessmentRequirementStatusID *uint                                              `json:"compliance_assessment_requirement_status_id"`
	ComplianceAssessmentRequirementStatus   *ComplianceAssessmentRequirementStatusDTO          `json:"compliance_assessment_requirement_status"`
	CompletedByUser                         *organizationDTOs.UserMinimalDTO                   `json:"completed_by_user"`
	ResidualRiskScore                       *float64                                           `json:"residual_risk_score"`
	ComplianceStatusID                      *uint                                              `json:"compliance_status_id"`
	ComplianceStatus                        *ComplianceStatusDTO                               `json:"compliance_status"`
	ComplianceRequirement                   ComplianceRequirementDetailDTO                     `json:"compliance_requirement"`
	ComplianceAssessment                    *ComplianceAssessmentDetailDTO                     `json:"compliance_assessment"`
	ComplianceRegister                      ComplianceRegisterDetailDTO                        `json:"compliance_register"`
	RequiredDocumentTitles                  []*ComplianceAssessmentRequirementDocumentTitleDTO `json:"required_document_titles"`
	ActionPlans                             []*ComplianceActionPlanDetailDTO                   `json:"action_plans"`
	RequirementComments                     []*ComplianceAssessmentRequirementCommentDetailDTO `json:"requirement_comments"`
	CreatedAt                               time.Time                                          `json:"created_at"`
	UpdatedAt                               *time.Time                                         `json:"updated_at"`
	CompletedAt                             *time.Time                                         `json:"completed_at"`
}
