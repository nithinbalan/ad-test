package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type ComplianceAssessmentListDTO struct {
	ID                                  uint       `json:"id"`
	ComplianceRegisterID                uint       `json:"compliance_register_id"`
	ComplianceRegisterTitleEn           string     `json:"compliance_register_title_en"`
	ComplianceRegisterTitleAr           string     `json:"compliance_register_title_ar"`
	IsRegulationNotificationAssessment  bool       `json:"is_regulation_notification_assessment"`
	IsRiskBasedAssessment               bool       `json:"is_risk_based_assessment"`
	Description                         *string    `json:"description"`
	ReviewUserID                        uint       `json:"review_user_id"`
	ReviewUserFirstNameEn               string     `json:"review_user_first_name_en"`
	ReviewUserFirstNameAr               string     `json:"review_user_first_name_ar"`
	ReviewUserLastNameEn                string     `json:"review_user_last_name_en"`
	ReviewUserLastNameAr                string     `json:"review_user_last_name_ar"`
	ReviewUserInitial                   string     `json:"review_user_initial"`
	ReviewUserColor                     string     `json:"review_user_color"`
	ReviewUserEmail                     string     `json:"review_user_email"`
	ReviewUserProfileImagePath          string     `json:"review_user_profile_image_path"`
	ComplianceAssessmentStatusTitleEn   string     `json:"compliance_assessment_status_title_en"`
	ComplianceAssessmentStatusTitleAr   string     `json:"compliance_assessment_status_title_ar"`
	ComplianceAssessmentStatusColorCode string     `json:"compliance_assessment_status_color_code"`
	ComplianceAssessmentStatusType      string     `json:"compliance_assessment_status_type"`
	CreatedAt                           time.Time  `json:"created_at"`
	UpdatedAt                           *time.Time `json:"updated_at"`
	CreatedUserID                       uint       `json:"created_user_id"`
	CreatedUserFirstNameEn              string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr              string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn               string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr               string     `json:"created_user_last_name_ar"`
	CreatedUserInitial                  string     `json:"created_user_initial"`
	CreatedUserColor                    string     `json:"created_user_color"`
	CreatedUserEmail                    string     `json:"created_user_email"`
	CreatedUserProfileImagePath         string     `json:"created_user_profile_image_path"`
	UpdatedUserID                       *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn              *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr              *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn               *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr               *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                  *string    `json:"updated_user_initial"`
	UpdatedUserColor                    *string    `json:"updated_user_color"`
	UpdatedUserEmail                    *string    `json:"updated_user_email"`
	UpdatedUserProfileImagePath         *string    `json:"updated_user_profile_image_path"`
	AverageResidualRiskScore            *float64   `json:"average_residual_risk_score"`
	AverageResidualRiskRatingID         *uint      `json:"average_residual_risk_rating_id"`
	AverageResidualRiskRatingTitleEn    *string    `json:"average_residual_risk_rating_title_en"`
	AverageResidualRiskRatingTitleAr    *string    `json:"average_residual_risk_rating_title_ar"`
	AverageResidualRiskRatingColorCode  *string    `json:"average_residual_risk_rating_color_code"`
	AverageResidualRiskRatingType       *string    `json:"average_residual_risk_rating_type"`
}

type ComplianceAssessmentListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type ComplianceAssessmentListResponseDTO struct {
	Items      []*ComplianceAssessmentListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO   `json:"pagination"`
}

type ComplianceAssessmentFormDTO struct {
	ComplianceRegisterID uint `json:"compliance_register_id" binding:"required"`
	// OrganizationUnitIDs  []*uint `json:"organization_unit_ids" binding:"required"`
	ReviewUserID uint    `json:"review_user_id" binding:"required"`
	Description  *string `json:"description"`
}

type ComplianceAssessmentUpdateFormDTO struct {
	ComplianceRegisterID uint    `json:"compliance_register_id" binding:"required"`
	ReviewUserID         uint    `json:"review_user_id" binding:"required"`
	Description          *string `json:"description"`
}

