package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
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
	Export      bool
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
	WorkflowStatusID               *uint                                        `json:"workflow_status_id"`
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
	CreatedAt                    time.Time                      `json:"created_at"`
	UpdatedAt                    *time.Time                     `json:"updated_at"`
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

// Auto-fetch risks based on business unit (first call)
type AutoFetchRisksResponseDTO struct {
	TotalRisks  int    `json:"total_risks"`
	MappedCount int    `json:"mapped_count"`
	FailedCount int    `json:"failed_count"`
	FailedIDs   []uint `json:"failed_ids"`
}

// Manual mapping of additional risks (second call)
type MapRisksToUpdatePlanDTO struct {
	RiskRegisterIDs []uint `json:"risk_register_ids" binding:"required,min=1"`
}

type MapRisksToUpdatePlanResponseDTO struct {
	SuccessCount int    `json:"success_count"`
	FailedCount  int    `json:"failed_count"`
	FailedIDs    []uint `json:"failed_ids"`
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
	IsStrategic                     bool    `json:"is_strategic"`
	IsOperational                   bool    `json:"is_operational"`
	ExpectedCause                   *string `json:"expected_cause"`
	Consequence                     *string `json:"consequence"`
	IsInternal                      bool    `json:"is_internal"`
	IsExternal                      bool    `json:"is_external"`
}

type ERMDeptEvaluationListQueryParamsDTO struct {
	Sort                         string
	SortBy                       string
	Keyword                      string
	Page                         int
	PageSize                     int
	Value                        int
	Filters                      map[string][]string
	QuickFilter                  []string
	ReportedAtFrom               string
	ReportedAtTo                 string
	NextReviewDateFrom           string
	NextReviewDateTo             string
	Export                       bool
	DateFrom                     string
	DateTo                       string
	RiskAppetiteToleranceLevelID uint
}

type ERMDeptEvaluationListResponseDTO struct {
	Items      []*ERMDeptEvaluationRiskListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO    `json:"pagination"`
}

// Introductory Meeting Document DTOs
type IntroductoryMeetingDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         *string `json:"name"`
	Token        *string `json:"token"`
	OriginalName *string `json:"original_name"`
	Type         *string `json:"type"`
	Size         *int    `json:"size"`
}

type IntroductoryMeetingDocumentFormDTO struct {
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"required"`
}

// ERM Department Evaluation Document DTOs
type ERMDeptEvaluationDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         *string `json:"name"`
	Token        *string `json:"token"`
	OriginalName *string `json:"original_name"`
	Type         *string `json:"type"`
	Size         *int    `json:"size"`
}

type ERMDeptEvaluationDocumentFormDTO struct {
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"required"`
}

// Self-Awareness Training Document DTOs
type SelfAwarenessTrainingDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         *string `json:"name"`
	Token        *string `json:"token"`
	OriginalName *string `json:"original_name"`
	Type         *string `json:"type"`
	Size         *int    `json:"size"`
}

