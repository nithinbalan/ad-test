package dtos

import "time"

type PaginationInfoDTO struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int   `json:"total_pages"`
	TotalCount  int64 `json:"total_count"`
}

type IDResponse struct {
	ID uint `json:"id"`
}

type SuccessfulIDsResponse struct {
	SuccessfulIds []uint `json:"successfull_ids"`
}

type FailedIDsResponse struct {
	FailedIds []uint `json:"failed_ids"`
}

type FileDTO struct {
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	FileType     string `json:"type"`
	SizeInBytes  int64  `json:"size"`
}

type FilePreviewDTO struct {
	URL   string `json:"url" binding:"required"`
	Token string `json:"token" binding:"required"`
}

type DocumentFormDTO struct {
	ID            *uint   `json:"id,omitempty"`
	URL           *string `json:"url"`
	FileName      string  `json:"file_name"`
	OriginalName  string  `json:"original_name"`
	Type          string  `json:"type"`
	Size          int     `json:"size"`
	IsNewFile     *bool   `json:"is_new_file,omitempty"` // Optional field for indicating if the file is new or updated
	IsDeletedFile *bool   `json:"is_deleted_file,omitempty"`
}

type TaskFormDTO struct {
	ID                              *uint      `json:"id,omitempty"`
	ModuleIdentifier                string     `json:"module_identifier" binding:"required"`
	ComplianceActionPlanID          *uint      `json:"compliance_action_plan_id,omitempty"`
	SLAActionPlanID                 *uint      `json:"sla_action_plan_id,omitempty"`
	IncidentActionPlanID            *uint      `json:"incident_action_plan_id,omitempty"`
	ControlActionPlanID             *uint      `json:"control_action_plan_id,omitempty"`
	RiskActionPlanID                *uint      `json:"risk_action_plan_id,omitempty"`
	MeetingActionPlanID             *uint      `json:"meeting_action_plan_id,omitempty"`
	AuditFindingActionPlanID        *uint      `json:"audit_finding_action_plan_id,omitempty"`
	ProjectActionPlanID             *uint      `json:"project_action_plan_id,omitempty"`
	DocumentID                      *uint      `json:"document_id,omitempty"`
	DocumentWorkflowID              *uint      `json:"document_workflow_id,omitempty"`
	DocumentChangeRequestID         *uint      `json:"document_change_request_id,omitempty"`
	DocumentChangeRequestWorkflowID *uint      `json:"document_change_request_workflow_id,omitempty"`
	AuditDocumentRequestID          *uint      `json:"audit_document_request_id,omitempty"`
	AuditFinalReportID              *uint      `json:"audit_final_report_id,omitempty"`
	AuditDraftReportID              *uint      `json:"audit_draft_report_id,omitempty"`
	AuditReportWorkflowID           *uint      `json:"audit_report_workflow_id,omitempty"`
	RecoveryActionPlanID            *uint      `json:"recovery_action_plan_id,omitempty"`
	StakeholderActionPlanID         *uint      `json:"stakeholder_action_plan_id,omitempty"`
	StrategyProfileActionPlanID     *uint      `json:"strategy_profile_action_plan_id,omitempty"`
	ThreatActionPlanID              *uint      `json:"threat_action_plan_id,omitempty"`
	VulnerabilityActionPlanID       *uint      `json:"vulnerability_action_plan_id,omitempty"`
	Title                           string     `json:"title" binding:"required"`
	Description                     *string    `json:"description,omitempty"`
	TargetDate                      *time.Time `json:"target_date,omitempty"`
	URLParam                        *string    `json:"url_param,omitempty"`
	QueryParam                      *string    `json:"query_param,omitempty"`
	SeverityLevelID                 *uint      `json:"severity_level_id,omitempty"`
	ResponsibleUserID               uint       `json:"responsible_user_id" binding:"required"`
	TaskStatusID                    *uint      `json:"task_status_id"`
	ViolationReportActionPlanID     *uint      `json:"violation_report_action_plan_id,omitempty"`
}

type AssociatedControlListDTO struct {
	ID                                   uint       `json:"id"`
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
	ControlEffectivenessRatingTitleEn    string     `json:"control_effectiveness_rating_title_en"`
	ControlEffectivenessRatingTitleAr    string     `json:"control_effectiveness_rating_title_ar"`
	ControlEffectivenessRatingColorCode  string     `json:"control_effectiveness_rating_color_code"`
	ControlEffectivenessRatingType       string     `json:"control_effectiveness_rating_type"`
	ControlEffectivenessScore            *float64   `json:"control_effectiveness_score"`
	RiskReduction                        *float64   `json:"risk_reduction"`
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
	ControlStatusTitleEn                 string     `json:"control_status_title_en"`
	ControlStatusTitleAr                 string     `json:"control_status_title_ar"`
	ControlStatusColorCode               string     `json:"control_status_color_code"`
	ControlStatusType                    string     `json:"control_status_type"`
	ControlTestDescription               string     `json:"control_test_description"`
}

