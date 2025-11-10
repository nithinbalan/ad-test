package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type PolicyReportListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type PolicyReportListDTO struct {
	ID                          uint   `json:"id"`
	Title                       string `json:"title"`
	StartDate                   string `json:"start_date"`
	EndDate                     string `json:"end_date"`
	DocumentID                  uint   `json:"document_id"`
	DocumentTitleEn             string `json:"document_title_en"`
	DocumentTitleAr             string `json:"document_title_ar"`
	CreatedAt                   string `json:"created_at"`
	CreatedUserID               uint   `json:"created_user_id"`
	CreatedUserFirstNameEn      string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string `json:"created_user_last_name_ar"`
	CreatedUserInitial          string `json:"created_user_initial"`
	CreatedUserColor            string `json:"created_user_color"`
	CreatedUserEmail            string `json:"created_user_email"`
	CreatedUserProfileImagePath string `json:"created_user_profile_image_path"`
}

type PolicyReportLisResponseDTO struct {
	Items      []*PolicyReportListDTO       `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type PolicyReportFormDTO struct {
	Title      string `json:"title" binding:"required"`
	StartDate  string `json:"start_date" binding:"required"`
	EndDate    string `json:"end_date" binding:"required"`
	DocumentID uint   `json:"document_id" binding:"required"`
}

type PolicyReportParsedDates struct {
	StartDate *time.Time
	EndDate   *time.Time
}

type PolicyReportQueryParamsDTO struct {
	DocumentID uint
	StartDate  string
	EndDate    string
}

type PolicyReportPolicyOverviewDTO struct {
	ID                              uint       `json:"id"`
	DocumentTypeID                  uint       `json:"document_type_id"`
	ReferenceNumber                 string     `json:"reference_number"`
	TitleEn                         string     `json:"title_en"`
	TitleAr                         string     `json:"title_ar"`
	DescriptionEn                   string     `json:"description_en"`
	DescriptionAr                   string     `json:"description_ar"`
	DocumentCategoryTitleEn         string     `json:"document_category_title_en"`
	DocumentCategoryTitleAr         string     `json:"document_category_title_ar"`
	DocumentTypeTitleEn             string     `json:"document_type_title_en"`
	DocumentTypeTitleAr             string     `json:"document_type_title_ar"`
	DocumentAccessTypeTitleEn       string     `json:"document_access_type_title_en"`
	DocumentAccessTypeTitleAr       string     `json:"document_access_type_title_ar"`
	DocumentAccessTypeColorCode     string     `json:"document_access_type_color_code"`
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
	NextReviewDate                  *string    `json:"next_review_date"`
	PublishedAt                     *string    `json:"published_at"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedAt                       *time.Time `json:"updated_at"`
}

type PolicyReportKeySummaryDTO struct {
	ClausesDefined            int64   `json:"clauses_defined"`
	AcknowledgmentRate        float64 `json:"acknowledgment_rate"`
	TotalAcknowledgments      int64   `json:"total_acknowledgments"`
	AcknowledgmentsCompleted  int64   `json:"acknowledgments_completed"`
	AcknowledgmentsPending    int64   `json:"acknowledgments_pending"`
	RiskReported              int64   `json:"risk_reported"`
	AssociatedStandardsCount  int64   `json:"associated_standards_count"`
	AssociatedClausesCount    int64   `json:"associated_clauses_count"`
	StrategicObjectivesLinked int64   `json:"strategic_objectives_linked"`
	ApprovalCycle             int64   `json:"approval_cycle"`
}

type PolicyAcknowledgementListDTO struct {
	UserID                  uint    `json:"user_id"`
	FirstNameEn             string  `json:"first_name_en"`
	FirstNameAr             string  `json:"first_name_ar"`
	LastNameEn              string  `json:"last_name_en"`
	LastNameAr              string  `json:"last_name_ar"`
	OrganizationUnitTitleEn string  `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr string  `json:"organization_unit_title_ar"`
	Acknowledged            bool    `json:"acknowledged"`
	AcknowledgedOn          *string `json:"acknowledged_on"`
	Comment                 *string `json:"comment"`
}

type PolicyAcknowledgementOrganizatationUnitDTO struct {
	OrganizationUnitID      *uint   `json:"organization_unit_id"`
	OrganizationUnitTitleEn string  `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr string  `json:"organization_unit_title_ar"`
	TotalUsers              int64   `json:"total_users"`
	AcknowledgedUsers       int64   `json:"acknowledged_users"`
	NotAcknowledgedUsers    int64   `json:"not_acknowledged_users"`
	AcknowledgedPercent     float64 `json:"acknowledged_percent"`
}