type SelfAwarenessTrainingDocumentFormDTO struct {
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"required"`
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
	Rto                              *uint                `json:"rto"`
	RtoUnitId                        *uint                `json:"rto_unit_id"`
	RtoUnitTitleEn                   *string              `json:"rto_unit_title_en"`
	RtoUnitTitleAr                   *string              `json:"rto_unit_title_ar"`
	Mtpd                             *uint                `json:"mtpd"`
	MtpdUnitId                       *uint                `json:"mtpd_unit_id"`
	MtpdUnitTitleEn                  *string              `json:"mtpd_unit_title_en"`
	MtpdUnitTitleAr                  *string              `json:"mtpd_unit_title_ar"`
	Rpo                              *uint                `json:"rpo"`
	RpoUnitId                        *uint                `json:"rpo_unit_id"`
	RpoUnitTitleEn                   *string              `json:"rpo_unit_title_en"`
	RpoUnitTitleAr                   *string              `json:"rpo_unit_title_ar"`
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
	ID                                   uint       `json:"id"`
	RiskUpdatePlanTimelineID             uint       `json:"risk_update_plan_timeline_id"`
	ControlID                            uint       `json:"control_id"`
	ControlTypeTitleEn                   string     `json:"control_type_title_en"`
	ControlTypeTitleAr                   string     `json:"control_type_title_ar"`
	ControlCategoryTitleEn               string     `json:"control_category_title_en"`
	ControlCategoryTitleAr               string     `json:"control_category_title_ar"`
	OwnerUnitTitleEn                     *string    `json:"owner_unit_title_en"`
	OwnerUnitTitleAr                     *string    `json:"owner_unit_title_ar"`
	OwnerID                              *uint      `json:"owner_id"`
	OwnerFirstNameEn                     *string    `json:"owner_first_name_en"`
	OwnerFirstNameAr                     *string    `json:"owner_first_name_ar"`
	OwnerLastNameEn                      *string    `json:"owner_last_name_en"`
	OwnerLastNameAr                      *string    `json:"owner_last_name_ar"`
	OwnerInitial                         string     `json:"owner_initial"`
	OwnerColor                           string     `json:"owner_color"`
	OwnerEmail                           string     `json:"owner_email"`
	OwnerProfileImagePath                string     `json:"owner_profile_image_path"`
	ControlImplementationStatusTitleEn   string     `json:"control_implementation_status_title_en"`
	ControlImplementationStatusTitleAr   string     `json:"control_implementation_status_title_ar"`
	ControlImplementationStatusColorCode string     `json:"control_implementation_status_color_code"`
	ControlEfficiencyStatusTitleEn       string     `json:"control_efficiency_status_title_en"`
	ControlEfficiencyStatusTitleAr       string     `json:"control_efficiency_status_title_ar"`
	ControlEfficiencyStatusColorCode     string     `json:"control_efficiency_status_color_code"`
	ReferenceNumber                      string     `json:"reference_number"`
	Code                                 *string    `json:"code"`
	TitleEn                              string     `json:"title_en"`
	TitleAr                              string     `json:"title_ar"`
	Description                          *string    `json:"description"`
	Objective                            *string    `json:"objective"`
	ReviewFrequencyTitleEn               string     `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr               string     `json:"review_frequency_title_ar"`
	LastReviewDate                       *time.Time `json:"last_review_date"`
	NextReviewDate                       *time.Time `json:"next_review_date"`
	LastReviewNote                       *string    `json:"last_review_note"`
	TestResponsibleUserID                uint       `json:"test_responsible_user_id"`
	TestResponsibleUserFirstNameEn       string     `json:"test_responsible_user_first_name_en"`
	TestResponsibleUserFirstNameAr       string     `json:"test_responsible_user_first_name_ar"`
	TestResponsibleUserLastNameEn        string     `json:"test_responsible_user_last_name_en"`
	TestResponsibleUserLastNameAr        string     `json:"test_responsible_user_last_name_ar"`
	TestResponsibleUserInitial           string     `json:"test_responsible_user_initial"`
	TestResponsibleUserColor             string     `json:"test_responsible_user_color"`
	TestResponsibleUserEmail             string     `json:"test_responsible_user_email"`
	TestResponsibleUserProfileImagePath  string     `json:"test_responsible_user_profile_image_path"`
	TestFrequencyTitleEn                 *string    `json:"test_frequency_title_en"`
	TestFrequencyTitleAr                 *string    `json:"test_frequency_title_ar"`
	LastTestDate                         *time.Time `json:"last_test_date"`
	NextTestDate                         *time.Time `json:"next_test_date"`
	ImplementedAt                        *time.Time `json:"implemented_at"`
	IsApplicable                         bool       `json:"is_applicable"`
	Justification                        *string    `json:"justification"`
	ImplementationGuideline              *string    `json:"implementation_guideline"`
	ApplicabilityChangeDate              *string    `json:"applicability_change_date"`
	LifecycleStatusTitleEn               string     `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr               string     `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode             string     `json:"lifecycle_status_color_code"`
	ControlStatusID                      *uint      `json:"control_status_id"`
	ControlStatusTitleEn                 *string    `json:"control_status_title_en"`
	ControlStatusTitleAr                 *string    `json:"control_status_title_ar"`
	ControlStatusColorCode               *string    `json:"control_status_color_code"`
	ControlStatusType                    *string    `json:"control_status_type"`
	CreatedUserID                        uint       `json:"created_user_id"`
	CreatedUserFirstNameEn               string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr               string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                string     `json:"created_user_last_name_ar"`
	CreatedUserInitial                   string     `json:"created_user_initial"`
	CreatedUserColor                     string     `json:"created_user_color"`
	CreatedUserEmail                     string     `json:"created_user_email"`
	CreatedUserProfileImagePath          string     `json:"created_user_profile_image_path"`
	UpdatedUserID                        *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn               *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr               *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                   string     `json:"updated_user_initial"`
	UpdatedUserColor                     string     `json:"updated_user_color"`
	UpdatedUserEmail                     string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath          string     `json:"updated_user_profile_image_path"`
	CreatedAt                            time.Time  `json:"created_at"`
	UpdatedAt                            *time.Time `json:"updated_at"`
	RiskRatingTitleEn                    *string    `json:"risk_rating_title_en"`
	RiskRatingTitleAr                    *string    `json:"risk_rating_title_ar"`
	RiskRatingColorCode                  *string    `json:"risk_rating_color_code"`
	RiskRatingValue                      *uint      `json:"risk_rating_value"`
	ControlAutomationID                  *uint      `json:"control_automation_id"`
	ControlAutomationTitleEn             *string    `json:"control_automation_title_en"`
	ControlAutomationTitleAr             *string    `json:"control_automation_title_ar"`
	ControlAutomationType                *string    `json:"control_automation_type"`
	ControlAutomationScore               *float64   `json:"control_automation_score"`
	ControlAutomationStatus              *string    `json:"control_automation_status"`
	ControlConsistencyID                 *uint      `json:"control_consistency_id"`
	ControlConsistencyTitleEn            *string    `json:"control_consistency_title_en"`
	ControlConsistencyTitleAr            *string    `json:"control_consistency_title_ar"`
	ControlConsistencyType               *string    `json:"control_consistency_type"`
	ControlConsistencyScore              *float64   `json:"control_consistency_score"`
	ControlConsistencyStatus             *string    `json:"control_consistency_status"`
	ControlConsistencyDescriptionEn      *string    `json:"control_consistency_description_en"`
	ControlConsistencyDescriptionAr      *string    `json:"control_consistency_description_ar"`
	ControlEffectivenessScore            *float64   `json:"control_effectiveness_score"`
	ControlEffectivenessRatingID         *uint      `json:"control_effectiveness_rating_id"`
	ControlEffectivenessRatingTitleEn    *string    `json:"control_effectiveness_rating_title_en"`
	ControlEffectivenessRatingTitleAr    *string    `json:"control_effectiveness_rating_title_ar"`
	ControlEffectivenessRatingMinScore   *float64   `json:"control_effectiveness_rating_min_score"`
	ControlEffectivenessRatingMaxScore   *float64   `json:"control_effectiveness_rating_max_score"`
	ControlEffectivenessRatingType       *string    `json:"control_effectiveness_rating_type"`
	RiskReduction                        *float64   `json:"risk_reduction"`
}