type AssociatedRiskListDTO struct {
	ID                              uint    `json:"id"`
	RiskRegisterID                  uint    `json:"risk_register_id"`
	ReferenceNumber                 string  `json:"reference_number"`
	RiskCategoryTitleEn             string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr             string  `json:"risk_category_title_ar"`
	RiskSourceTitleEn               *string `json:"risk_source_title_en"`
	RiskSourceTitleAr               *string `json:"risk_source_title_ar"`
	ResponsibleUnitID               uint    `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn          string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string  `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	RiskStatusTitleEn               string  `json:"risk_status_title_en"`
	RiskStatusTitleAr               string  `json:"risk_status_title_ar"`
	RiskResponseStrategyTitleEn     string  `json:"risk_response_strategy_title_en"`
	RiskResponseStrategyTitleAr     string  `json:"risk_response_strategy_title_ar"`
	RiskStatusColorCode             string  `json:"risk_status_color_code"`
	Title                           string  `json:"title"`
	Description                     *string `json:"description"`
	ReportedAt                      *string `json:"reported_at"`
	ReviewFrequencyTitleEn          *string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr          *string `json:"review_frequency_title_ar"`
	NextReviewDate                  *string `json:"next_review_date"`
	CurrentRiskRatingValueID        *uint   `json:"current_risk_rating_value_id"`
	Comments                        *string `json:"comments"`
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
	AppetiteAlignment               *string `json:"appetite_alignment"`
	ProcessID                       *uint   `json:"process_id"`
	ComplianceRegisterID            *uint   `json:"compliance_register_id"`
	ControlID                       *uint   `json:"control_id"`
	DocumentID                      *uint   `json:"document_id"`
	IncidentRegisterID              *uint   `json:"incident_register_id"`
	StrategicObjectiveId            *uint   `json:"strategic_objective_id"`
	KPIRegisterID                   *uint   `json:"kpi_register_id"`
}

type AssociatedRiskListResponseDTO struct {
	Items      []*AssociatedRiskListDTO `json:"items"`
	Pagination PaginationInfoDTO        `json:"pagination"`
}

type AssociatedClauseDTO struct {
	ComplianceClauseID uint   `json:"compliance_clause_id"`
	ClauseNo           string `json:"clause_no"`
	ClauseTitleEn      string `json:"clause_title_en"`
	ClauseTitleAr      string `json:"clause_title_ar"`
}

type AssociatedClauseListDTO struct {
	ComplianceRegisterID      uint                  `json:"compliance_register_id"`
	ComplianceRegisterTitleEn string                `json:"compliance_register_title_en"`
	ComplianceRegisterTitleAr string                `json:"compliance_register_title_ar"`
	Clauses                   []AssociatedClauseDTO `json:"clauses"`
}
type AssociatedClauseRowDTO struct {
	ComplianceRegisterID      uint
	ComplianceRegisterTitleEn string
	ComplianceRegisterTitleAr string
	ComplianceClauseID        uint
	ClauseNo                  string
	ClauseTitleEn             string
	ClauseTitleAr             string
}
type AssociatedItemClauseListResponseDTO struct {
	Items      []*AssociatedClauseListDTO `json:"items"`
	Pagination PaginationInfoDTO          `json:"pagination"`
}
type AssociatedClauseFormDTO struct {
	ClauseIDs []uint `json:"clause_ids" binding:"required"`
}
type AssociatedItemClauseListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
}
type AssociatedRiskMapDTO struct {
	RiskIDs []uint `json:"risk_ids" binding:"required"`
}

type AssociatedDocumentMapDTO struct {
	DocumentIDs []uint `json:"document_ids" binding:"required"`
}

type AssociatedIncidentMapDTO struct {
	IncidentIDs []uint `json:"incident_ids" binding:"required"`
}

type AssociatedStrategicObjectiveMapDTO struct {
	StrategicObjectiveIDs []uint `json:"strategic_objective_ids" binding:"required"`
}

type AssociatedAssetMapDTO struct {
	AssetIDs []uint `json:"asset_register_ids" binding:"required"`
}

type AssociatedProcessMapDTO struct {
	ProcessIDs []uint `json:"process_ids" binding:"required"`
}

type AssociatedComplianceMapDTO struct {
	ClauseIDs []uint `json:"clause_ids" binding:"required"`
}

type AssociatedComplianceRegisterMapDTO struct {
	ComplianceRegisterIDs []uint `json:"compliance_register_ids" binding:"required"`
}