type ComplianceAssessmentDetailDTO struct {
	ID                   int                          `json:"id"`
	ComplianceRegisterID *uint                        `json:"compliance_register_id"`
	ComplianceRegister   *ComplianceRegisterDetailDTO `json:"compliance_register"`
	// OrganizationUnitID uint                         `json:"organization_unit_id"`
	ReviewUserID      uint `json:"review_user_id"`
	SubmittedByUserID uint `json:"submitted_by_user_id"`
	// OrganizationUnit           settingsDTOs.OrganizationUnitDTO      `json:"organization_unit"`
	ReviewUser                 *organizationDTOs.UserMinimalDTO      `json:"review_user"`
	Description                *string                               `json:"description"`
	ComplianceAssessmentStatus ComplianceAssessmentStatusDTO         `json:"compliance_assessment_status"`
	AverageResidualRiskScore   *float64                              `json:"average_residual_risk_score"`
	AverageResidualRiskRating  *ComplianceAssessmentRatingDTO        `json:"average_residual_risk_rating"`
	SubmittedByUser            *organizationDTOs.UserMinimalDTO      `json:"submitted_by_user"`
	SubmittedAt                *time.Time                            `json:"submitted_at"`
	ReviewedByUser             *organizationDTOs.UserMinimalDTO      `json:"reviewed_by_user"`
	CreatedAt                  time.Time                             `json:"created_at"`
	UpdatedAt                  *time.Time                            `json:"updated_at"`
	CreatedByUser              organizationDTOs.UserMinimalDTO       `json:"created_by_user"`
	UpdatedByUser              *organizationDTOs.UserMinimalDTO      `json:"updated_by_user"`
	AssessmentRequirements     []*ComplianceAssessmentRequirementDTO `json:"compliance_requirements"`
	IsCompleted                bool                                  `json:"is_completed"`
	IsApproved                 bool                                  `json:"is_approved"`
	AssessmentControls         []*ComplianceAssessmentControlDTO     `json:"compliance_controls"` // ✅ new field
}

type ComplianceAssessmentRequirementDTO struct {
	ID                             *uint                                                `json:"id"`
	ComplianceAssessmentID         uint                                                 `json:"compliance_assessment_id"`
	ComplianceRequirementID        uint                                                 `json:"compliance_requirement_id"`
	ComplianceRequirement          ComplianceRequirementDetailDTO                       `json:"compliance_requirement"`
	IsCompleted                    bool                                                 `json:"is_completed"`
	Score                          *int                                                 `json:"score"`
	Comments                       *string                                              `json:"comments"`
	ComplianceStatusID             *uint                                                `json:"compliance_status_id"`
	ComplianceStatus               *ComplianceStatusDTO                                 `json:"compliance_status"`
	ResponsibleUnits               []*ComplianceAssessmentRequirementResponsibleUnitDTO `json:"responsible_units"`
	InherentRiskAnalysis           []*ComplianceRequirementInherentRiskDTO              `json:"inherent_risk_analysis"`
	AverageResidualRatingRiskScore *float64                                             `json:"average_residual_rating_risk_score"`
	AverageResidualRiskRating      *ComplianceAssessmentRatingDTO                       `json:"average_residual_risk_rating"`
}
type ComplianceAssessmentRequirementResponsibleUnitDTO struct {
	ID                                    uint                                               `json:"id"`
	ResponsibleUnitID                     uint                                               `json:"responsible_unit_id"`
	ResponsibleUnit                       settingsDTOs.OrganizationUnitDTO                   `json:"responsible_unit"`
	ComplianceAssessmentFrameworkValue    *ComplianceAssessmentFrameworkValueDetailDTO       `json:"compliance_assessment_framework_value"`
	ComplianceAssessmentRequirement       *ComplianceAssessmentRequirementDetailDTO          `json:"compliance_assessment_requirement"`
	RequiredDocumentTitles                []*ComplianceAssessmentRequirementDocumentTitleDTO `json:"required_document_titles"`
	ActionPlans                           []*ComplianceActionPlanDetailDTO                   `json:"action_plans"`
	RequirementComments                   []*ComplianceAssessmentRequirementCommentDetailDTO `json:"requirement_comments"`
	ResidualRiskAnalysis                  []*ComplianceAssessmentRequirementResidualRiskDTO  `json:"residual_risk_analysis"`
	ResidualControls                      []*ComplianceAssessmentResidualControlDTO          `json:"residual_controls"`
	CompletedByUser                       *organizationDTOs.UserMinimalDTO                   `json:"completed_by_user"`
	ComplianceStatusID                    *uint                                              `json:"compliance_status_id"`
	IsCompleted                           bool                                               `json:"is_completed"`
	Score                                 *int                                               `json:"score"`
	Comments                              *string                                            `json:"comments"`
	ComplianceAssessmentRequirementStatus *ComplianceAssessmentRequirementStatusDTO          `json:"compliance_assessment_requirement_status"`
	CompletedAt                           *time.Time                                         `json:"completed_at"` // ⬅️ new field
	ResidualRiskScore                     *float64                                           `json:"residual_risk_score"`
	ResidualRiskRating                    *ComplianceAssessmentRatingDTO                     `json:"residual_risk_rating"`
}