type AssociatedRiskListDTO struct {
	Title                      string  `json:"title"`
	Description                *string `json:"description"`
	RiskCategoryTitleEn        string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr        string  `json:"risk_category_title_ar"`
	RiskStatusTitleEn          string  `json:"risk_status_title_en"`
	RiskStatusTitleAr          string  `json:"risk_status_title_ar"`
	RiskStatusColorCode        string  `json:"risk_status_color_code"`
	CurrentLikelihoodTitleEn   *string `json:"current_likelihood_title_en"`
	CurrentLikelihoodTitleAr   *string `json:"current_likelihood_title_ar"`
	CurrentImpactTitleEn       *string `json:"current_impact_title_en"`
	CurrentImpactTitleAr       *string `json:"current_impact_title_ar"`
	CurrentRiskRatingTitleEn   *string `json:"current_risk_rating_title_en"`
	CurrentRiskRatingTitleAr   *string `json:"current_risk_rating_title_ar"`
	CurrentRiskRatingColorCode *string `json:"current_risk_rating_color_code"`
}

type AssociatedStrategicObjectiveListDTO struct {
	TitleEn                           string  `json:"title_en"`
	TitleAr                           string  `json:"title_ar"`
	DescriptionEn                     *string `json:"description_en"`
	DescriptionAr                     *string `json:"description_ar"`
	StartDate                         *string `json:"start_date"`
	EndDate                           *string `json:"end_date"`
	StrategicObjectiveCategoryTitleEn string  `json:"strategic_objective_category_title_en"`
	StrategicObjectiveCategoryTitleAr string  `json:"strategic_objective_category_title_ar"`
	ResponsibleUserFirstNameEn        *string `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn         *string `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr        *string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr         *string `json:"responsible_user_last_name_ar"`
	ResponsibleUnitTitleEn            *string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr            *string `json:"responsible_unit_title_ar"`
}

type PolicyVersionDTO struct {
	ID                     uint    `json:"id"`
	DocumentID             uint    `json:"document_id"`
	Version                float64 `json:"version"`
	PublishedAt            *string `json:"published_at"`
	PublishedByFirstNameEn *string `json:"published_by_first_name_en"`
	PublishedByFirstNameAr *string `json:"published_by_first_name_ar"`
	PublishedByLastNameEn  *string `json:"published_by_last_name_en"`
	PublishedByLastNameAr  *string `json:"published_by_last_name_ar"`
}

type PolicyAccessUserDTO struct {
	UserID                  uint   `json:"user_id"`
	FirstNameEn             string `json:"first_name_en"`
	LastNameEn              string `json:"last_name_en"`
	FirstNameAr             string `json:"first_name_ar"`
	LastNameAr              string `json:"last_name_ar"`
	DesignationTitleEn      string `json:"designation_title_en"`
	DesignationTitleAr      string `json:"designation_title_ar"`
	OrganizationUnitTitleEn string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr string `json:"organization_unit_title_ar"`
}

type PolicyApprovalHistoryDTO struct {
	ID                    uint    `json:"id"`
	DocumentWorkflowID    uint    `json:"document_workflow_id"`
	WorkflowLevel         int     `json:"workflow_level"`
	ApprovedByID          uint    `json:"approved_by_id"`
	ApprovedByFirstNameEn string  `json:"approved_by_first_name_en"`
	ApprovedByFirstNameAr string  `json:"approved_by_first_name_ar"`
	ApprovedByLastNameEn  string  `json:"approved_by_last_name_en"`
	ApprovedByLastNameAr  string  `json:"approved_by_last_name_ar"`
	DesignationTitleEn    string  `json:"designation_title_en"`
	DesignationTitleAr    string  `json:"designation_title_ar"`
	ActionTitleEn         string  `json:"action_title_en"`
	ActionTitleAr         string  `json:"action_title_ar"`
	ActionType            string  `json:"action_type"`
	ActionColorCode       string  `json:"action_color_code"`
	ActionDate            *string `json:"action_date"`
	Comment               string  `json:"comment"`
}