type AssociatedDocumentListDTO struct {
	ID                              uint       `json:"id"`
	DocumentID                      uint       `json:"document_id"`
	DocumentTypeID                  uint       `json:"document_type_id"`
	TitleEn                         string     `json:"title_en"`
	TitleAr                         string     `json:"title_ar"`
	DocumentCategoryTitleEn         string     `json:"document_category_title_en"`
	DocumentCategoryTitleAr         string     `json:"document_category_title_ar"`
	DocumentTypeTitleEn             string     `json:"document_type_title_en"`
	DocumentTypeTitleAr             string     `json:"document_type_title_ar"`
	DocumentAccessTypeTitleEn       string     `json:"document_access_type_title_en"`
	DocumentAccessTypeTitleAr       string     `json:"document_access_type_title_ar"`
	ResponsibleUnitTitleEn          string     `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string     `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint       `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string     `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string     `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string     `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string     `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string     `json:"responsible_user_initial"`
	ResponsibleUserColor            string     `json:"responsible_user_color"`
	ResponsibleUserEmail            string     `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string     `json:"responsible_user_profile_image_path"`
	DocumentStatusTitleEn           string     `json:"document_status_title_en"`
	DocumentStatusTitleAr           string     `json:"document_status_title_ar"`
	DocumentStatusColorCode         string     `json:"document_status_color_code"`
	DocumentStatusType              string     `json:"document_status_type"`
	DocumentVersion                 float64    `json:"document_version"`
	DocumentURL                     string     `json:"document_url"`
	DocumentName                    string     `json:"document_name"`
	DocumentSize                    int        `json:"document_size"`
	DocumentToken                   string     `json:"document_token"`
	DocumentOriginalName            *string    `json:"document_original_name"`
	DocumentType                    string     `json:"document_type"`
	CreatedUserID                   uint       `json:"created_user_id"`
	CreatedUserFirstNameEn          string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string     `json:"created_user_last_name_ar"`
	CreatedUserInitial              string     `json:"created_user_initial"`
	CreatedUserColor                string     `json:"created_user_color"`
	CreatedUserEmail                string     `json:"created_user_email"`
	CreatedUserProfileImagePath     string     `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn          string     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          string     `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           string     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           string     `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string     `json:"updated_user_initial"`
	UpdatedUserColor                string     `json:"updated_user_color"`
	UpdatedUserEmail                string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string     `json:"updated_user_profile_image_path"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedAt                       *time.Time `json:"updated_at"`
	Acknowledged                    *bool      `json:"acknowledged,omitempty"`
	PublishedUserFirstNameEn        string     `json:"published_user_first_name_en"`
	PublishedUserFirstNameAr        string     `json:"published_user_first_name_ar"`
	PublishedUserLastNameEn         string     `json:"published_user_last_name_en"`
	PublishedUserLastNameAr         string     `json:"published_user_last_name_ar"`
	PublishedUserInitial            string     `json:"published_user_initial"`
	PublishedUserColor              string     `json:"published_user_color"`
	PublishedUserEmail              string     `json:"published_user_email"`
	PublishedUserProfileImagePath   string     `json:"published_user_profile_image_path"`
	PublishedAt                     *time.Time `json:"published_at"`
	UserAccessTypes                 *string    `json:"user_access_types"`
}

type AssociatedDocumentListResponseDTO struct {
	Items      []*AssociatedDocumentListDTO `json:"items"`
	Pagination PaginationInfoDTO            `json:"pagination"`
}

type AssociatedIncidentListDTO struct {
	ID                              uint       `json:"id"`
	IncidentRegisterID              uint       `json:"incident_register_id"`
	ReferenceNumber                 string     `json:"reference_number"`
	Title                           string     `json:"title"`
	Description                     *string    `json:"description"`
	IncidentAt                      time.Time  `json:"incident_at"`
	Location                        *string    `json:"location"`
	ReportedAt                      *time.Time `json:"reported_at"`
	QuickAction                     *string    `json:"quick_action"`
	IncidentCategoryID              uint       `json:"incident_category_id"`
	IncidentTypeID                  uint       `json:"incident_type_id"`
	SeverityLevelID                 uint       `json:"severity_level_id"`
	IncidentStatusID                uint       `json:"incident_status_id"`
	ResponsibleUnitTitleEn          string     `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string     `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint       `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string     `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string     `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string     `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string     `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string     `json:"responsible_user_initial"`
	ResponsibleUserColor            string     `json:"responsible_user_color"`
	ResponsibleUserEmail            string     `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string     `json:"responsible_user_profile_image_path"`
	IncidentCategoryTitleEn         string     `json:"incident_category_title_en"`
	IncidentCategoryTitleAr         string     `json:"incident_category_title_ar"`
	IncidentTypeTitleEn             string     `json:"incident_type_title_en"`
	IncidentTypeTitleAr             string     `json:"incident_type_title_ar"`
	SeverityLevelTitleEn            string     `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string     `json:"severity_level_title_ar"`
	IncidentStatusTitleEn           string     `json:"incident_status_title_en"`
	IncidentStatusTitleAr           string     `json:"incident_status_title_ar"`
	IncidentStatusColorCode         string     `json:"incident_status_color_code"`
	SeverityLevelColorCode          string     `json:"severity_level_color_code"`
	ReportedUserID                  uint       `json:"reported_user_id"`
	ReportedUserFirstNameEn         string     `json:"reported_user_first_name_en"`
	ReportedUserFirstNameAr         string     `json:"reported_user_first_name_ar"`
	ReportedUserLastNameEn          string     `json:"reported_user_last_name_en"`
	ReportedUserLastNameAr          string     `json:"reported_user_last_name_ar"`
	ReportedUserInitial             string     `json:"reported_user_initial"`
	ReportedUserColor               string     `json:"reported_user_color"`
	ReportedUserEmail               string     `json:"reported_user_email"`
	ReportedUserProfileImagePath    string     `json:"reported_user_profile_image_path"`
	InvestigationInitiatedAt        *string    `json:"investigation_initiated_at"`
	InvestigationDescription        *string    `json:"investigation_description"`
	CreatedAt                       string     `json:"created_at"`
	UpdatedAt                       string     `json:"updated_at"`
	CreatedUserID                   uint       `json:"created_user_id"`
	CreatedUserFirstNameEn          string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string     `json:"created_user_last_name_ar"`
	CreatedUserInitial              string     `json:"created_user_initial"`
	CreatedUserColor                string     `json:"created_user_color"`
	CreatedUserEmail                string     `json:"created_user_email"`
	CreatedUserProfileImagePath     string     `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn          string     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          string     `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           string     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           string     `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string     `json:"updated_user_initial"`
	UpdatedUserColor                string     `json:"updated_user_color"`
	UpdatedUserEmail                string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string     `json:"updated_user_profile_image_path"`
}