type AssessmentCompletionSummaryDTO struct {
	Total        int `json:"total"`
	Completed    int `json:"completed"`
	NotCompleted int `json:"not_completed"`
}
type AssessmentCompletionRollupDTO struct {
	Total                int     `json:"total"`
	Completed            int     `json:"completed"`
	NotCompleted         int     `json:"not_completed"`
	CompletionPercentage float64 `json:"completion_percentage"` // 0–100
}

type AssessmentCompletionSummaryResponsibleUnitDTO struct {
	ResponsibleUnitID      uint   `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr string `json:"responsible_unit_title_ar"`
	Total                  int    `json:"total"`
	Completed              int    `json:"completed"`
	NotCompleted           int    `json:"not_completed"`
}
type AssessmentCompletionSummaryGroupedByResponsibleUnitDTO struct {
	ResponsibleUnitID      uint                           `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn string                         `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr string                         `json:"responsible_unit_title_ar"`
	Summary                AssessmentCompletionSummaryDTO `json:"summary"`
}

// type ComplianceStatusSummaryDTO struct {
// 	TotalRequirements int            `json:"total_requirements"`
// 	StatusCounts      map[string]int `json:"status_counts"`
// 	PendingCount      int            `json:"pending_count"`
// }

type ComplianceStatusCountDTO struct {
	ComplianceStatusID uint   `json:"compliance_status_id"`
	TitleEn            string `json:"title_en"`
	TitleAr            string `json:"title_ar"`
	Type               string `json:"type"`
	Count              int    `json:"count"`
	ColorCode          string `json:"color_code"`
}

type ComplianceScoreSummaryDTO struct {
	TotalRequirements int     `json:"total_requirements_assessed"`
	AverageScore      float64 `json:"average_score"`
}

type ResponsibleUnitComplianceScoreDTO struct {
	ComplianceRegisterID              uint    `json:"compliance_register_id"`
	ComplianceAssessmentRequirementID uint    `json:"compliance_assessment_requirement_id"`
	TotalRequirements                 int     `json:"total_requirements"`
	ScoredRequirements                int     `json:"scored_requirements"`
	TotalScore                        float64 `json:"total_score"`
	MaxPossibleScore                  float64 `json:"max_possible_score"`
	CompliancePercentage              float64 `json:"compliance_percentage"`
	Score                             float64 `json:"score"`
	FrameworkType                     string  `json:"framework_type"`
	MappedFrameworkValueID            *uint   `json:"mapped_framework_value_id"`
	MappedComplianceStatusID          *uint   `json:"mapped_compliance_status_id"`
	ComplianceClauseID                uint    `json:"compliance_clause_id"`
}

type CompletionStatusResponsibleUnitDTO struct {
	ResponsibleUnitID      uint   `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr string `json:"responsible_unit_title_ar"`
	TitleEn                string `json:"title_en"`
	TitleAr                string `json:"title_ar"`
	ColorCode              string `json:"color_code"`
	Count                  int    `json:"count"`
}

// Grouped output for API
type CompletionStatusGroupedByResponsibleUnitDTO struct {
	ResponsibleUnitID      uint                       `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn string                     `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr string                     `json:"responsible_unit_title_ar"`
	StatusCount            []ComplianceStatusCountDTO `json:"status_count"`
}

type ComplianceStatusByRequirementRowDTO struct {
	ComplianceAssessmentRequirementID uint   `json:"compliance_assessment_requirement_id"`
	ComplianceStatusID                uint   `json:"compliance_status_id"`
	TitleEn                           string `json:"title_en"`
	TitleAr                           string `json:"title_ar"`
	ColorCode                         string `json:"color_code"`
	Count                             int    `json:"count"`
}

