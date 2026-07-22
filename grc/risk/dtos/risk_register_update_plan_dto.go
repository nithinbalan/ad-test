package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
)

type RiskRegisterUpdatePlanListDTO struct {
	ID                             uint    `json:"id"`
	DivisionID                     uint    `json:"division_id"`
	DivisionTitleEn                *string `json:"division_title_en"`
	DivisionTitleAr                *string `json:"division_title_ar"`
	BusinessUnitID                 uint    `json:"business_unit_id"`
	BusinessUnitTitleEn            *string `json:"business_unit_title_en"`
	BusinessUnitTitleAr            *string `json:"business_unit_title_ar"`
	RiskRegisterUpdatePlanStatusID uint    `json:"risk_register_update_plan_status_id"`
	StatusTitleEn                  string  `json:"status_title_en"`
	StatusTitleAr                  string  `json:"status_title_ar"`
	StatusType                     string  `json:"status_type"`
	StatusColorCode                string  `json:"status_color_code"`
	NumberOfOperatingEntities      *int    `json:"number_of_operating_entities"`
	Workdays                       *int    `json:"workdays"`
	RiskChampionID                 uint    `json:"risk_champion_id"`
	RiskChampionFirstNameEn        string  `json:"risk_champion_first_name_en"`
	RiskChampionFirstNameAr        string  `json:"risk_champion_first_name_ar"`
	RiskChampionLastNameEn         string  `json:"risk_champion_last_name_en"`
	RiskChampionLastNameAr         string  `json:"risk_champion_last_name_ar"`
	RiskChampionEmail              string  `json:"risk_champion_email"`
	RiskChampionInitial            string  `json:"risk_champion_initial"`
	RiskChampionColor              string  `json:"risk_champion_color"`
	RiskChampionProfileImagePath   string  `json:"risk_champion_profile_image_path"`
	StartDate                      string  `json:"start_date"`
	EndDate                        string  `json:"end_date"`
	Notes                          *string `json:"notes"`
	TotalPercentageCompleted       *int    `json:"total_percentage_completed"`
	CreatedAt                      string  `json:"created_at"`
	UpdatedAt                      *string `json:"updated_at"`
	CreatedUserID                  uint    `json:"created_user_id"`
	CreatedUserFirstNameEn         string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr         string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn          string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr          string  `json:"created_user_last_name_ar"`
	CreatedUserInitial             string  `json:"created_user_initial"`
	CreatedUserColor               string  `json:"created_user_color"`
	CreatedUserEmail               string  `json:"created_user_email"`
	CreatedUserProfileImagePath    string  `json:"created_user_profile_image_path"`
	UpdatedUserID                  *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn         string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr         string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn          string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr          string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial             string  `json:"updated_user_initial"`
	UpdatedUserColor               string  `json:"updated_user_color"`
	UpdatedUserEmail               string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath    string  `json:"updated_user_profile_image_path"`
}

type RiskRegisterUpdatePlanListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type RiskRegisterUpdatePlanListResponseDTO struct {
	Items      []*RiskRegisterUpdatePlanListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}

type RiskRegisterUpdatePlanDetailDTO struct {
	ID                             uint                                         `json:"id"`
	DivisionID                     uint                                         `json:"division_id"`
	Division                       *settingsDTOs.OrganizationUnitDTO            `json:"division"`
	BusinessUnitID                 uint                                         `json:"business_unit_id"`
	BusinessUnit                   *settingsDTOs.OrganizationUnitDTO            `json:"business_unit"`
	RiskRegisterUpdatePlanStatusID uint                                         `json:"risk_register_update_plan_status_id"`
	RiskRegisterUpdatePlanStatus   *RiskRegisterUpdatePlanStatusDTO             `json:"risk_register_update_plan_status"`
	NumberOfOperatingEntities      *int                                         `json:"number_of_operating_entities"`
	Workdays                       *int                                         `json:"workdays"`
	RiskChampionID                 uint                                         `json:"risk_champion_id"`
	RiskChampion                   *organizationDTOs.UserMinimalDTO             `json:"risk_champion"`
	StartDateString                string                                       `json:"start_date"`
	EndDateString                  string                                       `json:"end_date"`
	Notes                          *string                                      `json:"notes"`
	TotalPercentageCompleted       *int                                         `json:"total_percentage_completed"`
	IsBusinessContextUpdated       bool                                         `json:"is_business_context_updated"`
	TeamMembers                    []*RiskRegisterUpdatePlanTeamMemberDetailDTO `json:"team_members"`
	Timelines                      []*RiskUpdatePlanTimelineDetailDTO           `json:"timelines"`
	CreatedAt                      time.Time                                    `json:"created_at"`
	UpdatedAt                      *time.Time                                   `json:"updated_at"`
	CreatedByUser                  *organizationDTOs.UserMinimalDTO             `json:"created_by"`
	UpdatedByUser                  *organizationDTOs.UserMinimalDTO             `json:"updated_by"`
	SubmittedByUser                *organizationDTOs.UserMinimalDTO             `json:"submitted_by"`
	SubmittedAt                    *time.Time                                   `json:"submitted_at"`
	ApprovedByUser                 *organizationDTOs.UserMinimalDTO             `json:"approved_by"`
	ApprovedAt                     *time.Time                                   `json:"approved_at"`
	NextReviewLevel                *int                                         `json:"next_review_level"`
	WorkflowStatus                 *masterDTOs.WorkflowStatusDTO                `json:"workflow_status"`
}

type RiskRegisterUpdatePlanTeamMemberDetailDTO struct {
	ID        uint                             `json:"id"`
	UserID    uint                             `json:"user_id"`
	User      *organizationDTOs.UserMinimalDTO `json:"user"`
	CreatedAt time.Time                        `json:"created_at"`
	UpdatedAt *time.Time                       `json:"updated_at"`
}

type RiskUpdatePlanTimelineDetailDTO struct {
	ID                           uint                           `json:"id"`
	RiskUpdatePlanID             uint                           `json:"risk_update_plan_id"`
	RiskUpdatePlanTimelineTypeID uint                           `json:"risk_update_plan_timeline_type_id"`
	RiskUpdatePlanTimelineType   *RiskUpdatePlanTimelineTypeDTO `json:"risk_update_plan_timeline_type"`
	IsBusinessContextUpdated     bool                           `json:"is_business_context_updated"`
	IsNotApplicable              bool                           `json:"is_not_applicable"`
	IsCompleted                  bool                           `json:"is_completed"`
	Justification                *string                        `json:"justification"`
	CreatedAt                    time.Time                      `json:"created_at"`
	UpdatedAt                    *time.Time                     `json:"updated_at"`
	CreatedBYID                  uint                           `json:"created_by_id"`
	UpdatedBYID                  *uint                          `json:"updated_by_id"`
	// Relationships
	ErmDeptDocuments                []*RiskUpdatePlanErmDeptDocumentDTO                `json:"erm_dept_documents"`
	IntroductoryMeetingDocuments    []*RiskUpdatePlanIntroductoryMeetingDocumentDTO    `json:"introductory_meeting_documents"`
	BusinessContexts                []*RiskUpdatePlanBusinessContextDTO                `json:"business_contexts"`
	Risks                           []*RiskUpdatePlanRiskDTO                           `json:"risks"`
	Controls                        []*RiskUpdatePlanControlDTO                        `json:"controls"`
	Processes                       []*RiskUpdatePlanProcessDTO                        `json:"processes"`
	LeadershipAssessments           []*RiskUpdatePlanLeadershipAssessmentDTO           `json:"leadership_assessments"`
	LeadershipDocuments             []*RiskUpdatePlanLeadershipDocumentDTO             `json:"leadership_documents"`
	RiskAssessmentWorkshopDocuments []*RiskUpdatePlanRiskAssessmentWorkshopDocumentDTO `json:"risk_assessment_workshop_documents"`
	RiskTeamAssessments             []*RiskUpdatePlanRiskTeamAssessmentDTO             `json:"risk_team_assessments"`
	RiskTeamDocuments               []*RiskUpdatePlanRiskTeamDocumentDTO               `json:"risk_team_documents"`
	RiskWorkshopAssessmentDocuments []*RiskUpdatePlanRiskWorkshopAssessmentDocumentDTO `json:"risk_workshop_assessment_documents"`
	SelfAwarenessTrainings          []*RiskUpdatePlanSelfAwarenessTrainingDTO          `json:"self_awareness_trainings"`
	CreatedByUser                   *organizationDTOs.UserMinimalDTO                   `json:"created_by"`
	UpdatedByUser                   *organizationDTOs.UserMinimalDTO                   `json:"updated_by"`
}