type AssociatedIncidentListResponseDTO struct {
	Items      []*AssociatedIncidentListDTO `json:"items"`
	Pagination PaginationInfoDTO            `json:"pagination"`
}

type AssociatedStrategicObjectiveListDTO struct {
	ID                   uint    `json:"id"`
	StrategicObjectiveID uint    `json:"strategic_objective_id"`
	TitleEn              string  `json:"title_en"`
	TitleAr              string  `json:"title_ar"`
	ReferenceNumber      string  `json:"reference_number"`
	DescriptionEn        *string `json:"description_en"`
	DescriptionAr        *string `json:"description_ar"`
	StartDate            *string `json:"start_date"`
	EndDate              *string `json:"end_date"`
	Status               string  `json:"status"`

	StrategicObjectiveCategoryID      uint   `json:"strategic_objective_category_id"`
	StrategicObjectiveCategoryTitleEn string `json:"strategic_objective_category_title_en"`
	StrategicObjectiveCategoryTitleAr string `json:"strategic_objective_category_title_ar"`

	ResponsibleUserID               *uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      *string `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn       *string `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr      *string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr       *string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string `json:"responsible_user_initial"`
	ResponsibleUserEmail            *string `json:"responsible_user_email"`
	ResponsibleUserColor            *string `json:"responsible_user_color"`
	ResponsibleUserProfileImagePath *string `json:"responsible_user_profile_image_path"`

	ResponsibleUnitID      *uint   `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn *string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr *string `json:"responsible_unit_title_ar"`

	ReviewUserID          *uint   `json:"review_user_id"`
	ReviewUserFirstNameEn *string `json:"review_user_first_name_en"`
	ReviewUserLastNameEn  *string `json:"review_user_last_name_en"`
	ReviewUserFirstNameAr *string `json:"review_user_first_name_ar"`
	ReviewUserLastNameAr  *string `json:"review_user_last_name_ar"`

	ReviewFrequencyID      *uint   `json:"review_frequency_id"`
	ReviewFrequencyTitleEn *string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr *string `json:"review_frequency_title_ar"`

	CreatedByID                 uint   `json:"created_by_id"`
	CreatedUserFirstNameEn      string `json:"created_user_first_name_en"`
	CreatedUserLastNameEn       string `json:"created_user_last_name_en"`
	CreatedUserFirstNameAr      string `json:"created_user_first_name_ar"`
	CreatedUserLastNameAr       string `json:"created_user_last_name_ar"`
	CreatedUserInitial          string `json:"created_user_initial"`
	CreatedUserColor            string `json:"created_user_color"`
	CreatedUserEmail            string `json:"created_user_email"`
	CreatedUserProfileImagePath string `json:"created_user_profile_image_path"`

	UpdatedByID                 *uint   `json:"updated_by_id"`
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
}

type AssociatedItemStrategicObjectiveListResponseDTO struct {
	Items      []*AssociatedStrategicObjectiveListDTO `json:"items"`
	Pagination PaginationInfoDTO                      `json:"pagination"`
}

type AssociatedAssetListDTO struct {
	ID                                         uint       `json:"id"`
	AssetRegisterID                            uint       `json:"asset_register_id"`
	ReferenceNumber                            string     `json:"reference_number"`
	TitleEn                                    string     `json:"title_en"`
	TitleAr                                    string     `json:"title_ar"`
	AssetCategoryTitleEn                       *string    `json:"asset_category_title_en"`
	AssetCategoryTitleAr                       *string    `json:"asset_category_title_ar"`
	AssetSubCategoryTitleEn                    *string    `json:"asset_sub_category_title_en"`
	AssetSubCategoryTitleAr                    *string    `json:"asset_sub_category_title_ar"`
	AssetStatusTitleEn                         string     `json:"asset_status_title_en"`
	AssetStatusTitleAr                         string     `json:"asset_status_title_ar"`
	AssetStatusColorCode                       string     `json:"asset_status_color_code"`
	ResponsibleUnitTitleEn                     string     `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr                     string     `json:"responsible_unit_title_ar"`
	ResponsibleUserID                          uint       `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn                 string     `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr                 string     `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn                  string     `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr                  string     `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial                     string     `json:"responsible_user_initial"`
	ResponsibleUserColor                       string     `json:"responsible_user_color"`
	ResponsibleUserEmail                       string     `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath            string     `json:"responsible_user_profile_image_path"`
	Location                                   *string    `json:"location"`
	Description                                *string    `json:"description"`
	AcquisitionDate                            string     `json:"acquisition_date"`
	AssignmentDate                             *string    `json:"assignment_date"`
	WarrantyExpiryDate                         *string    `json:"warranty_expiry_date"`
	Manufacturer                               *string    `json:"manufacturer"`
	Model                                      *string    `json:"model"`
	SerialNumber                               *string    `json:"serial_number"`
	TotalScore                                 *float64   `json:"total_score"`
	AssetCriticalityMatrixRatingScaleTitleEn   string     `json:"asset_criticality_matrix_rating_scale_title_en"`
	AssetCriticalityMatrixRatingScaleTitleAr   string     `json:"asset_criticality_matrix_rating_scale_title_ar"`
	AssetCriticalityMatrixRatingScaleColorCode string     `json:"asset_criticality_matrix_rating_scale_color_code"`
	CreatedUserID                              uint       `json:"created_user_id"`
	CreatedUserFirstNameEn                     string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr                     string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                      string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                      string     `json:"created_user_last_name_ar"`
	CreatedUserInitial                         string     `json:"created_user_initial"`
	CreatedUserColor                           string     `json:"created_user_color"`
	CreatedUserEmail                           string     `json:"created_user_email"`
	CreatedUserProfileImagePath                string     `json:"created_user_profile_image_path"`
	CreatedAt                                  time.Time  `json:"created_at"`
	UpdatedUserID                              *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn                     *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr                     *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                      *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                      *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                         string     `json:"updated_user_initial"`
	UpdatedUserColor                           string     `json:"updated_user_color"`
	UpdatedUserEmail                           string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath                string     `json:"updated_user_profile_image_path"`
	UpdatedAt                                  *time.Time `json:"updated_at"`
}