// grouped output per requirement
type CompletionStatusGroupedByRequirementDTO struct {
	ComplianceAssessmentRequirementID uint                       `json:"compliance_assessment_requirement_id"`
	StatusCount                       []ComplianceStatusCountDTO `json:"status_count"`
}

type ComplianceAssessmentApprovalStatusDTO struct {
	AssessmentID uint `json:"assessment_id"`
	IsApproved   bool `json:"is_approved"`
}
type ClauseStatusAggregateDTO struct {
	TotalCount          int
	CompliantCount      int
	PartialCount        int
	NotImplementedCount int
}
type ComplianceRequirementScoreStatusDTO struct {
	ComplianceStatusID *uint    `json:"compliance_status_id"`
	Score              *float64 `json:"score"`
}

type ComplianceAssessmentControlDTO struct {
	ID                                      *uint                                        `json:"id"`
	ComplianceAssessmentID                  uint                                         `json:"compliance_assessment_id"`
	ComplianceControlID                     uint                                         `json:"compliance_control_id"`
	ComplianceControl                       ComplianceClauseControlDTO                   `json:"compliance_control"` // nested details
	CompletedByUserID                       *uint                                        `json:"completed_by_user_id"`
	CompletedByUser                         *organizationDTOs.UserMinimalDTO             `json:"completed_by_user"`
	CompletedAt                             *time.Time                                   `json:"completed_at"`
	ComplianceAssessmentFrameworkValueID    *uint                                        `json:"compliance_assessment_framework_value_id"`
	ComplianceAssessmentFrameworkValue      *ComplianceAssessmentFrameworkValueDetailDTO `json:"framework_value"`
	Comments                                *string                                      `json:"comments"`
	IsCompleted                             bool                                         `json:"is_completed"`
	Score                                   *float64                                     `json:"score"`
	ComplianceAssessmentRequirementStatusID *uint                                        `json:"compliance_assessment_requirement_status_id"`
	ComplianceAssessmentRequirementStatus   *ComplianceAssessmentRequirementStatusDTO    `json:"requirement_status"`
	CreatedAt                               time.Time                                    `json:"created_at"`
	UpdatedAt                               *time.Time                                   `json:"updated_at"`
	CreatedByUser                           *organizationDTOs.UserMinimalDTO             `json:"created_by_user"`
	UpdatedByUser                           *organizationDTOs.UserMinimalDTO             `json:"updated_by_user"`
}

type ControlComplianceScoreDTO struct {
	ComplianceRegisterID   uint   `json:"compliance_register_id"`
	ComplianceAssessmentID uint   `json:"compliance_assessment_id"`
	ComplianceControlID    uint   `json:"compliance_control_id"`
	ComplianceClauseID     uint   `json:"compliance_clause_id"`
	FrameworkType          string `json:"framework_type"`

	// scoring
	TotalControls        int     `json:"total_controls"`
	ScoredControls       int     `json:"scored_controls"`
	TotalScore           float64 `json:"total_score"`
	MaxPossibleScore     float64 `json:"max_possible_score"`
	ComplianceScore      float64 `json:"compliance_score"`
	CompliancePercentage float64 `json:"compliance_percentage"`
	Score                float64 `json:"score"`

	// mappings
	MappedFrameworkValueID   *uint `json:"mapped_framework_value_id"`
	MappedComplianceStatusID *uint `json:"mapped_compliance_status_id"`
}
type ClauseComplianceScoreDTO struct {
	ComplianceClauseID       *uint   `json:"compliance_clause_id"`
	ControlCount             int     `json:"control_count"`
	SumPct                   float64 `json:"-"`
	AveragePct               float64 `json:"average_pct"`
	MappedFrameworkValueID   *uint   `json:"mapped_framework_value_id,omitempty"`
	MappedComplianceStatusID *uint   `json:"mapped_compliance_status_id,omitempty"`
}
type ControlSyncUpdateDTO struct {
	ControlID                            uint
	ComplianceControlID                  uint
	ComplianceAssessmentFrameworkValueID *uint
	CompletedByUserID                    *uint
	CompletedAt                          *time.Time
	ComplianceStatusID                   *uint
}

type ComplianceAssessmentScoreSummaryQueryParamsDTO struct {
	Filters map[string][]string `json:"filters"`
}
