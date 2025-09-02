package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ComplianceAssessmentControlListDTO struct {
	ID                                             uint     `json:"id"`
	ClauseID                                       uint     `json:"clause_id"`
	ComplianceClauseNo                             string   `json:"compliance_clause_no"`
	ComplianceClauseTitleEn                        string   `json:"compliance_clause_title_en"`
	ComplianceClauseTitleAr                        string   `json:"compliance_clause_title_ar"`
	ComplianceControlID                            uint     `json:"compliance_control_id"`
	ControlTitleEn                                 string   `json:"control_title_en"`
	ControlTitleAr                                 string   `json:"control_title_ar"`
	ComplianceAssessmentRequirementStatusID        *uint    `json:"compliance_assessment_requirement_status_id"`
	ComplianceAssessmentRequirementStatusTitleEn   string   `json:"compliance_assessment_requirement_status_title_en"`
	ComplianceAssessmentRequirementStatusTitleAr   string   `json:"compliance_assessment_requirement_status_title_ar"`
	ComplianceAssessmentRequirementStatusColorCode string   `json:"compliance_assessment_requirement_status_color_code"`
	ComplianceAssessmentRequirementStatusType      string   `json:"compliance_assessment_requirement_status_type"`
	IsCompleted                                    bool     `json:"is_completed"`
	Score                                          *float64 `json:"score"`
	OwnerID                                        *uint    `json:"owner_id"`
	OwnerFirstNameEn                               *string  `json:"owner_first_name_en"`
	OwnerFirstNameAr                               *string  `json:"owner_first_name_ar"`
	OwnerLastNameEn                                *string  `json:"owner_last_name_en"`
	OwnerLastNameAr                                *string  `json:"owner_last_name_ar"`
	OwnerInitial                                   string   `json:"owner_initial"`
	OwnerColor                                     string   `json:"owner_color"`
	OwnerEmail                                     string   `json:"owner_email"`
	OwnerProfileImagePath                          string   `json:"owner_profile_image_path"`
	CreatedAt                                      string   `json:"created_at"`
	UpdatedAt                                      *string  `json:"updated_at"`
}

type ComplianceAssessmentControlListResponseDTO struct {
	Items      []*ComplianceAssessmentControlListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO          `json:"pagination"`
}

type ComplianceAssessmentControlDetailDTO struct {
	ID                                    int                                            `json:"id"`
	ComplianceAssessmentID                uint                                           `json:"compliance_assessment_id"`
	ComplianceControlID                   uint                                           `json:"compliance_control_id"`
	ComplianceControl                     ComplianceControlDTO                           `json:"compliance_control"`
	Comments                              *string                                        `json:"comments"`
	IsCompleted                           bool                                           `json:"is_completed"`
	Score                                 *float64                                       `json:"score"`
	ComplianceAssessmentRequirementStatus *ComplianceAssessmentRequirementStatusDTO      `json:"compliance_assessment_requirement_status"`
	ComplianceAssessmentFrameworkValueID  *uint                                          `json:"compliance_assessment_framework_value_id"`
	ComplianceAssessmentFrameworkValue    *ComplianceAssessmentFrameworkValueDetailDTO   `json:"compliance_assessment_framework_value"`
	CompletedByUser                       *organizationDTOs.UserMinimalDTO               `json:"completed_by_user"`
	CreatedAt                             time.Time                                      `json:"created_at"`
	UpdatedAt                             *time.Time                                     `json:"updated_at"`
	CreatedByUser                         organizationDTOs.UserMinimalDTO                `json:"created_by_user"`
	UpdatedByUser                         *organizationDTOs.UserMinimalDTO               `json:"updated_by_user"`
	ControlComments                       []*ComplianceAssessmentControlCommentDetailDTO `json:"control_comments"`
	ActionPlans                           []*ComplianceActionPlanDetailDTO               `json:"action_plans"`
}

type ComplianceAssessmentControlUpdateFormDTO struct {
	Score                      *float64 `json:"score"`
	Comments                   *string  `json:"comments"`
	AssessmentFrameworkValueID uint     `json:"assessment_framework_value_id" binding:"required"`
}

type ComplianceAssessmentControlClausesListDTO struct {
	ID                uint                                         `json:"id"`
	ClauseNo          string                                       `json:"clause_no"`
	TitleEn           string                                       `json:"title_en"`
	TitleAr           *string                                      `json:"title_ar"`
	ChildClauses      []*ComplianceAssessmentControlClausesListDTO `json:"child_clauses"`
	TotalControls     int                                          `json:"total_controls"`
	CompletedControls int                                          `json:"completed_controls"`
}

type ClauseControlCountDTO struct {
	ID                uint `json:"id"`
	TotalControls     int  `json:"total_controls"`
	CompletedControls int  `json:"completed_controls"`
}

type ComplianceAssessmentControlListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Filters     map[string][]string
	QuickFilter []string
	Page        int
	PageSize    int
	ClauseID    uint
}