type AssociatedItemAssetListResponseDTO struct {
	Items      []*AssociatedAssetListDTO `json:"items"`
	Pagination PaginationInfoDTO         `json:"pagination"`
}

type AssociatedProcessListDTO struct {
	ID                               uint       `json:"id"`
	ProcessID                        uint       `json:"process_id"`
	ProcessGroupTitleEn              string     `json:"process_group_title_en"`
	ProcessGroupTitleAr              string     `json:"process_group_title_ar"`
	OwnerUnitTitleEn                 *string    `json:"owner_unit_title_en"`
	OwnerUnitTitleAr                 *string    `json:"owner_unit_title_ar"`
	OwnerID                          *uint      `json:"owner_id"`
	OwnerFirstNameEn                 *string    `json:"owner_first_name_en"`
	OwnerFirstNameAr                 *string    `json:"owner_first_name_ar"`
	OwnerLastNameEn                  *string    `json:"owner_last_name_en"`
	OwnerLastNameAr                  *string    `json:"owner_last_name_ar"`
	OwnerInitial                     string     `json:"owner_initial"`
	OwnerColor                       string     `json:"owner_color"`
	OwnerEmail                       string     `json:"owner_email"`
	OwnerProfileImagePath            string     `json:"owner_profile_image_path"`
	LifecycleStatusTitleEn           string     `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr           string     `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode         string     `json:"lifecycle_status_color_code"`
	ReferenceNumber                  string     `json:"reference_number"`
	TitleEn                          string     `json:"title_en"`
	TitleAr                          string     `json:"title_ar"`
	Description                      *string    `json:"description"`
	MarkedAsAuditableItem            bool       `json:"marked_as_auditable_item"`
	RtoFromValue                     *int       `json:"rto_from_value"`
	RtoToValue                       *int       `json:"rto_to_value"`
	RtoUnitId                        *uint      `json:"rto_unit_id"`
	RtoUnitTitleEn                   *string    `json:"rto_unit_title_en"`
	RtoUnitTitleAr                   *string    `json:"rto_unit_title_ar"`
	RtoCondition                     *string    `json:"rto_condition"`
	MtpdFromValue                    *int       `json:"mtpd_from_value"`
	MtpdToValue                      *int       `json:"mtpd_to_value"`
	MtpdUnitId                       *uint      `json:"mtpd_unit_id"`
	MtpdUnitTitleEn                  *string    `json:"mtpd_unit_title_en"`
	MtpdUnitTitleAr                  *string    `json:"mtpd_unit_title_ar"`
	MtpdCondition                    *string    `json:"mtpd_condition"`
	RpoFromValue                     *int       `json:"rpo_from_value"`
	RpoToValue                       *int       `json:"rpo_to_value"`
	RpoUnitId                        *uint      `json:"rpo_unit_id"`
	RpoUnitTitleEn                   *string    `json:"rpo_unit_title_en"`
	RpoUnitTitleAr                   *string    `json:"rpo_unit_title_ar"`
	RpoCondition                     *string    `json:"rpo_condition"`
	RspFromValue                     *int       `json:"rsp_from_value"`
	RspToValue                       *int       `json:"rsp_to_value"`
	RspUnitId                        *uint      `json:"rsp_unit_id"`
	RspUnitTitleEn                   *string    `json:"rsp_unit_title_en"`
	RspUnitTitleAr                   *string    `json:"rsp_unit_title_ar"`
	RspCondition                     *string    `json:"rsp_condition"`
	ImpactAnalysisLevelId            *uint      `json:"impact_analysis_level_id"`
	ImpactAnalysisLevelTitleEn       *string    `json:"impact_analysis_level_title_en"`
	ImpactAnalysisLevelTitleAr       *string    `json:"impact_analysis_level_title_ar"`
	ImpactAnalysisLevelColorCode     *string    `json:"impact_analysis_level_color_code"`
	ImpactAnalysisLevelType          *string    `json:"impact_analysis_level_type"`
	RiskLevelId                      *uint      `json:"risk_level_id"`
	RiskLevelTitleEn                 *string    `json:"risk_level_title_en"`
	RiskLevelTitleAr                 *string    `json:"risk_level_title_ar"`
	RiskLevelColorCode               *string    `json:"risk_level_color_code"`
	RiskLevelType                    *string    `json:"risk_level_type"`
	LastAssessedAt                   *time.Time `json:"last_assessed_at"`
	LastAssessedUserID               *uint      `json:"last_assessed_user_id"`
	LastAssessedUserFirstNameEn      *string    `json:"last_assessed_user_first_name_en"`
	LastAssessedUserFirstNameAr      *string    `json:"last_assessed_user_first_name_ar"`
	LastAssessedUserLastNameEn       *string    `json:"last_assessed_user_last_name_en"`
	LastAssessedUserLastNameAr       *string    `json:"last_assessed_user_last_name_ar"`
	LastAssessedUserInitial          *string    `json:"last_assessed_user_initial"`
	LastAssessedUserColor            *string    `json:"last_assessed_user_color"`
	LastAssessedUserEmail            *string    `json:"last_assessed_user_email"`
	LastAssessedUserProfileImagePath *string    `json:"last_assessed_user_profile_image_path"`
	IsBiaPerformed                   bool       `json:"is_bia_performed"`
	ProcessCriticalityLevelID        uint       `json:"process_criticality_level_id"`
	ProcessCriticalityLevelTitleEn   string     `json:"process_criticality_level_title_en"`
	ProcessCriticalityLevelTitleAr   string     `json:"process_criticality_level_title_ar"`
	CreatedUserID                    uint       `json:"created_user_id"`
	CreatedUserFirstNameEn           string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr           string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn            string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr            string     `json:"created_user_last_name_ar"`
	CreatedUserInitial               string     `json:"created_user_initial"`
	CreatedUserColor                 string     `json:"created_user_color"`
	CreatedUserEmail                 string     `json:"created_user_email"`
	CreatedUserProfileImagePath      string     `json:"created_user_profile_image_path"`
	UpdatedUserID                    *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn           *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr           *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn            *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr            *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial               string     `json:"updated_user_initial"`
	UpdatedUserColor                 string     `json:"updated_user_color"`
	UpdatedUserEmail                 string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath      string     `json:"updated_user_profile_image_path"`
	CreatedAt                        time.Time  `json:"created_at"`
	UpdatedAt                        *time.Time `json:"updated_at"`
	RiskRatingTitleEn                *string    `json:"risk_rating_title_en"`
	RiskRatingTitleAr                *string    `json:"risk_rating_title_ar"`
	RiskRatingColorCode              *string    `json:"risk_rating_color_code"`
	RiskRatingValue                  *uint      `json:"risk_rating_value"`
}