type RiskRegisterUpdatePlanFormDTO struct {
	DivisionID                     uint    `json:"division_id" binding:"required"`
	BusinessUnitID                 uint    `json:"business_unit_id" binding:"required"`
	RiskRegisterUpdatePlanStatusID uint    `json:"risk_register_update_plan_status_id" binding:"required"`
	NumberOfOperatingEntities      *int    `json:"number_of_operating_entities"`
	Workdays                       *int    `json:"workdays"`
	RiskChampionID                 uint    `json:"risk_champion_id" binding:"required"`
	StartDate                      string  `json:"start_date" binding:"required"`
	EndDate                        string  `json:"end_date" binding:"required"`
	Notes                          *string `json:"notes"`
	TeamMemberIds                  []*uint `json:"team_member_ids" binding:"required"`
}

type RiskRegisterUpdatePlanParsedDates struct {
	StartDate time.Time
	EndDate   time.Time
}

// ERM DEPARTMENT RISK
type MapRisksToUpdatePlanDTO struct {
	RiskRegisterIDs []uint `json:"risk_register_ids" binding:"required,min=1"`
}

// List mapped risks
type ERMDeptEvaluationRiskListDTO struct {
	ID                              uint    `json:"id"`
	RiskRegisterUpdatePlanID        uint    `json:"risk_register_update_plan_id"`
	RiskRegisterID                  uint    `json:"risk_register_id"`
	ReferenceNumber                 string  `json:"reference_number"`
	RiskCategoryTitleEn             *string `json:"risk_category_title_en"`
	RiskCategoryTitleAr             *string `json:"risk_category_title_ar"`
	RiskSourceTitleEn               *string `json:"risk_source_title_en"`
	RiskSourceTitleAr               *string `json:"risk_source_title_ar"`
	ResponsibleUnitID               *uint   `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn          *string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          *string `json:"responsible_unit_title_ar"`
	ResponsibleUserID               *uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      *string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      *string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       *string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       *string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string `json:"responsible_user_initial"`
	ResponsibleUserColor            *string `json:"responsible_user_color"`
	ResponsibleUserEmail            *string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath *string `json:"responsible_user_profile_image_path"`
	RiskStatusTitleEn               string  `json:"risk_status_title_en"`
	RiskStatusTitleAr               string  `json:"risk_status_title_ar"`
	RiskResponseStrategyTitleEn     string  `json:"risk_response_strategy_title_en"`
	RiskResponseStrategyTitleAr     string  `json:"risk_response_strategy_title_ar"`
	RiskStatusType                  string  `json:"risk_status_type"`
	RiskStatusColorCode             string  `json:"risk_status_color_code"`
	Title                           string  `json:"title"`
	Description                     *string `json:"description"`
	ReportedAt                      *string `json:"reported_at"`
	ReviewFrequencyTitleEn          *string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr          *string `json:"review_frequency_title_ar"`
	NextReviewDate                  *string `json:"next_review_date"`
	CurrentRiskRatingValueID        *uint   `json:"current_risk_rating_value_id"`
	SeverityLevelTitleEn            *string `json:"severity_level_title_en"`
	SeverityLevelTitleAr            *string `json:"severity_level_title_ar"`
	CreatedAt                       string  `json:"created_at"`
	UpdatedAt                       *string `json:"updated_at"`
	CreatedUserID                   uint    `json:"created_user_id"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string  `json:"created_user_last_name_ar"`
	CreatedUserInitial              string  `json:"created_user_initial"`
	CreatedUserColor                string  `json:"created_user_color"`
	CreatedUserEmail                string  `json:"created_user_email"`
	CreatedUserProfileImagePath     string  `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn          *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string  `json:"updated_user_initial"`
	UpdatedUserColor                string  `json:"updated_user_color"`
	UpdatedUserEmail                string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string  `json:"updated_user_profile_image_path"`
	CurrentLikelihoodTitleEn        *string `json:"current_likelihood_title_en"`
	CurrentLikelihoodTitleAr        *string `json:"current_likelihood_title_ar"`
	CurrentImpactTitleEn            *string `json:"current_impact_title_en"`
	CurrentImpactTitleAr            *string `json:"current_impact_title_ar"`
	CurrentRiskRatingTitleEn        *string `json:"current_risk_rating_title_en"`
	CurrentRiskRatingTitleAr        *string `json:"current_risk_rating_title_ar"`
	CurrentRiskRatingValue          *uint   `json:"current_risk_rating_value"`
	CurrentRiskRatingColorCode      *string `json:"current_risk_rating_color_code"`
	InherentLikelihoodTitleEn       *string `json:"inherent_likelihood_title_en"`
	InherentLikelihoodTitleAr       *string `json:"inherent_likelihood_title_ar"`
	InherentImpactTitleEn           *string `json:"inherent_impact_title_en"`
	InherentImpactTitleAr           *string `json:"inherent_impact_title_ar"`
	InherentRiskRatingTitleEn       *string `json:"inherent_risk_rating_title_en"`
	InherentRiskRatingTitleAr       *string `json:"inherent_risk_rating_title_ar"`
	InherentRiskRatingValue         *uint   `json:"inherent_risk_rating_value"`
	InherentRiskRatingColorCode     *string `json:"inherent_risk_rating_color_code"`
	ResidualLikelihoodTitleEn       *string `json:"residual_likelihood_title_en"`
	ResidualLikelihoodTitleAr       *string `json:"residual_likelihood_title_ar"`
	ResidualImpactTitleEn           *string `json:"residual_impact_title_en"`
	ResidualImpactTitleAr           *string `json:"residual_impact_title_ar"`
	ResidualRiskRatingTitleEn       *string `json:"residual_risk_rating_title_en"`
	ResidualRiskRatingTitleAr       *string `json:"residual_risk_rating_title_ar"`
	ResidualRiskRatingValue         *uint   `json:"residual_risk_rating_value"`
	ResidualRiskRatingColorCode     *string `json:"residual_risk_rating_color_code"`
	RiskCode                        *string `json:"risk_code"`
	RegistrationDate                *string `json:"registration_date"`
	DivisionID                      *uint   `json:"division_id"`
	DivisionTitleEn                 *string `json:"division_title_en"`
	DivisionTitleAr                 *string `json:"division_title_ar"`
	SectionID                       *uint   `json:"section_id"`
	SectionTitleEn                  *string `json:"section_title_en"`
	SectionTitleAr                  *string `json:"section_title_ar"`
	SubSectionID                    *uint   `json:"sub_section_id"`
	SubSectionTitleEn               *string `json:"sub_section_title_en"`
	SubSectionTitleAr               *string `json:"sub_section_title_ar"`
	IsStrategic                     bool    `json:"is_strategic"`
	IsOperational                   bool    `json:"is_operational"`
	ExpectedCause                   *string `json:"expected_cause"`
	Consequence                     *string `json:"consequence"`
	IsInternal                      bool    `json:"is_internal"`
	IsExternal                      bool    `json:"is_external"`
}

type ERMDeptEvaluationListResponseDTO struct {
	Items      []*ERMDeptEvaluationRiskListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO    `json:"pagination"`
}

// ERM Department Evaluation - Process DTOs
type ProcessActivityDTO struct {
	ID                uint       `json:"id"`
	TitleEn           string     `json:"title_en"`
	TitleAr           string     `json:"title_ar"`
	DescriptionEn     *string    `json:"description_en"`
	DescriptionAr     *string    `json:"description_ar"`
	Sequence          int        `json:"sequence"`
	EstimatedDuration *string    `json:"estimated_duration"`
	ExpectedOutput    string     `json:"expected_output"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
}

