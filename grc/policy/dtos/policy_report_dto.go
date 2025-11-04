package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"
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

type PolicySummaryReportDTO struct {
	Title                                 string                                             `json:"title"`
	DocumentTitleEn                       string                                             `json:"document_title_en"`
	DocumentTitleAr                       string                                             `json:"document_title_ar"`
	GeneratedOn                           string                                             `json:"generated_on"`
	GeneratedByUserFirstNameEn            string                                             `json:"generated_by_user_first_name_en"`
	GeneratedByUserFirstNameAr            string                                             `json:"generated_by_user_first_name_ar"`
	GeneratedByUserLastNameEn             string                                             `json:"generated_by_user_last_name_en"`
	GeneratedByUserLastNameAr             string                                             `json:"generated_by_user_last_name_ar"`
	GeneratedByUserInitial                string                                             `json:"generated_by_user_initial"`
	GeneratedByUserColor                  string                                             `json:"generated_by_user_color"`
	GeneratedByUserEmail                  string                                             `json:"generated_by_user_email"`
	GeneratedByUserProfileImagePath       string                                             `json:"generated_by_user_profile_image_path"`
	PolicyOverview                        *PolicyReportPolicyOverviewDTO                     `json:"policy_overview"`
	PolicyKeySummary                      *PolicyReportKeySummaryDTO                         `json:"policy_key_summary"`
	ClauseList                            []*documentDTOs.DocumentClauseListDTO              `json:"clause_list"`
	AcknowledgementUserList               []*PolicyAcknowledgementListDTO                    `json:"acknowldegement_user_list"`
	OrganizatationUnitAcknowledgementList *PolicyAcknowledgementOrganizatationUnitSummaryDTO `json:"orgaonizatoin_unit_acknowledgement_list"`
	AssociatedRiskList                    []*commonDTOs.AssociatedRiskListDTO                `json:"associated_risk_list"`
	AssociatedStrategicObjectiveLIst      []*commonDTOs.AssociatedStrategicObjectiveListDTO  `json:"associated_strategic_objective_list"`
	AssociatedComplianceList              []*commonDTOs.AssociatedClauseListDTO              `json:"associated_compliance_list"`
	PolicyVersionList                     []*PolicyVersionDTO                                `json:"policy_version_list"`
	PolicyAccessUserList                  []*PolicyAccessUserDTO                             `json:"policy_access_user_list"`
	PolicyApprovalHistroyList             []*PolicyApprovalHistoryDTO                        `json:"policy_apporval_histroy_list"`
	PolicyAttachmentList                  []*PolicyAttachmentDTO                             `json:"policy_attachment_list"`
	AppUrl                                string                                             `json:"url"`
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
	PublishedAt                     *string    `json:"published_at"`
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
	UserID            uint    `json:"user_id"`
	FirstNameEn       string  `json:"first_name_en"`
	FirstNameAr       string  `json:"first_name_ar"`
	LastNameEn        string  `json:"last_name_en"`
	LastNameAr        string  `json:"last_name_ar"`
	Email             string  `json:"email"`
	Initial           string  `json:"initial"`
	Color             string  `json:"color"`
	ProfileImagePath  string  `json:"profile_image_path"`
	DepartmentTitleEn string  `json:"department_title_en"`
	DepartmentTitleAr string  `json:"department_title_ar"`
	Acknowledged      bool    `json:"acknowledged"`
	AcknowledgedOn    *string `json:"acknowledged_on"`
	Comment           *string `json:"comment"`
}

type PolicyAcknowledgementOrganizatationUnitSummaryItemDTO struct {
	OrganizationUnitID      *uint   `json:"organization_unit_id"`
	OrganizationUnitTitleEn string  `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr string  `json:"organization_unit_title_ar"`
	TotalUsers              int64   `json:"total_users"`
	AcknowledgedUsers       int64   `json:"acknowledged_users"`
	NotAcknowledgedUsers    int64   `json:"not_acknowledged_users"`
	AcknowledgedPercent     float64 `json:"acknowledged_percent"`
}

type PolicyAcknowledgementOrganizatationUnitSummaryDTO struct {
	Items                     []PolicyAcknowledgementOrganizatationUnitSummaryItemDTO `json:"items"`
	TotalUsers                int64                                                   `json:"total_users"`
	TotalAcknowledgedUsers    int64                                                   `json:"total_acknowledged_users"`
	TotalNotAcknowledgedUsers int64                                                   `json:"total_not_acknowledged_users"`
	TotalAcknowledgedPercent  float64                                                 `json:"total_acknowledged_percent"`
}

type PolicyVersionDTO struct {
	ID                     uint    `json:"id"`
	DocumentID             uint    `json:"document_id"`
	Version                float64 `json:"version"`
	IsLatestVersion        bool    `json:"is_latest_version"`
	PublishedAt            string  `json:"published_at"`
	EnableAcknowledgement  bool    `json:"enable_acknowledgement"`
	CreatedAt              string  `json:"created_at"`
	UpdatedAt              *string `json:"updated_at"`
	PublishedByFirstNameEn *string `json:"published_by_first_name_en"`
	PublishedByLastNameEn  *string `json:"published_by_last_name_en"`
	PublishedByEmail       *string `json:"published_by_email"`
}

type PolicyAccessUserDTO struct {
	UserID                  uint    `json:"user_id"`
	FirstNameEn             string  `json:"first_name_en"`
	LastNameEn              string  `json:"last_name_en"`
	FirstNameAr             string  `json:"first_name_ar"`
	LastNameAr              string  `json:"last_name_ar"`
	Email                   string  `json:"email"`
	Phone                   *string `json:"phone"`
	ProfileImagePath        *string `json:"profile_image_path"`
	DesignationTitleEn      string  `json:"designation_title_en"`
	DesignationTitleAr      string  `json:"designation_title_ar"`
	OrganizationUnitTitleEn string  `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr string  `json:"organization_unit_title_ar"`
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
	ActionDate            *string `json:"date"`
	Comment               string  `json:"comment"`
}

type PolicyAttachmentDTO struct {
	DocumentId          uint    `json:"document_id"`
	DocumentTitleEn     string  `json:"document_title_en"`
	DocumentTitleAr     string  `json:"document_title_ar"`
	DocumentTypeTitleEn string  `json:"document_type_title_en"`
	DocumentTypeTitleAr string  `json:"document_type_title_ar"`
	CreatedAt           *string `json:"created_at"`
}