type ERMDeptEvaluationControlListResponseDTO struct {
	Items      []*ERMDeptEvaluationControlListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO       `json:"pagination"`
}

type MapControlsToUpdatePlanDTO struct {
	ControlIDs []uint `json:"control_ids" binding:"required,min=1"`
}

// Timeline Not Applicable DTO
type MarkTimelineNotApplicableDTO struct {
	Justification string `json:"justification" binding:"required"`
}

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

type LeadershipAssessmentDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         *string `json:"name"`
	Size         *int    `json:"size"`
	Token        *string `json:"token"`
	OriginalName *string `json:"original_name"`
	Type         *string `json:"type"`
}

type RiskTeamAssessmentUserListDTO struct {
	ID                       uint                             `json:"id"`
	RiskUpdatePlanTimelineID uint                             `json:"risk_update_plan_timeline_id"`
	RiskChampionID           uint                             `json:"risk_champion_id"`
	RiskChampion             *organizationDTOs.UserDetailsDTO `json:"risk_champion"`
	CreatedAt                time.Time                        `json:"created_at"`
	UpdatedAt                *time.Time                       `json:"updated_at"`
}

type RiskTeamAssessmentDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         *string `json:"name"`
	Size         *int    `json:"size"`
	Token        *string `json:"token"`
	OriginalName *string `json:"original_name"`
	Type         *string `json:"type"`
}

type UserListQueryParamsDTO struct {
	Sort        string              `json:"sort"`
	SortBy      string              `json:"sort_by"`
	Keyword     string              `json:"keyword"`
	Page        int                 `json:"page"`
	PageSize    int                 `json:"page_size"`
	QuickFilter []string            `json:"quick_filter"`
	DateFrom    string              `json:"date_from"`
	DateTo      string              `json:"date_to"`
	Filters     map[string][]string `json:"filters"`
	Export      bool                `json:"export"`
}