type ERMDeptEvaluationProcessListDTO struct {
	ID                               uint                 `json:"id"`
	RiskUpdatePlanTimelineID         uint                 `json:"risk_update_plan_timeline_id"`
	ProcessID                        uint                 `json:"process_id"`
	ReferenceNumber                  string               `json:"reference_number"`
	TitleEn                          string               `json:"title_en"`
	TitleAr                          string               `json:"title_ar"`
	Description                      *string              `json:"description"`
	MarkedAsAuditableItem            bool                 `json:"marked_as_auditable_item"`
	ProcessGroupTitleEn              string               `json:"process_group_title_en"`
	ProcessGroupTitleAr              string               `json:"process_group_title_ar"`
	RtoFromValue                     *int                 `json:"rto_from_value"`
	RtoToValue                       *int                 `json:"rto_to_value"`
	RtoUnitId                        *uint                `json:"rto_unit_id"`
	RtoUnitTitleEn                   *string              `json:"rto_unit_title_en"`
	RtoUnitTitleAr                   *string              `json:"rto_unit_title_ar"`
	RtoCondition                     *string              `json:"rto_condition"`
	MtpdFromValue                    *int                 `json:"mtpd_from_value"`
	MtpdToValue                      *int                 `json:"mtpd_to_value"`
	MtpdUnitId                       *uint                `json:"mtpd_unit_id"`
	MtpdUnitTitleEn                  *string              `json:"mtpd_unit_title_en"`
	MtpdUnitTitleAr                  *string              `json:"mtpd_unit_title_ar"`
	MtpdCondition                    *string              `json:"mtpd_condition"`
	RpoFromValue                     *int                 `json:"rpo_from_value"`
	RpoToValue                       *int                 `json:"rpo_to_value"`
	RpoUnitId                        *uint                `json:"rpo_unit_id"`
	RpoUnitTitleEn                   *string              `json:"rpo_unit_title_en"`
	RpoUnitTitleAr                   *string              `json:"rpo_unit_title_ar"`
	RpoCondition                     *string              `json:"rpo_condition"`
	RspFromValue                     *int                 `json:"rsp_from_value"`
	RspToValue                       *int                 `json:"rsp_to_value"`
	RspUnitId                        *uint                `json:"rsp_unit_id"`
	RspUnitTitleEn                   *string              `json:"rsp_unit_title_en"`
	RspUnitTitleAr                   *string              `json:"rsp_unit_title_ar"`
	RspCondition                     *string              `json:"rsp_condition"`
	ImpactAnalysisLevelId            *uint                `json:"impact_analysis_level_id"`
	ImpactAnalysisLevelTitleEn       *string              `json:"impact_analysis_level_title_en"`
	ImpactAnalysisLevelTitleAr       *string              `json:"impact_analysis_level_title_ar"`
	ImpactAnalysisLevelColorCode     *string              `json:"impact_analysis_level_color_code"`
	ImpactAnalysisLevelType          *string              `json:"impact_analysis_level_type"`
	IsBiaPerformed                   bool                 `json:"is_bia_performed"`
	LastAssessedAt                   *time.Time           `json:"last_assessed_at"`
	LastAssessedUserID               *uint                `json:"last_assessed_user_id"`
	LastAssessedUserFirstNameEn      *string              `json:"last_assessed_user_first_name_en"`
	LastAssessedUserFirstNameAr      *string              `json:"last_assessed_user_first_name_ar"`
	LastAssessedUserLastNameEn       *string              `json:"last_assessed_user_last_name_en"`
	LastAssessedUserLastNameAr       *string              `json:"last_assessed_user_last_name_ar"`
	LastAssessedUserInitial          *string              `json:"last_assessed_user_initial"`
	LastAssessedUserColor            *string              `json:"last_assessed_user_color"`
	LastAssessedUserEmail            *string              `json:"last_assessed_user_email"`
	LastAssessedUserProfileImagePath *string              `json:"last_assessed_user_profile_image_path"`
	OwnerUnitTitleEn                 *string              `json:"owner_unit_title_en"`
	OwnerUnitTitleAr                 *string              `json:"owner_unit_title_ar"`
	OwnerID                          *uint                `json:"owner_id"`
	OwnerFirstNameEn                 *string              `json:"owner_first_name_en"`
	OwnerFirstNameAr                 *string              `json:"owner_first_name_ar"`
	OwnerLastNameEn                  *string              `json:"owner_last_name_en"`
	OwnerLastNameAr                  *string              `json:"owner_last_name_ar"`
	OwnerInitial                     string               `json:"owner_initial"`
	OwnerColor                       string               `json:"owner_color"`
	OwnerEmail                       string               `json:"owner_email"`
	OwnerProfileImagePath            string               `json:"owner_profile_image_path"`
	LifecycleStatusTitleEn           string               `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr           string               `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode         string               `json:"lifecycle_status_color_code"`
	ProcessCriticalityLevelID        uint                 `json:"process_criticality_level_id"`
	ProcessCriticalityLevelTitleEn   string               `json:"process_criticality_level_title_en"`
	ProcessCriticalityLevelTitleAr   string               `json:"process_criticality_level_title_ar"`
	CreatedUserID                    uint                 `json:"created_user_id"`
	CreatedUserFirstNameEn           string               `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr           string               `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn            string               `json:"created_user_last_name_en"`
	CreatedUserLastNameAr            string               `json:"created_user_last_name_ar"`
	CreatedUserInitial               string               `json:"created_user_initial"`
	CreatedUserColor                 string               `json:"created_user_color"`
	CreatedUserEmail                 string               `json:"created_user_email"`
	CreatedUserProfileImagePath      string               `json:"created_user_profile_image_path"`
	UpdatedUserID                    *uint                `json:"updated_user_id"`
	UpdatedUserFirstNameEn           *string              `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr           *string              `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn            *string              `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr            *string              `json:"updated_user_last_name_ar"`
	UpdatedUserInitial               string               `json:"updated_user_initial"`
	UpdatedUserColor                 string               `json:"updated_user_color"`
	UpdatedUserEmail                 string               `json:"updated_user_email"`
	UpdatedUserProfileImagePath      string               `json:"updated_user_profile_image_path"`
	CreatedAt                        time.Time            `json:"created_at"`
	UpdatedAt                        *time.Time           `json:"updated_at"`
	RiskRatingTitleEn                *string              `json:"risk_rating_title_en"`
	RiskRatingTitleAr                *string              `json:"risk_rating_title_ar"`
	RiskRatingColorCode              *string              `json:"risk_rating_color_code"`
	RiskRatingValue                  *uint                `json:"risk_rating_value"`
	RiskLevelID                      *uint                `json:"risk_level_id"`
	RiskLevelTitleEn                 *string              `json:"risk_level_title_en"`
	RiskLevelTitleAr                 *string              `json:"risk_level_title_ar"`
	RiskLevelColorCode               *string              `json:"risk_level_color_code"`
	RiskLevelType                    *string              `json:"risk_level_type"`
	Activities                       []ProcessActivityDTO `json:"activities" gorm:"-"`
}

type ERMDeptEvaluationProcessListResponseDTO struct {
	Items      []*ERMDeptEvaluationProcessListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO       `json:"pagination"`
}

type MapProcessesToUpdatePlanDTO struct {
	ProcessIDs []uint `json:"process_ids" binding:"required,min=1"`
}

// ERM Department Evaluation - Control DTOs
type ERMDeptEvaluationControlListDTO struct {
	ID                                   uint                                            `json:"id"`
	RiskUpdatePlanTimelineID             uint                                            `json:"risk_update_plan_timeline_id"`
	ControlID                            uint                                            `json:"control_id"`
	ControlTypeTitleEn                   string                                          `json:"control_type_title_en"`
	ControlTypeTitleAr                   string                                          `json:"control_type_title_ar"`
	ControlType                          string                                          `json:"control_type"`
	ControlCategoryTitleEn               string                                          `json:"control_category_title_en"`
	ControlCategoryTitleAr               string                                          `json:"control_category_title_ar"`
	OwnerUnitTitleEn                     *string                                         `json:"owner_unit_title_en"`
	OwnerUnitTitleAr                     *string                                         `json:"owner_unit_title_ar"`
	OwnerID                              *uint                                           `json:"owner_id"`
	OwnerFirstNameEn                     *string                                         `json:"owner_first_name_en"`
	OwnerFirstNameAr                     *string                                         `json:"owner_first_name_ar"`
	OwnerLastNameEn                      *string                                         `json:"owner_last_name_en"`
	OwnerLastNameAr                      *string                                         `json:"owner_last_name_ar"`
	OwnerInitial                         string                                          `json:"owner_initial"`
	OwnerColor                           string                                          `json:"owner_color"`
	OwnerEmail                           string                                          `json:"owner_email"`
	OwnerProfileImagePath                string                                          `json:"owner_profile_image_path"`
	ControlImplementationStatusTitleEn   string                                          `json:"control_implementation_status_title_en"`
	ControlImplementationStatusTitleAr   string                                          `json:"control_implementation_status_title_ar"`
	ControlImplementationStatusColorCode string                                          `json:"control_implementation_status_color_code"`
	ControlEfficiencyStatusTitleEn       string                                          `json:"control_efficiency_status_title_en"`
	ControlEfficiencyStatusTitleAr       string                                          `json:"control_efficiency_status_title_ar"`
	ControlEfficiencyStatusColorCode     string                                          `json:"control_efficiency_status_color_code"`
	ReferenceNumber                      string                                          `json:"reference_number"`
	Code                                 *string                                         `json:"code"`
	TitleEn                              string                                          `json:"title_en"`
	TitleAr                              string                                          `json:"title_ar"`
	Description                          *string                                         `json:"description"`
	Objective                            *string                                         `json:"objective"`
	ReviewFrequencyTitleEn               string                                          `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr               string                                          `json:"review_frequency_title_ar"`
	LastReviewDate                       *time.Time                                      `json:"last_review_date"`
	NextReviewDate                       *time.Time                                      `json:"next_review_date"`
	LastReviewNote                       *string                                         `json:"last_review_note"`
	TestResponsibleUserID                uint                                            `json:"test_responsible_user_id"`
	TestResponsibleUserFirstNameEn       string                                          `json:"test_responsible_user_first_name_en"`
	TestResponsibleUserFirstNameAr       string                                          `json:"test_responsible_user_first_name_ar"`
	TestResponsibleUserLastNameEn        string                                          `json:"test_responsible_user_last_name_en"`
	TestResponsibleUserLastNameAr        string                                          `json:"test_responsible_user_last_name_ar"`
	TestResponsibleUserInitial           string                                          `json:"test_responsible_user_initial"`
	TestResponsibleUserColor             string                                          `json:"test_responsible_user_color"`
	TestResponsibleUserEmail             string                                          `json:"test_responsible_user_email"`
	TestResponsibleUserProfileImagePath  string                                          `json:"test_responsible_user_profile_image_path"`
	TestFrequencyTitleEn                 *string                                         `json:"test_frequency_title_en"`
	TestFrequencyTitleAr                 *string                                         `json:"test_frequency_title_ar"`
	LastTestDate                         *time.Time                                      `json:"last_test_date"`
	NextTestDate                         *time.Time                                      `json:"next_test_date"`
	ImplementedAt                        *time.Time                                      `json:"implemented_at"`
	IsApplicable                         bool                                            `json:"is_applicable"`
	Justification                        *string                                         `json:"justification"`
	ImplementationGuideline              *string                                         `json:"implementation_guideline"`
	ApplicabilityChangeDate              *string                                         `json:"applicability_change_date"`
	LifecycleStatusTitleEn               string                                          `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr               string                                          `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode             string                                          `json:"lifecycle_status_color_code"`
	ControlStatusID                      *uint                                           `json:"control_status_id"`
	ControlStatusTitleEn                 *string                                         `json:"control_status_title_en"`
	ControlStatusTitleAr                 *string                                         `json:"control_status_title_ar"`
	ControlStatusColorCode               *string                                         `json:"control_status_color_code"`
	ControlStatusType                    *string                                         `json:"control_status_type"`
	CreatedUserID                        uint                                            `json:"created_user_id"`
	CreatedUserFirstNameEn               string                                          `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr               string                                          `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                string                                          `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                string                                          `json:"created_user_last_name_ar"`
	CreatedUserInitial                   string                                          `json:"created_user_initial"`
	CreatedUserColor                     string                                          `json:"created_user_color"`
	CreatedUserEmail                     string                                          `json:"created_user_email"`
	CreatedUserProfileImagePath          string                                          `json:"created_user_profile_image_path"`
	UpdatedUserID                        *uint                                           `json:"updated_user_id"`
	UpdatedUserFirstNameEn               *string                                         `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr               *string                                         `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                *string                                         `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                *string                                         `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                   string                                          `json:"updated_user_initial"`
	UpdatedUserColor                     string                                          `json:"updated_user_color"`
	UpdatedUserEmail                     string                                          `json:"updated_user_email"`
	UpdatedUserProfileImagePath          string                                          `json:"updated_user_profile_image_path"`
	CreatedAt                            time.Time                                       `json:"created_at"`
	UpdatedAt                            *time.Time                                      `json:"updated_at"`
	RiskRatingTitleEn                    *string                                         `json:"risk_rating_title_en"`
	RiskRatingTitleAr                    *string                                         `json:"risk_rating_title_ar"`
	RiskRatingColorCode                  *string                                         `json:"risk_rating_color_code"`
	RiskRatingValue                      *uint                                           `json:"risk_rating_value"`
	ControlAutomationID                  *uint                                           `json:"control_automation_id"`
	ControlAutomationTitleEn             *string                                         `json:"control_automation_title_en"`
	ControlAutomationTitleAr             *string                                         `json:"control_automation_title_ar"`
	ControlAutomationType                *string                                         `json:"control_automation_type"`
	ControlAutomationScore               *float64                                        `json:"control_automation_score"`
	ControlAutomationStatus              *string                                         `json:"control_automation_status"`
	ControlConsistencyID                 *uint                                           `json:"control_consistency_id"`
	ControlConsistencyTitleEn            *string                                         `json:"control_consistency_title_en"`
	ControlConsistencyTitleAr            *string                                         `json:"control_consistency_title_ar"`
	ControlConsistencyType               *string                                         `json:"control_consistency_type"`
	ControlConsistencyScore              *float64                                        `json:"control_consistency_score"`
	ControlConsistencyStatus             *string                                         `json:"control_consistency_status"`
	ControlConsistencyDescriptionEn      *string                                         `json:"control_consistency_description_en"`
	ControlConsistencyDescriptionAr      *string                                         `json:"control_consistency_description_ar"`
	ControlEffectivenessScore            *float64                                        `json:"control_effectiveness_score"`
	ControlEffectivenessRatingID         *uint                                           `json:"control_effectiveness_rating_id"`
	ControlEffectivenessRatingTitleEn    *string                                         `json:"control_effectiveness_rating_title_en"`
	ControlEffectivenessRatingTitleAr    *string                                         `json:"control_effectiveness_rating_title_ar"`
	ControlEffectivenessRatingMinScore   *float64                                        `json:"control_effectiveness_rating_min_score"`
	ControlEffectivenessRatingMaxScore   *float64                                        `json:"control_effectiveness_rating_max_score"`
	ControlEffectivenessRatingType       *string                                         `json:"control_effectiveness_rating_type"`
	ControlEffectivenessRatingColorCode  *string                                         `json:"control_effectiveness_rating_color_code"`
	RiskReduction                        *float64                                        `json:"risk_reduction"`
	ControlTestID                        *uint                                           `json:"control_test_id"`
	ControlTestTitle                     *string                                         `json:"control_test_title"`
	ControlTestDescription               *string                                         `json:"control_test_description"`
	ControlTestComments                  *string                                         `json:"control_test_comments"`
	ControlTestScore                     *float64                                        `json:"control_test_score"`
	ControlTestEffectivenessScore        *float64                                        `json:"control_test_effectiveness_score"`
	ControlTestRiskReduction             *float64                                        `json:"control_test_risk_reduction"`
	ControlTestTestedAt                  *string                                         `json:"control_test_tested_at"`
	ControlTestTestedByID                *uint                                           `json:"control_test_tested_by_id"`
	ControlTestTestedByFirstNameEn       *string                                         `json:"control_test_tested_by_first_name_en"`
	ControlTestTestedByFirstNameAr       *string                                         `json:"control_test_tested_by_first_name_ar"`
	ControlTestTestedByLastNameEn        *string                                         `json:"control_test_tested_by_last_name_en"`
	ControlTestTestedByLastNameAr        *string                                         `json:"control_test_tested_by_last_name_ar"`
	ControlTestTestedByInitial           *string                                         `json:"control_test_tested_by_initial"`
	ControlTestTestedByColor             *string                                         `json:"control_test_tested_by_color"`
	ControlTestTestedByEmail             *string                                         `json:"control_test_tested_by_email"`
	ControlTestTestedByProfileImagePath  *string                                         `json:"control_test_tested_by_profile_image_path"`
	ControlTestDocuments                 []organizationDTOs.ControlTestDocumentDTO       `json:"control_test_documents" gorm:"-"`
	AssociatedDocuments                  []organizationDTOs.ControlAssociatedDocumentDTO `json:"associated_documents" gorm:"-"`
}

type ERMDeptEvaluationControlListResponseDTO struct {
	Items      []*ERMDeptEvaluationControlListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO       `json:"pagination"`
}

type MapControlsToUpdatePlanDTO struct {
	ControlIDs []uint `json:"control_ids" binding:"required,min=1"`
}

// ERM Department Evaluation DTOs
type UpdateBusinessContextStatusFormDTO struct {
	IsBusinessContextUpdated *bool `json:"is_business_context_updated" binding:"required"`
}

type UpdateBusinessContextFormDTO struct {
	InternalAuditNote          *string `json:"internal_audit_note"`
	ExternalContext            *string `json:"external_context"`
	AnnualObjectiveDescription *string `json:"annual_objective_description"`
}

type BusinessContextDetailDTO struct {
	ID                         uint                             `json:"id"`
	RiskUpdatePlanTimeline     *RiskUpdatePlanTimelineDetailDTO `json:"risk_update_plan_timeline"`
	InternalAuditNote          *string                          `json:"internal_audit_note"`
	ExternalContext            *string                          `json:"external_context"`
	AnnualObjectiveDescription *string                          `json:"annual_objective_description"`
	CreatedAt                  time.Time                        `json:"created_at"`
	CreatedByUser              *organizationDTOs.UserMinimalDTO `json:"created_by"`
}

// Timeline Not Applicable DTO
type MarkTimelineNotApplicableDTO struct {
	Justification string `json:"justification" binding:"required"`
}

// LeaderShip Assessment DTO
type LeaderShipAssessmentUserListDTO struct {
	ID                       uint    `json:"id"`
	FirstNameEn              string  `json:"first_name_en"`
	LastNameEn               string  `json:"last_name_en"`
	FirstNameAr              string  `json:"first_name_ar"`
	LastNameAr               string  `json:"last_name_ar"`
	OrganizationUnit1TitleEn string  `json:"organization_unit_1_title_en" gorm:"column:organization_unit_1_title_en"`
	OrganizationUnit2TitleEn string  `json:"organization_unit_2_title_en,omitempty" gorm:"column:organization_unit_2_title_en"`
	OrganizationUnit3TitleEn string  `json:"organization_unit_3_title_en,omitempty" gorm:"column:organization_unit_3_title_en"`
	OrganizationUnit4TitleEn string  `json:"organization_unit_4_title_en,omitempty" gorm:"column:organization_unit_4_title_en"`
	OrganizationUnit5TitleEn string  `json:"organization_unit_5_title_en,omitempty" gorm:"column:organization_unit_5_title_en"`
	OrganizationUnit1TitleAr string  `json:"organization_unit_1_title_ar" gorm:"column:organization_unit_1_title_ar"`
	OrganizationUnit2TitleAr string  `json:"organization_unit_2_title_ar,omitempty" gorm:"column:organization_unit_2_title_ar"`
	OrganizationUnit3TitleAr string  `json:"organization_unit_3_title_ar,omitempty" gorm:"column:organization_unit_3_title_ar"`
	OrganizationUnit4TitleAr string  `json:"organization_unit_4_title_ar,omitempty" gorm:"column:organization_unit_4_title_ar"`
	OrganizationUnit5TitleAr string  `json:"organization_unit_5_title_ar,omitempty" gorm:"column:organization_unit_5_title_ar"`
	DesignationTitleEn       string  `json:"designation_title_en"`
	DesignationTitleAr       string  `json:"designation_title_ar"`
	RolesEn                  string  `json:"roles_en"`
	RolesAr                  string  `json:"roles_ar"`
	LicenseTypeID            *uint   `json:"license_type_id"`
	LicenseTypeTitleEn       *string `json:"license_type_title_en"`
	LicenseTypeTitleAr       *string `json:"license_type_title_ar"`
	LicenseTypeType          *string `json:"license_type_type"`
	LicenseTypeColorCode     *string `json:"license_type_color_code"`
	Email                    string  `json:"email"`
	ProfileImagePath         *string `json:"profile_image_path"`
	Initial                  string  `json:"initial"`
	Color                    string  `json:"color"`
	MobileNumber             string  `json:"mobile_number"`
	UserStatusTitleEn        string  `json:"user_status_title_en"`
	UserStatusTitleAr        string  `json:"user_status_title_ar"`
	UserStatusColorCode      string  `json:"user_status_color_code"`
	UserStatusIsActive       bool    `json:"user_status_is_active"`
	ReportingUsersCount      int     `json:"reporting_users_count"`
}

type LeaderShipAssessmentUserListResponseDTO struct {
	Items      []*LeaderShipAssessmentUserListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO       `json:"pagination"`
}

// Risk Team Assessment DTO
type RiskTeamAssessmentUserListDTO struct {
	ID                       uint                             `json:"id"`
	RiskUpdatePlanTimelineID uint                             `json:"risk_update_plan_timeline_id"`
	RiskChampionID           uint                             `json:"risk_champion_id"`
	RiskChampion             *organizationDTOs.UserDetailsDTO `json:"risk_champion"`
	CreatedAt                time.Time                        `json:"created_at"`
	UpdatedAt                *time.Time                       `json:"updated_at"`
}

// Risk Identification

type RiskIdentificationListDTO struct {
	ID                              uint    `json:"id"`
	ReferenceNumber                 string  `json:"reference_number"`
	RiskCategoryTitleEn             *string `json:"risk_category_title_en"`
	RiskCategoryTitleAr             *string `json:"risk_category_title_ar"`
	RiskSourceTitleEn               *string `json:"risk_source_title_en"`
	RiskSourceTitleAr               *string `json:"risk_source_title_ar"`
	ResponsibleUnitID               *uint   `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn          *string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          *string `json:"responsible_unit_title_ar"`
	ResponsibleUserID               *uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      *string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      *string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       *string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       *string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string `json:"responsible_user_initial"`
	ResponsibleUserColor            *string `json:"responsible_user_color"`
	ResponsibleUserEmail            *string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath *string `json:"responsible_user_profile_image_path"`
	RiskStatusTitleEn               string  `json:"risk_status_title_en"`
	RiskStatusTitleAr               string  `json:"risk_status_title_ar"`
	RiskResponseStrategyTitleEn     string  `json:"risk_response_strategy_title_en"`
	RiskResponseStrategyTitleAr     string  `json:"risk_response_strategy_title_ar"`
	RiskStatusType                  string  `json:"risk_status_type"`
	RiskStatusColorCode             string  `json:"risk_status_color_code"`
	Title                           string  `json:"title"`
	Description                     *string `json:"description"`
	ReportedAt                      *string `json:"reported_at"`
	ReviewFrequencyTitleEn          *string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr          *string `json:"review_frequency_title_ar"`
	NextReviewDate                  *string `json:"next_review_date"`
	CurrentRiskRatingValueID        *uint   `json:"current_risk_rating_value_id"`
	Comments                        *string `json:"comments"`
	SeverityLevelTitleEn            *string `json:"severity_level_title_en"`
	SeverityLevelTitleAr            *string `json:"severity_level_title_ar"`
	CreatedAt                       string  `json:"created_at"`
	UpdatedAt                       *string `json:"updated_at"`
	CreatedUserID                   uint    `json:"created_user_id"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string  `json:"created_user_last_name_ar"`
	CreatedUserInitial              string  `json:"created_user_initial"`
	CreatedUserColor                string  `json:"created_user_color"`
	CreatedUserEmail                string  `json:"created_user_email"`
	CreatedUserProfileImagePath     string  `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn          *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string  `json:"updated_user_initial"`
	UpdatedUserColor                string  `json:"updated_user_color"`
	UpdatedUserEmail                string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string  `json:"updated_user_profile_image_path"`
	CurrentLikelihoodTitleEn        *string `json:"current_likelihood_title_en"`
	CurrentLikelihoodTitleAr        *string `json:"current_likelihood_title_ar"`
	CurrentImpactTitleEn            *string `json:"current_impact_title_en"`
	CurrentImpactTitleAr            *string `json:"current_impact_title_ar"`
	CurrentRiskRatingTitleEn        *string `json:"current_risk_rating_title_en"`
	CurrentRiskRatingTitleAr        *string `json:"current_risk_rating_title_ar"`
	CurrentRiskRatingValue          *uint   `json:"current_risk_rating_value"`
	CurrentRiskRatingColorCode      *string `json:"current_risk_rating_color_code"`
	InherentLikelihoodTitleEn       *string `json:"inherent_likelihood_title_en"`
	InherentLikelihoodTitleAr       *string `json:"inherent_likelihood_title_ar"`
	InherentImpactTitleEn           *string `json:"inherent_impact_title_en"`
	InherentImpactTitleAr           *string `json:"inherent_impact_title_ar"`
	InherentRiskRatingTitleEn       *string `json:"inherent_risk_rating_title_en"`
	InherentRiskRatingTitleAr       *string `json:"inherent_risk_rating_title_ar"`
	InherentRiskRatingValue         *uint   `json:"inherent_risk_rating_value"`
	InherentRiskRatingColorCode     *string `json:"inherent_risk_rating_color_code"`
	ResidualLikelihoodTitleEn       *string `json:"residual_likelihood_title_en"`
	ResidualLikelihoodTitleAr       *string `json:"residual_likelihood_title_ar"`
	ResidualImpactTitleEn           *string `json:"residual_impact_title_en"`
	ResidualImpactTitleAr           *string `json:"residual_impact_title_ar"`
	ResidualRiskRatingTitleEn       *string `json:"residual_risk_rating_title_en"`
	ResidualRiskRatingTitleAr       *string `json:"residual_risk_rating_title_ar"`
	ResidualRiskRatingValue         *uint   `json:"residual_risk_rating_value"`
	ResidualRiskRatingColorCode     *string `json:"residual_risk_rating_color_code"`
	RiskCode                        *string `json:"risk_code"`
	RegistrationDate                *string `json:"registration_date"`
	DivisionID                      *uint   `json:"division_id"`
	DivisionTitleEn                 *string `json:"division_title_en"`
	DivisionTitleAr                 *string `json:"division_title_ar"`
	SectionID                       *uint   `json:"section_id"`
	SectionTitleEn                  *string `json:"section_title_en"`
	SectionTitleAr                  *string `json:"section_title_ar"`
	SubSectionID                    *uint   `json:"sub_section_id"`
	SubSectionTitleEn               *string `json:"sub_section_title_en"`
	SubSectionTitleAr               *string `json:"sub_section_title_ar"`
	IsStrategic                     bool    `json:"is_strategic"`
	IsOperational                   bool    `json:"is_operational"`
	ExpectedCause                   *string `json:"expected_cause"`
	Consequence                     *string `json:"consequence"`
	IsInternal                      bool    `json:"is_internal"`
	IsExternal                      bool    `json:"is_external"`
}