type PolicyAttachmentDTO struct {
	DocumentId          uint   `json:"document_id"`
	DocumentTitleEn     string `json:"document_title_en"`
	DocumentTitleAr     string `json:"document_title_ar"`
	DocumentTypeTitleEn string `json:"document_type_title_en"`
	DocumentTypeTitleAr string `json:"document_type_title_ar"`
	CreatedAt           string `json:"created_at"`
}

type PolicySummaryReportDTO struct {
	ID                                    uint                                                           `json:"id"`
	Title                                 string                                                         `json:"title"`
	DocumentTitleEn                       string                                                         `json:"document_title_en"`
	DocumentTitleAr                       string                                                         `json:"document_title_ar"`
	CreatedAt                             string                                                         `json:"created_at"`
	CreatedUserFirstNameEn                string                                                         `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr                string                                                         `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                 string                                                         `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                 string                                                         `json:"created_user_last_name_ar"`
	PolicyOverview                        *PolicyReportOverviewResponseDTO                               `json:"policy_overview"`
	PolicyKeyMetrics                      *PolicyReportKeyMetricsResponseDTO                             `json:"policy_key_metrics"`
	ClauseList                            []*PolicyReportClauseResponseDTO                               `json:"clause_list"`
	AcknowledgementUserList               []*PolicyReportAcknowledgementUserResponseDTO                  `json:"acknowldegement_user_list"`
	OrganizatationUnitAcknowledgementList *PolicyReportAcknowledgementOrganizationUnitSummaryResponseDTO `json:"organization_unit_acknowledgement_list"`
	AssociatedRiskList                    []*PolicyReportRiskResponseDTO                                 `json:"associated_risk_list"`
	AssociatedStrategicObjectiveLIst      []*PolicyReportStrategicObjectiveResponseDTO                   `json:"associated_strategic_objective_list"`
	AssociatedComplianceList              []*PolicyReportComplianceResponseDTO                           `json:"associated_compliance_list"`
	PolicyVersionList                     []*PolicyReportDocumentVersionResponseDTO                      `json:"policy_version_list"`
	PolicyAccessUserList                  []*PolicyReportDocumentAccessResponseDTO                       `json:"policy_access_user_list"`
	PolicyApprovalHistroyList             []*PolicyReportDocumentApprovalHistoryResponseDTO              `json:"policy_apporval_histroy_list"`
	PolicyAttachmentList                  []*PolicyReportDocumentAttachmentResponseDTO                   `json:"policy_attachment_list"`
	AppUrl                                string                                                         `json:"url"`
}

type PolicyReportOverviewResponseDTO struct {
	ID                          uint    `json:"id"`
	PolicyReportID              uint    `json:"policy_report_id"`
	ReferenceNumber             string  `json:"reference_number"`
	PolicyTitleEn               string  `json:"policy_title_en"`
	PolicyTitleAr               string  `json:"policy_title_ar"`
	DescriptionEn               string  `json:"description_en"`
	DescriptionAr               string  `json:"description_ar"`
	ResponsibleUnitTitleEn      string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr      string  `json:"responsible_unit_title_ar"`
	ResponsibleUserNameEn       string  `json:"responsible_user_name_en"`
	ResponsibleUserNameAr       string  `json:"responsible_user_name_ar"`
	PublishedDate               *string `json:"published_date"`
	NextReviewDate              *string `json:"next_review_date"`
	DocumentVersion             string  `json:"document_version"`
	DocumentStatusTitleEn       string  `json:"document_status_title_en"`
	DocumentStatusTitleAr       string  `json:"document_status_title_ar"`
	DocumentStatusColorCode     string  `json:"document_status_color_code"`
	DocumentAccessTypeTitleEn   string  `json:"document_access_type_title_en"`
	DocumentAccessTypeTitleAr   string  `json:"document_access_type_title_ar"`
	DocumentAccessTypeColorCode string  `json:"document_access_type_color_code"`
}

type PolicyReportKeyMetricsResponseDTO struct {
	ID                             uint    `json:"id"`
	PolicyReportID                 uint    `json:"policy_report_id"`
	ClausesDefinedCount            int     `json:"clauses_defined_count"`
	TotalAcknowledgementCount      int     `json:"total_acknowledgement_count"`
	AcknowledgementsCompletedCount int     `json:"acknowledgements_completed_count"`
	AcknowledgementsPendingCount   int     `json:"acknowledgements_pending_count"`
	AcknowledgementPercentage      float32 `json:"acknowledgement_percentage"`
	RisksReportedCount             int     `json:"risks_reported_count"`
	AssociatedStandardsCount       int     `json:"associated_standards_count"`
	AssociatedClausesCount         int     `json:"associated_clauses_count"`
	StrategicObjectivesLinkedCount int     `json:"strategic_objectives_linked_count"`
	ApprovalCycleCount             int     `json:"approval_cycle_count"`
}

type PolicyReportClauseResponseDTO struct {
	ID             uint                                      `json:"id"`
	PolicyReportID uint                                      `json:"policy_report_id"`
	ClauseNumber   string                                    `json:"clause_number"`
	ClauseTitleEn  string                                    `json:"clause_title_en"`
	ClauseTitleAr  string                                    `json:"clause_title_ar"`
	Description    string                                    `json:"description"`
	Controls       *string                                   `json:"controls"`
	ControlsParsed []*organizationDTOs.ControlMinimalListDTO `json:"controls_parsed" gorm:"-"`
}

type PolicyReportAcknowledgementUserResponseDTO struct {
	ID                      uint    `json:"id"`
	PolicyReportID          uint    `json:"policy_report_id"`
	UserNameEn              string  `json:"user_name_en"`
	UserNameAr              string  `json:"user_name_ar"`
	OrganizationUnitTitleEn string  `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr string  `json:"organization_unit_title_ar"`
	Acknowledged            bool    `json:"acknowledged"`
	AcknowledgedOn          *string `json:"acknowledged_on"`
	Comment                 *string `json:"comment"`
}