type AssociatedItemProcessListResponseDTO struct {
	Items      []*AssociatedProcessListDTO `json:"items"`
	Pagination PaginationInfoDTO           `json:"pagination"`
}
type AssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Value       int
	QuickFilter []string
	Filters     map[string][]string
}

type AssociatedKPIMapDTO struct {
	KpiIDs []uint `json:"kpi_ids" binding:"required"`
}

type AssociatedControlMapDTO struct {
	ControlIDs []uint `json:"control_ids" binding:"required"`
}

type AssociatedKPIRegisterListDTO struct {
	ID                              int      `json:"id"`
	KPIRegisterID                   uint     `json:"kpi_register_id"`
	ReferenceNumber                 *string  `json:"reference_number"`
	Title                           string   `json:"title"`
	Description                     *string  `json:"description"`
	KPICategoryID                   uint     `json:"kpi_category_id"`
	KPICategoryTitleEn              string   `json:"kpi_category_title_en"`
	KPICategoryTitleAr              string   `json:"kpi_category_title_ar"`
	ResponsibleUserID               *uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      *string  `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn       *string  `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr      *string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr       *string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string  `json:"responsible_user_initial"`
	ResponsibleUserEmail            *string  `json:"responsible_user_email"`
	ResponsibleUserColor            *string  `json:"responsible_user_color"`
	ResponsibleUserProfileImagePath *string  `json:"responsible_user_profile_image_path"`
	KPIStatusID                     uint     `json:"kpi_status_id"`
	KPIStatusTitleEn                string   `json:"kpi_status_title_en"`
	KPIStatusTitleAr                string   `json:"kpi_status_title_ar"`
	KPIStatusColorCode              string   `json:"kpi_status_color_code"`
	KPIStatusType                   string   `json:"kpi_status_type"`
	ReviewFrequencyTitleEn          *string  `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr          *string  `json:"review_frequency_title_ar"`
	ReviewFrequencyType             *string  `json:"review_frequency_type"`
	CalculationTypeTitleEn          *string  `json:"calculation_type_title_en"`
	CalculationTypeTitleAr          *string  `json:"calculation_type_title_ar"`
	CalculationTypeType             *string  `json:"calculation_type_type"`
	TargetValue                     *float64 `json:"target_value"`
	UnitTitleEn                     *string  `json:"unit_title_en"`
	UnitTitleAr                     *string  `json:"unit_title_ar"`
	UnitCodeEn                      *string  `json:"unit_code_en"`
	UnitCodeAr                      *string  `json:"unit_code_ar"`
	UnitColorCode                   *string  `json:"unit_color_code"`
	Weightage                       *float64 `json:"weightage"`
	StartDate                       *string  `json:"start_date"`
	EndDate                         *string  `json:"end_date"`
	Remarks                         *string  `json:"remarks"`
	CreatedUserID                   uint     `json:"created_user_id"`
	CreatedUserFirstNameEn          string   `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string   `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string   `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string   `json:"created_user_last_name_ar"`
	CreatedUserInitial              string   `json:"created_user_initial"`
	CreatedUserColor                string   `json:"created_user_color"`
	CreatedUserEmail                string   `json:"created_user_email"`
	CreatedUserProfileImagePath     string   `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint    `json:"updated_user_id"`
	UpdatedUserFirstNameEn          *string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              *string  `json:"updated_user_initial"`
	UpdatedUserColor                *string  `json:"updated_user_color"`
	UpdatedUserEmail                *string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath     *string  `json:"updated_user_profile_image_path"`
}
type AssociatedKPIRegisterListResponseDTO struct {
	Items      []*AssociatedKPIRegisterListDTO `json:"items"`
	Pagination PaginationInfoDTO               `json:"pagination"`
}

type AssociatedControlListResponseDTO struct {
	Items      []*AssociatedControlListDTO `json:"items"`
	Pagination PaginationInfoDTO           `json:"pagination"`
}

type AssociatedBCPMapDTO struct {
	BcpIDs []uint `json:"business_continuity_plan_ids" binding:"required"`
}

type AssociatedProjectMapDTO struct {
	ProjectIDs []uint `json:"project_ids" binding:"required"`
}

type AssociatedBusinessContinuityPlanListDTO struct {
	ID                              uint    `json:"id"`
	ReferenceNumber                 string  `json:"reference_number"`
	BusinessFunctionID              uint    `json:"business_function_id"`
	BusinessFunctionTitle           string  `json:"business_function_title"`
	Version                         *int    `json:"version"`
	IsLatestVersion                 *bool   `json:"is_latest_version"`
	Purpose                         *string `json:"purpose"`
	Scope                           *string `json:"scope"`
	Objectives                      *string `json:"objectives"`
	Assumptions                     *string `json:"assumptions"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	BcpStatusTitleEn                string  `json:"bcp_status_title_en"`
	BcpStatusTitleAr                string  `json:"bcp_status_title_ar"`
	BcpStatusColorCode              string  `json:"bcp_status_color_code"`
	BcpStatusType                   string  `json:"bcp_status_type"`
	Rto                             *int    `json:"rto"`
	RtoUnitId                       *uint   `json:"rto_unit_id"`
	RtoUnitTitleEn                  *string `json:"rto_unit_title_en"`
	RtoUnitTitleAr                  *string `json:"rto_unit_title_ar"`
	Mtpd                            *int    `json:"mtpd"`
	MtpdUnitId                      *uint   `json:"mtpd_unit_id"`
	MtpdUnitTitleEn                 *string `json:"mtpd_unit_title_en"`
	MtpdUnitTitleAr                 *string `json:"mtpd_unit_title_ar"`
	Rpo                             *int    `json:"rpo"`
	RpoUnitId                       *uint   `json:"rpo_unit_id"`
	RpoUnitTitleEn                  *string `json:"rpo_unit_title_en"`
	RpoUnitTitleAr                  *string `json:"rpo_unit_title_ar"`
	CreatedAt                       string  `json:"created_at"`
	UpdatedAt                       string  `json:"updated_at"`
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
	UpdatedUserProfileImagePath     *string `json:"updated_user_profile_image_path"`
}