type RiskIdentificationListResponseDTO struct {
	Items      []*RiskIdentificationListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
type UpdatePlanRiskRegisterFormDTO struct {
	RiskCategoryID                 *uint   `json:"risk_category_id"`
	RiskSourceID                   *uint   `json:"risk_source_id"`
	ResponsibleUserID              *uint   `json:"responsible_user_id"`
	ResponsibleUnitID              *uint   `json:"responsible_unit_id" binding:"required"`
	Title                          string  `json:"title" binding:"required,max=255"`
	Description                    *string `json:"description"`
	ReportedAt                     *string `json:"reported_at"`
	ReviewFrequencyID              *uint   `json:"review_frequency_id"`
	NextReviewDate                 *string `json:"next_review_date" binding:"omitempty,datetime=2006-01-02"`
	LastReviewNote                 *string `json:"last_review_note" binding:"omitempty"`
	InherentRiskMatrixLikelihoodID *uint   `json:"inherent_risk_matrix_likelihood_id"`
	InherentRiskMatrixImpactID     *uint   `json:"inherent_risk_matrix_impact_id"`
	RiskResponseStrategyID         *uint   `json:"risk_response_strategy_id"`
	ResidualRiskMatrixLikelihoodID *uint   `json:"residual_risk_matrix_likelihood_id"`
	ResidualRiskMatrixImpactID     *uint   `json:"residual_risk_matrix_impact_id"`
	RiskCode                       *string `json:"risk_code"`
	RegistrationDate               string  `json:"registration_date" binding:"required,datetime=2006-01-02"`
	DivisionID                     uint    `json:"division_id" binding:"required"`
	SectionID                      *uint   `json:"section_id"`
	SubSectionID                   *uint   `json:"sub_section_id"`
	IsStrategic                    bool    `json:"is_strategic"`
	IsOperational                  bool    `json:"is_operational"`
	ExpectedCause                  *string `json:"expected_cause"`
	Consequence                    *string `json:"consequence"`
	IsInternal                     bool    `json:"is_internal"`
	IsExternal                     bool    `json:"is_external"`
}

// CONTROL REGISTER
type ControlRegisterUpdatedDTO struct {
	NotImplemented       int64 `json:"not_implemented"`
	PartiallyImplemented int64 `json:"partially_implemented"`
	Implemented          int64 `json:"implemented"`
	NotApplicable        int64 `json:"not_applicable"`
}

type UpdatePlanControlTestFormDTO struct {
	ControlID                     uint                          `json:"control_id" binding:"required"`
	ControlImplementationStatusID uint                          `json:"control_implementation_status_id" binding:"required"`
	ControlEfficiencyStatusID     *uint                         `json:"control_efficiency_status_id,omitempty" binding:"omitempty"`
	Title                         string                        `json:"title" binding:"required"`
	Description                   *string                       `json:"description"  binding:"omitempty"`
	Score                         *uint                         `json:"score" binding:"omitempty"`
	ControlAutomationID           uint                          `json:"control_automation_id" binding:"required"`
	ControlConsistencyID          uint                          `json:"control_consistency_id" binding:"required"`
	Comments                      *string                       `json:"comments" binding:"omitempty"`
	TestedAt                      *string                       `json:"tested_at" binding:"required"` // assuming this will be a date in string format
	TestedBy                      uint                          `json:"tested_by" binding:"required"`
	Documents                     []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

// Risk Analysis

type RiskAnalysisListDTO struct {
	ID                          uint                            `json:"id"`
	ReferenceNumber             string                          `json:"reference_number"`
	Title                       string                          `json:"title"`
	Description                 *string                         `json:"description"`
	RiskEvaluationID            *uint                           `json:"-"`
	InherentLikelihoodID        *uint                           `json:"inherent_likelihood_id"`
	InherentImpactID            *uint                           `json:"inherent_impact_id"`
	InherentLikelihoodTitleEn   *string                         `json:"inherent_likelihood_title_en"`
	InherentLikelihoodTitleAr   *string                         `json:"inherent_likelihood_title_ar"`
	InherentImpactTitleEn       *string                         `json:"inherent_impact_title_en"`
	InherentImpactTitleAr       *string                         `json:"inherent_impact_title_ar"`
	InherentRiskRatingTitleEn   *string                         `json:"inherent_risk_rating_title_en"`
	InherentRiskRatingTitleAr   *string                         `json:"inherent_risk_rating_title_ar"`
	InherentRiskRatingValue     *uint                           `json:"inherent_risk_rating_value"`
	InherentRiskRatingColorCode *string                         `json:"inherent_risk_rating_color_code"`
	ImpactCategories            []RiskAnalysisImpactCategoryDTO `json:"impact_categories" gorm:"-"`
}

type RiskAnalysisImpactCategoryDTO struct {
	ID         uint                       `json:"id"`
	TitleEn    string                     `json:"title_en"`
	TitleAr    string                     `json:"title_ar"`
	Guidelines []RiskAnalysisGuidelineDTO `json:"guidelines" gorm:"-"`
}

type RiskAnalysisGuidelineDTO struct {
	ID        uint   `json:"id"`
	Guideline string `json:"guideline"`
}

type RiskAnalysisListResponseDTO struct {
	Items      []*RiskAnalysisListDTO       `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type RiskAnalysisFormDTO struct {
	Analysis []*UpdatePlanInherentEvaluationItemDTO `json:"analysis" binding:"required"`
}

type UpdatePlanInherentEvaluationItemDTO struct {
	RiskRegisterID         uint                                   `json:"risk_register_id" binding:"required"`
	RiskMatrixLikelihoodID uint                                   `json:"risk_matrix_likelihood_id" binding:"required"`
	RiskMatrixImpactID     uint                                   `json:"risk_matrix_impact_id" binding:"required"`
	RiskImpactAnalyses     []*UpdatePlanRiskImpactAnalysisFormDTO `json:"risk_impact_analyses"`
}

type UpdatePlanRiskImpactAnalysisFormDTO struct {
	RiskImpactCategoryID uint `json:"risk_impact_category_id"`
	RiskMatrixImpactID   uint `json:"risk_matrix_impact_id"`
}

// Risk Control Associations

type RiskControlAssociationFormDTO struct {
	Residual []*UpdatePlanResidualControlItemDTO `json:"residual" binding:"required"`
}
type UpdatePlanResidualControlItemDTO struct {
	RiskRegisterID uint `json:"risk_register_id" binding:"required"`
	ControlID      uint `json:"control_id" binding:"required"`
}

type RiskControlAssociationListDTO struct {
	ID                              uint                                  `json:"id"`
	ReferenceNumber                 string                                `json:"reference_number"`
	RiskCategoryTitleEn             *string                               `json:"risk_category_title_en"`
	RiskCategoryTitleAr             *string                               `json:"risk_category_title_ar"`
	RiskSourceTitleEn               *string                               `json:"risk_source_title_en"`
	RiskSourceTitleAr               *string                               `json:"risk_source_title_ar"`
	ResponsibleUnitID               *uint                                 `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn          *string                               `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          *string                               `json:"responsible_unit_title_ar"`
	ResponsibleUserID               *uint                                 `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      *string                               `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      *string                               `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       *string                               `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       *string                               `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string                               `json:"responsible_user_initial"`
	ResponsibleUserColor            *string                               `json:"responsible_user_color"`
	ResponsibleUserEmail            *string                               `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath *string                               `json:"responsible_user_profile_image_path"`
	RiskStatusTitleEn               string                                `json:"risk_status_title_en"`
	RiskStatusTitleAr               string                                `json:"risk_status_title_ar"`
	RiskResponseStrategyID          *uint                                 `json:"risk_response_strategy_id"`
	RiskResponseStrategyTitleEn     string                                `json:"risk_response_strategy_title_en"`
	RiskResponseStrategyTitleAr     string                                `json:"risk_response_strategy_title_ar"`
	RiskStatusType                  string                                `json:"risk_status_type"`
	RiskStatusColorCode             string                                `json:"risk_status_color_code"`
	Title                           string                                `json:"title"`
	Description                     *string                               `json:"description"`
	ReportedAt                      *string                               `json:"reported_at"`
	ReviewFrequencyTitleEn          *string                               `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr          *string                               `json:"review_frequency_title_ar"`
	NextReviewDate                  *string                               `json:"next_review_date"`
	CurrentRiskRatingValueID        *uint                                 `json:"current_risk_rating_value_id"`
	Comments                        *string                               `json:"comments"`
	SeverityLevelTitleEn            *string                               `json:"severity_level_title_en"`
	SeverityLevelTitleAr            *string                               `json:"severity_level_title_ar"`
	CreatedAt                       string                                `json:"created_at"`
	UpdatedAt                       *string                               `json:"updated_at"`
	CreatedUserID                   uint                                  `json:"created_user_id"`
	CreatedUserFirstNameEn          string                                `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string                                `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string                                `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string                                `json:"created_user_last_name_ar"`
	CreatedUserInitial              string                                `json:"created_user_initial"`
	CreatedUserColor                string                                `json:"created_user_color"`
	CreatedUserEmail                string                                `json:"created_user_email"`
	CreatedUserProfileImagePath     string                                `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint                                 `json:"updated_user_id"`
	UpdatedUserFirstNameEn          *string                               `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string                               `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string                               `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string                               `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string                                `json:"updated_user_initial"`
	UpdatedUserColor                string                                `json:"updated_user_color"`
	UpdatedUserEmail                string                                `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string                                `json:"updated_user_profile_image_path"`
	CurrentLikelihoodTitleEn        *string                               `json:"current_likelihood_title_en"`
	CurrentLikelihoodTitleAr        *string                               `json:"current_likelihood_title_ar"`
	CurrentImpactTitleEn            *string                               `json:"current_impact_title_en"`
	CurrentImpactTitleAr            *string                               `json:"current_impact_title_ar"`
	CurrentRiskRatingTitleEn        *string                               `json:"current_risk_rating_title_en"`
	CurrentRiskRatingTitleAr        *string                               `json:"current_risk_rating_title_ar"`
	CurrentRiskRatingValue          *uint                                 `json:"current_risk_rating_value"`
	CurrentRiskRatingColorCode      *string                               `json:"current_risk_rating_color_code"`
	InherentLikelihoodTitleEn       *string                               `json:"inherent_likelihood_title_en"`
	InherentLikelihoodTitleAr       *string                               `json:"inherent_likelihood_title_ar"`
	InherentImpactTitleEn           *string                               `json:"inherent_impact_title_en"`
	InherentImpactTitleAr           *string                               `json:"inherent_impact_title_ar"`
	InherentRiskRatingTitleEn       *string                               `json:"inherent_risk_rating_title_en"`
	InherentRiskRatingTitleAr       *string                               `json:"inherent_risk_rating_title_ar"`
	InherentRiskRatingValue         *uint                                 `json:"inherent_risk_rating_value"`
	InherentRiskRatingColorCode     *string                               `json:"inherent_risk_rating_color_code"`
	ResidualLikelihoodTitleEn       *string                               `json:"residual_likelihood_title_en"`
	ResidualLikelihoodTitleAr       *string                               `json:"residual_likelihood_title_ar"`
	ResidualImpactTitleEn           *string                               `json:"residual_impact_title_en"`
	ResidualImpactTitleAr           *string                               `json:"residual_impact_title_ar"`
	ResidualRiskRatingTitleEn       *string                               `json:"residual_risk_rating_title_en"`
	ResidualRiskRatingTitleAr       *string                               `json:"residual_risk_rating_title_ar"`
	ResidualRiskRatingValue         *uint                                 `json:"residual_risk_rating_value"`
	ResidualRiskRatingColorCode     *string                               `json:"residual_risk_rating_color_code"`
	RiskCode                        *string                               `json:"risk_code"`
	RegistrationDate                *string                               `json:"registration_date"`
	DivisionID                      *uint                                 `json:"division_id"`
	DivisionTitleEn                 *string                               `json:"division_title_en"`
	DivisionTitleAr                 *string                               `json:"division_title_ar"`
	SectionID                       *uint                                 `json:"section_id"`
	SectionTitleEn                  *string                               `json:"section_title_en"`
	SectionTitleAr                  *string                               `json:"section_title_ar"`
	SubSectionID                    *uint                                 `json:"sub_section_id"`
	SubSectionTitleEn               *string                               `json:"sub_section_title_en"`
	SubSectionTitleAr               *string                               `json:"sub_section_title_ar"`
	IsStrategic                     bool                                  `json:"is_strategic"`
	IsOperational                   bool                                  `json:"is_operational"`
	ExpectedCause                   *string                               `json:"expected_cause"`
	Consequence                     *string                               `json:"consequence"`
	IsInternal                      bool                                  `json:"is_internal"`
	IsExternal                      bool                                  `json:"is_external"`
	ControlAnalysis                 *RiskControlAssociationAnalysisDTO    `json:"control_analysis" gorm:"-"`
	ActionPlans                     []RiskControlAssociationActionPlanDTO `json:"action_plans" gorm:"-"`
}

type RiskControlAssociationAnalysisDTO struct {
	ID                                 uint     `json:"id"`
	Type                               string   `json:"type"`
	RiskControlID                      uint     `json:"risk_control_id"`
	ControlID                          uint     `json:"control_id"`
	ControlTitleEn                     string   `json:"control_title_en"`
	ControlTitleAr                     string   `json:"control_title_ar"`
	ControlEfficiencyStatusID          *uint    `json:"control_efficiency_status_id"`
	ControlEfficiencyStatusTitleEn     *string  `json:"control_efficiency_status_title_en"`
	ControlEfficiencyStatusTitleAr     *string  `json:"control_efficiency_status_title_ar"`
	ControlImplementationStatusID      uint     `json:"control_implementation_status_id"`
	ControlImplementationStatusTitleEn string   `json:"control_implementation_status_title_en"`
	ControlImplementationStatusTitleAr string   `json:"control_implementation_status_title_ar"`
	Score                              *uint    `json:"score"`
	ControlEffectivenessScore          *float64 `json:"control_effectiveness_score"`
	ControlEffectivenessRatingID       *uint    `json:"control_effectiveness_rating_id"`
	ControlEffectivenessRatingTitleEn  *string  `json:"control_effectiveness_rating_title_en"`
	ControlEffectivenessRatingTitleAr  *string  `json:"control_effectiveness_rating_title_ar"`
	RiskReduction                      *float64 `json:"risk_reduction"`
}

type RiskControlAssociationActionPlanDTO struct {
	ID                              uint    `json:"id"`
	RiskRegisterID                  uint    `json:"risk_register_id"`
	Title                           string  `json:"title"`
	Description                     *string `json:"description"`
	SeverityLevelID                 *uint   `json:"severity_level_id"`
	SeverityLevelTitleEn            *string `json:"severity_level_title_en"`
	SeverityLevelTitleAr            *string `json:"severity_level_title_ar"`
	ResponsibleUnitID               *uint   `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn          *string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          *string `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	Comments                        *string `json:"comments"`
	StartDate                       *string `json:"start_date"`
	TargetDate                      *string `json:"target_date"`
	Cost                            *int    `json:"cost"`
	ActionPlanStatusID              uint    `json:"action_plan_status_id"`
	ActionPlanStatusTitleEn         string  `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr         string  `json:"action_plan_status_title_ar"`
	ActionPlanStatusType            string  `json:"action_plan_status_type"`
	ActionPlanStatusColorCode       string  `json:"action_plan_status_color_code"`
	ReviewUserID                    *uint   `json:"review_user_id"`
	ReviewUserFirstNameEn           *string `json:"review_user_first_name_en"`
	ReviewUserFirstNameAr           *string `json:"review_user_first_name_ar"`
	ReviewUserLastNameEn            *string `json:"review_user_last_name_en"`
	ReviewUserLastNameAr            *string `json:"review_user_last_name_ar"`
	ReviewUserInitial               *string `json:"review_user_initial"`
	ReviewUserColor                 *string `json:"review_user_color"`
	ReviewUserEmail                 *string `json:"review_user_email"`
	ReviewUserProfileImagePath      *string `json:"review_user_profile_image_path"`
	Percentage                      *uint   `json:"percentage"`
	Notes                           *string `json:"notes"`
	DelegateeUserID                 *uint   `json:"delegatee_user_id"`
	DelegateeUserFirstNameEn        *string `json:"delegatee_user_first_name_en"`
	DelegateeUserFirstNameAr        *string `json:"delegatee_user_first_name_ar"`
	DelegateeUserLastNameEn         *string `json:"delegatee_user_last_name_en"`
	DelegateeUserLastNameAr         *string `json:"delegatee_user_last_name_ar"`
	DelegateeUserInitial            *string `json:"delegatee_user_initial"`
	DelegateeUserColor              *string `json:"delegatee_user_color"`
	DelegateeUserEmail              *string `json:"delegatee_user_email"`
	DelegateeUserProfileImagePath   *string `json:"delegatee_user_profile_image_path"`
	CreatedAt                       string  `json:"created_at"`
	UpdatedAt                       *string `json:"updated_at"`
	CreatedUserID                   uint    `json:"created_user_id"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string  `json:"created_user_last_name_ar"`
	CreatedUserInitial              string  `json:"created_user_initial"`
	CreatedUserColor                string  `json:"created_user_color"`
	CreatedUserEmail                string  `json:"created_user_email"`
	CreatedUserProfileImagePath     string  `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn          *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              *string `json:"updated_user_initial"`
	UpdatedUserColor                *string `json:"updated_user_color"`
	UpdatedUserEmail                *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string  `json:"updated_user_profile_image_path"`
}

// Timeline Detail Related DTOs
type RiskUpdatePlanErmDeptDocumentDTO struct {
	ID                       uint       `json:"id"`
	RiskUpdatePlanTimelineID uint       `json:"risk_update_plan_timeline_id"`
	URL                      *string    `json:"url"`
	Name                     *string    `json:"name"`
	Size                     *int       `json:"size"`
	Token                    *string    `json:"token"`
	OriginalName             *string    `json:"original_name"`
	Type                     *string    `json:"type"`
	CreatedBYID              *uint      `json:"created_by_id"`
	UpdatedBYID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
}

type RiskUpdatePlanIntroductoryMeetingDocumentDTO struct {
	ID                       uint       `json:"id"`
	RiskUpdatePlanTimelineID uint       `json:"risk_update_plan_timeline_id"`
	URL                      *string    `json:"url"`
	Name                     *string    `json:"name"`
	Size                     *int       `json:"size"`
	Token                    *string    `json:"token"`
	OriginalName             *string    `json:"original_name"`
	Type                     *string    `json:"type"`
	CreatedBYID              *uint      `json:"created_by_id"`
	UpdatedBYID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
}

type RiskUpdatePlanBusinessContextDTO struct {
	ID                       uint       `json:"id"`
	RiskUpdatePlanTimelineID uint       `json:"risk_update_plan_timeline_id"`
	Title                    string     `json:"title"`
	Description              *string    `json:"description"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
}

type RiskUpdatePlanRiskDTO struct {
	ID                       uint                   `json:"id"`
	RiskUpdatePlanTimelineID uint                   `json:"risk_update_plan_timeline_id"`
	RiskRegisterID           uint                   `json:"risk_register_id"`
	RiskRegister             *RiskRegisterDetailDTO `json:"risk_register"`
	CreatedBYID              *uint                  `json:"created_by_id"`
	CreatedAt                time.Time              `json:"created_at"`
}

type RiskUpdatePlanControlDTO struct {
	ID                       uint                               `json:"id"`
	RiskUpdatePlanTimelineID uint                               `json:"risk_update_plan_timeline_id"`
	ControlID                uint                               `json:"control_id"`
	Control                  *organizationDTOs.ControlDetailDTO `json:"control"`
	CreatedBYID              *uint                              `json:"created_by_id"`
	CreatedAt                time.Time                          `json:"created_at"`
}

type RiskUpdatePlanProcessDTO struct {
	ID                       uint                               `json:"id"`
	RiskUpdatePlanTimelineID uint                               `json:"risk_update_plan_timeline_id"`
	ProcessID                uint                               `json:"process_id"`
	Process                  *organizationDTOs.ProcessDetailDTO `json:"process"`
	CreatedBYID              *uint                              `json:"created_by_id"`
	CreatedAt                time.Time                          `json:"created_at"`
}

type RiskUpdatePlanLeadershipAssessmentDTO struct {
	ID                       uint                             `json:"id"`
	RiskUpdatePlanTimelineID uint                             `json:"risk_update_plan_timeline_id"`
	UserID                   uint                             `json:"user_id"`
	User                     *organizationDTOs.UserMinimalDTO `json:"user"`
	CreatedBYID              *uint                            `json:"created_by_id"`
	UpdatedBYID              *uint                            `json:"updated_by_id"`
	CreatedAt                time.Time                        `json:"created_at"`
	UpdatedAt                *time.Time                       `json:"updated_at"`
}

type RiskUpdatePlanLeadershipDocumentDTO struct {
	ID                       uint       `json:"id"`
	RiskUpdatePlanTimelineID uint       `json:"risk_update_plan_timeline_id"`
	URL                      *string    `json:"url"`
	Name                     *string    `json:"name"`
	Size                     *int       `json:"size"`
	Token                    *string    `json:"token"`
	OriginalName             *string    `json:"original_name"`
	Type                     *string    `json:"type"`
	CreatedBYID              *uint      `json:"created_by_id"`
	UpdatedBYID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
}

type RiskUpdatePlanRiskAssessmentWorkshopDocumentDTO struct {
	ID                       uint       `json:"id"`
	RiskUpdatePlanTimelineID uint       `json:"risk_update_plan_timeline_id"`
	URL                      *string    `json:"url"`
	Name                     *string    `json:"name"`
	Size                     *int       `json:"size"`
	Token                    *string    `json:"token"`
	OriginalName             *string    `json:"original_name"`
	Type                     *string    `json:"type"`
	CreatedBYID              *uint      `json:"created_by_id"`
	UpdatedBYID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
}

type RiskUpdatePlanRiskTeamAssessmentDTO struct {
	ID                       uint                             `json:"id"`
	RiskUpdatePlanTimelineID uint                             `json:"risk_update_plan_timeline_id"`
	RiskChampionID           uint                             `json:"risk_champion_id"`
	RiskChampion             *organizationDTOs.UserMinimalDTO `json:"risk_champion"`
	CreatedBYID              *uint                            `json:"created_by_id"`
	UpdatedBYID              *uint                            `json:"updated_by_id"`
	CreatedAt                time.Time                        `json:"created_at"`
	UpdatedAt                *time.Time                       `json:"updated_at"`
}

type RiskUpdatePlanRiskTeamDocumentDTO struct {
	ID                       uint       `json:"id"`
	RiskUpdatePlanTimelineID uint       `json:"risk_update_plan_timeline_id"`
	URL                      *string    `json:"url"`
	Name                     *string    `json:"name"`
	Size                     *int       `json:"size"`
	Token                    *string    `json:"token"`
	OriginalName             *string    `json:"original_name"`
	Type                     *string    `json:"type"`
	CreatedBYID              *uint      `json:"created_by_id"`
	UpdatedBYID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
}

type RiskUpdatePlanRiskWorkshopAssessmentDocumentDTO struct {
	ID                       uint       `json:"id"`
	RiskUpdatePlanTimelineID uint       `json:"risk_update_plan_timeline_id"`
	URL                      *string    `json:"url"`
	Name                     *string    `json:"name"`
	Size                     *int       `json:"size"`
	Token                    *string    `json:"token"`
	OriginalName             *string    `json:"original_name"`
	Type                     *string    `json:"type"`
	CreatedBYID              *uint      `json:"created_by_id"`
	UpdatedBYID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
}

type RiskUpdatePlanSelfAwarenessTrainingDTO struct {
	ID                       uint       `json:"id"`
	RiskUpdatePlanTimelineID uint       `json:"risk_update_plan_timeline_id"`
	Title                    string     `json:"title"`
	Description              *string    `json:"description"`
	CreatedBYID              *uint      `json:"created_by_id"`
	UpdatedBYID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
}

//COMMON DTOS

type CommonListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	IsTested *bool
}

type CommonDocumentFormDTO struct {
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"required"`
}

type CommonDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         *string `json:"name"`
	Token        *string `json:"token"`
	OriginalName *string `json:"original_name"`
	Type         *string `json:"type"`
	Size         *int    `json:"size"`
}

type RiskMitigationUpdatePlanFormDTO struct {
	RiskRegisterID         uint  `json:"risk_register_id" validate:"required"`
	RiskResponseStrategyID *uint `json:"risk_response_strategy_id"`
}