type PolicyReportAcknowledgementOrganizationUnitResponseDTO struct {
	ID                      uint    `json:"id"`
	PolicyReportID          uint    `json:"policy_report_id"`
	OrganizationUnitTitleEn string  `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr string  `json:"organization_unit_title_ar"`
	TotalUsers              int     `json:"total_users"`
	AcknowledgedUsers       int     `json:"acknowledged_users"`
	NotAcknowledgedUsers    int     `json:"not_acknowledged_users"`
	AcknowledgedPercent     float32 `json:"acknowledged_percent"`
}

type PolicyReportAcknowledgementOrganizationUnitSummaryResponseDTO struct {
	Items                     []*PolicyReportAcknowledgementOrganizationUnitResponseDTO `json:"items"`
	TotalUsers                int                                                       `json:"total_users"`
	TotalAcknowledgedUsers    int                                                       `json:"total_acknowledged_users"`
	TotalNotAcknowledgedUsers int                                                       `json:"total_not_acknowledged_users"`
	TotalAcknowledgedPercent  float32                                                   `json:"total_acknowledged_percent"`
}

type PolicyReportRiskResponseDTO struct {
	ID                    uint   `json:"id"`
	PolicyReportID        uint   `json:"policy_report_id"`
	RiskTitle             string `json:"risk_title"`
	Description           string `json:"description"`
	RiskCategoryTitleEn   string `json:"risk_category_title_en"`
	RiskCategoryTitleAr   string `json:"risk_category_title_ar"`
	RiskStatusTitleEn     string `json:"risk_status_title_en"`
	RiskStatusTitleAr     string `json:"risk_status_title_ar"`
	RiskStatusColorCode   string `json:"risk_status_color_code"`
	RiskLikelihoodTitleEn string `json:"risk_likelihood_title_en"`
	RiskLikelihoodTitleAr string `json:"risk_likelihood_title_ar"`
	RiskImpactTitleEn     string `json:"risk_impact_title_en"`
	RiskImpactTitleAr     string `json:"risk_impact_title_ar"`
	RiskRatingTitleEn     string `json:"risk_rating_title_en"`
	RiskRatingTitleAr     string `json:"risk_rating_title_ar"`
	RiskRatingColorCode   string `json:"risk_rating_color_code"`
}

type PolicyReportStrategicObjectiveResponseDTO struct {
	ID                                uint    `json:"id"`
	PolicyReportID                    uint    `json:"policy_report_id"`
	StrategicObjectiveTitleEn         string  `json:"strategic_objective_title_en"`
	StrategicObjectiveTitleAr         string  `json:"strategic_objective_title_ar"`
	DescriptionEn                     *string `json:"description_en"`
	DescriptionAr                     *string `json:"description_ar"`
	StartDate                         *string `json:"start_date"`
	EndDate                           *string `json:"end_date"`
	StrategicObjectiveCategoryTitleEn string  `json:"strategic_objective_category_title_en"`
	StrategicObjectiveCategoryTitleAr string  `json:"strategic_objective_category_title_ar"`
	ResponsibleUserNameEn             string  `json:"responsible_user_name_en"`
	ResponsibleUserNameAr             string  `json:"responsible_user_name_ar"`
	ResponsibleUnitTitleEn            string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr            string  `json:"responsible_unit_title_ar"`
}

type PolicyReportComplianceResponseDTO struct {
	ID                uint                              `json:"id"`
	PolicyReportID    uint                              `json:"policy_report_id"`
	ComplianceTitleEn string                            `json:"compliance_title_en"`
	ComplianceTitleAr string                            `json:"compliance_title_ar"`
	Clauses           *string                           `json:"clauses"`
	ClausesParsed     []*commonDTOs.AssociatedClauseDTO `json:"clauses_parsed" gorm:"-"`
}

type PolicyReportDocumentVersionResponseDTO struct {
	ID                  uint    `json:"id"`
	PolicyReportID      uint    `json:"policy_report_id"`
	Version             int     `json:"version"`
	PublishedAt         *string `json:"published_at"`
	PublishedUserNameEn string  `json:"published_user_name_en"`
	PublishedUserNameAr string  `json:"published_user_name_ar"`
}

type PolicyReportDocumentAccessResponseDTO struct {
	ID                      uint   `json:"id"`
	PolicyReportID          uint   `json:"policy_report_id"`
	UserNameEn              string `json:"user_name_en"`
	UserNameAr              string `json:"user_name_ar"`
	DesignationTitleEn      string `json:"designation_title_en"`
	DesignationTitleAr      string `json:"designation_title_ar"`
	OrganizationUnitTitleEn string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr string `json:"organization_unit_title_ar"`
}

type PolicyReportDocumentApprovalHistoryResponseDTO struct {
	ID                 uint    `json:"id"`
	PolicyReportID     uint    `json:"policy_report_id"`
	ApprovedUserNameEn string  `json:"approved_user_name_en"`
	ApprovedUserNameAr string  `json:"approved_user_name_ar"`
	DesignationTitleEn string  `json:"designation_title_en"`
	DesignationTitleAr string  `json:"designation_title_ar"`
	ActionTitleEn      string  `json:"action_title_en"`
	ActionTitleAr      string  `json:"action_title_ar"`
	ActionColorCode    string  `json:"action_color_code"`
	ActionDate         *string `json:"action_date"`
	Comment            *string `json:"comments"`
}

type PolicyReportDocumentAttachmentResponseDTO struct {
	ID                  uint    `json:"id"`
	PolicyReportID      uint    `json:"policy_report_id"`
	DocumentTitleEn     string  `json:"document_title_en"`
	DocumentTitleAr     string  `json:"document_title_ar"`
	DocumentTypeTitleEn string  `json:"document_type_title_en"`
	DocumentTypeTitleAr string  `json:"document_type_title_ar"`
	LinkedDate          *string `json:"linked_date"`
}