type AssociatedBusinessContinuityPlanListResponseDTO struct {
	Items      []*AssociatedBusinessContinuityPlanListDTO `json:"items"`
	Pagination PaginationInfoDTO                          `json:"pagination"`
}

type AssociatedProjectListDTO struct {
	ID                              uint    `json:"id"`
	TitleEn                         string  `json:"title_en"`
	TitleAr                         string  `json:"title_ar"`
	DescriptionEn                   *string `json:"description_en"`
	DescriptionAr                   *string `json:"description_ar"`
	ProjectTypeTitleEn              string  `json:"project_type_title_en"`
	ProjectTypeTitleAr              string  `json:"project_type_title_ar"`
	ProjectStatusTitleEn            string  `json:"project_status_title_en"`
	ProjectStatusTitleAr            string  `json:"project_status_title_ar"`
	StartDate                       string  `json:"start_date"`
	EndDate                         string  `json:"end_date"`
	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
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
	UpdatedUserProfileImagePath     *string `json:"updated_user_profile_image_path"`
}

type AssociatedProjectListResponseDTO struct {
	Items      []*AssociatedProjectListDTO `json:"items"`
	Pagination PaginationInfoDTO           `json:"pagination"`
}

type AssociatedComplianceListDTO struct {
	ID                                    uint   `json:"id"`
	ComplianceRegisterID                  uint   `json:"compliance_register_id"`
	ComplianceCategoryTitleEn             string `json:"compliance_category_title_en"`
	ComplianceCategoryTitleAr             string `json:"compliance_category_title_ar"`
	ComplianceStatusTitleEn               string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr               string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode             string `json:"compliance_status_color_code"`
	ResponsibleUnitTitleEn                string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr                string `json:"responsible_unit_title_ar"`
	ResponsibleUserID                     uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn            string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr            string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn             string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr             string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial                string `json:"responsible_user_initial"`
	ResponsibleUserColor                  string `json:"responsible_user_color"`
	ResponsibleUserEmail                  string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath       string `json:"responsible_user_profile_image_path"`
	ReviewFrequencyTitleEn                string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr                string `json:"review_frequency_title_ar"`
	ReviewResponsibleUserID               uint   `json:"review_responsible_user_id"`
	ReviewResponsibleUserFirstNameEn      string `json:"review_responsible_user_first_name_en"`
	ReviewResponsibleUserFirstNameAr      string `json:"review_responsible_user_first_name_ar"`
	ReviewResponsibleUserLastNameEn       string `json:"review_responsible_user_last_name_en"`
	ReviewResponsibleUserLastNameAr       string `json:"review_responsible_user_last_name_ar"`
	ReviewResponsibleUserInitial          string `json:"review_responsible_user_initial"`
	ReviewResponsibleUserColor            string `json:"review_responsible_user_color"`
	ReviewResponsibleUserEmail            string `json:"review_responsible_user_email"`
	ReviewResponsibleUserProfileImagePath string `json:"review_responsible_user_profile_image_path"`
	SeverityLevelTitleEn                  string `json:"severity_level_title_en"`
	SeverityLevelTitleAr                  string `json:"severity_level_title_ar"`
	SeverityLevelColorCode                string `json:"severity_level_color_code"`
	ReferenceNumber                       string `json:"reference_number"`
	StandardCode                          string `json:"standard_code"`
	TitleEn                               string `json:"title_en"`
	TitleAr                               string `json:"title_ar"`
	DescriptionEn                         string `json:"description_en"`
	DescriptionAr                         string `json:"description_ar"`
	Version                               string `json:"version"`
	Deadline                              string `json:"deadline"`
	IssueDate                             string `json:"issue_date"`
	ExpiryDate                            string `json:"expiry_date"`
	NextReviewDate                        string `json:"next_review_date"`
	LastReviewedDate                      string `json:"last_reviewed_date"`
	CreatedAt                             string `json:"created_at"`
	UpdatedAt                             string `json:"updated_at"`
	CreatedUserID                         uint   `json:"created_user_id"`
	CreatedUserFirstNameEn                string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr                string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                 string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                 string `json:"created_user_last_name_ar"`
	CreatedUserInitial                    string `json:"created_user_initial"`
	CreatedUserColor                      string `json:"created_user_color"`
	CreatedUserEmail                      string `json:"created_user_email"`
	CreatedUserProfileImagePath           string `json:"created_user_profile_image_path"`
	UpdatedUserID                         *uint  `json:"updated_user_id"`
	UpdatedUserFirstNameEn                string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr                string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                 string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                 string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                    string `json:"updated_user_initial"`
	UpdatedUserColor                      string `json:"updated_user_color"`
	UpdatedUserEmail                      string `json:"updated_user_email"`
	UpdatedUserProfileImagePath           string `json:"updated_user_profile_image_path"`
	LifecycleStatusTitleEn                string `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr                string `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode              string `json:"lifecycle_status_color_code"`
	ComplianceAssessmentFrameworkID       *uint  `json:"compliance_assessment_framework_id"`
	ComplianceAssessmentFrameworkTitleEn  string `json:"compliance_assessment_framework_title_en"`
	ComplianceAssessmentFrameworkTitleAr  string `json:"compliance_assessment_framework_title_ar"`
}

type AssociatedComplianceListResponseDTO struct {
	Items      []*AssociatedComplianceListDTO `json:"items"`
	Pagination PaginationInfoDTO              `json:"pagination"`
}
