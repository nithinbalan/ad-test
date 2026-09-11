package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
)

type PlannedAdvisoryListDTO struct {
	ID                                uint    `json:"id"`
	ReferenceNumber                   string  `json:"reference_number"`
	CaseNumber                        string  `json:"case_number"`
	TitleEn                           string  `json:"title_en"`
	TitleAr                           string  `json:"title_ar"`
	ReceiptDate                       *string `json:"receipt_date"`
	CaseLink                          string  `json:"case_link"`
	ThreatType                        *string `json:"threat_type"`
	ThreatSeverityScore               *int    `json:"threat_severity_score"`
	LikelihoodScore                   *int    `json:"likelihood_score"`
	RiskLevelID                       *uint   `json:"risk_level_id"`
	RiskLevelValueID                  *uint   `json:"risk_level_value_id"`
	RiskLevelTitleEn                  *string `json:"risk_level_title_en"`
	RiskLevelTitleAr                  *string `json:"risk_level_title_ar"`
	RiskLevelColorCode                *string `json:"risk_level_color_code"`
	RiskLevelType                     *string `json:"risk_level_type"`
	ProcessingTime                    *string `json:"processing_time"`
	Description                       *string `json:"description"`
	AffectedAssets                    *string `json:"affected_assets"`
	Tags                              *string `json:"tags"`
	IOCs                              *string `json:"iocs"`
	AdvisorySourceID                  uint    `json:"advisory_source_id"`
	AdvisorySourceTitleEn             string  `json:"advisory_source_title_en"`
	AdvisorySourceTitleAr             string  `json:"advisory_source_title_ar"`
	CaseClassificationID              uint    `json:"case_classification_id"`
	CaseClassificationTitleEn         string  `json:"case_classification_title_en"`
	CaseClassificationTitleAr         string  `json:"case_classification_title_ar"`
	AdvisoryTypeID                    uint    `json:"advisory_type_id"`
	AdvisoryTypeTitleEn               string  `json:"advisory_type_title_en"`
	AdvisoryTypeTitleAr               string  `json:"advisory_type_title_ar"`
	ImportanceID                      uint    `json:"importance_id"`
	ImportanceTitleEn                 string  `json:"importance_title_en"`
	ImportanceTitleAr                 string  `json:"importance_title_ar"`
	SecurityClassificationID          uint    `json:"security_classification_id"`
	SecurityClassificationTitleEn     string  `json:"security_classification_title_en"`
	SecurityClassificationTitleAr     string  `json:"security_classification_title_ar"`
	AdvisoryRecipientID               uint    `json:"advisory_recipient_id"`
	AdvisoryRecipientFirstNameEn      *string `json:"advisory_recipient_first_name_en"`
	AdvisoryRecipientFirstNameAr      *string `json:"advisory_recipient_first_name_ar"`
	AdvisoryRecipientLastNameEn       *string `json:"advisory_recipient_last_name_en"`
	AdvisoryRecipientLastNameAr       *string `json:"advisory_recipient_last_name_ar"`
	AdvisoryRecipientEmail            *string `json:"advisory_recipient_email"`
	AdvisoryRecipientInitial          *string `json:"advisory_recipient_initial"`
	AdvisoryRecipientColor            *string `json:"advisory_recipient_color"`
	AdvisoryRecipientProfileImagePath *string `json:"advisory_recipient_profile_image_path"`
	OrganizationUnitID                uint    `json:"organization_unit_id"`
	OrganizationUnitTitleEn           string  `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr           string  `json:"organization_unit_title_ar"`
	ExecutedByID                      *uint   `json:"executed_by_id"`
	AdvisoryStatusID                  uint    `json:"advisory_status_id"`
	AdvisoryStatusTitleEn             string  `json:"advisory_status_title_en"`
	AdvisoryStatusTitleAr             string  `json:"advisory_status_title_ar"`
	InternalStatusID                  *uint   `json:"internal_status_id"`
	InternalStatusTitleEn             *string `json:"internal_status_title_en"`
	InternalStatusTitleAr             *string `json:"internal_status_title_ar"`
	NotifyDate                        *string `json:"notify_date"`
	DeadlineDate                      *string `json:"deadline_date"`
	ResolveDate                       *string `json:"resolve_date"`
	ReplyDate                         *string `json:"reply_date"`
	SAEDTicketNo                      *string `json:"saed_ticket_no"`
	Attachment                        *string `json:"attachment"`
	ReportLink                        *string `json:"report_link"`
	Notes                             *string `json:"notes"`
	CreatedAt                         string  `json:"created_at"`
	UpdatedAt                         *string `json:"updated_at"`
	CreatedUserID                     uint    `json:"created_user_id"`
	CreatedUserFirstNameEn            string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr            string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn             string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr             string  `json:"created_user_last_name_ar"`
	CreatedUserInitial                string  `json:"created_user_initial"`
	CreatedUserColor                  string  `json:"created_user_color"`
	CreatedUserEmail                  string  `json:"created_user_email"`
	CreatedUserProfileImagePath       string  `json:"created_user_profile_image_path"`
	UpdatedUserID                     *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn            *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr            *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn             *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr             *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                string  `json:"updated_user_initial"`
	UpdatedUserColor                  string  `json:"updated_user_color"`
	UpdatedUserEmail                  string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath       string  `json:"updated_user_profile_image_path"`
}

type PlannedAdvisoryListQueryParamsDTO struct {
	Sort            string
	SortBy          string
	Keyword         string
	Status          string
	Page            int
	PageSize        int
	Filters         map[string][]string
	QuickFilter     []string
	ReceiptDateFrom string
	ReceiptDateTo   string
	DateFrom        string
	DateTo          string
	Export          bool
}

type PlannedAdvisoryListResponseDTO struct {
	Items      []*PlannedAdvisoryListDTO    `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type PlannedAdvisoryDetailDTO struct {
	ID                     uint                                      `json:"id"`
	ReferenceNumber        string                                    `json:"reference_number"`
	CaseNumber             string                                    `json:"case_number"`
	TitleEn                string                                    `json:"title_en"`
	TitleAr                string                                    `json:"title_ar"`
	ReceiptDate            time.Time                                 `json:"receipt_date"`
	ReceiptDateString      *string                                   `json:"receipt_date_string"`
	CaseLink               string                                    `json:"case_link"`
	ThreatType             *string                                   `json:"threat_type"`
	ThreatSeverity         *ThreatIntelligenceSeverityDTO            `json:"threat_severity"`
	Likelihood             *ThreatIntelligenceLikelihoodDTO          `json:"likelihood"`
	RiskLevelID            *uint                                     `json:"risk_level_id"`
	RiskLevelValueID       *uint                                     `json:"risk_level_value_id"`
	RiskLevel              *ThreatIntelligenceRiskLevelDTO           `json:"risk_level"`
	RiskLevelValue         *ThreatIntelligenceRiskLevelValueDTO      `json:"risk_level_value"`
	ProcessingTime         *string                                   `json:"processing_time"`
	Description            *string                                   `json:"description"`
	AffectedAssets         *string                                   `json:"affected_assets"`
	Tags                   *string                                   `json:"tags"`
	IOCs                   *string                                   `json:"iocs"`
	AdvisorySource         *PlannedAdvisorySourceDTO                 `json:"advisory_source"`
	CaseClassification     *PlannedAdvisoryClassificationDTO         `json:"case_classification"`
	AdvisoryType           *PlannedAdvisoryTypeDTO                   `json:"advisory_type"`
	Importance             *PlannedAdvisoryImportanceDTO             `json:"importance"`
	SecurityClassification *PlannedAdvisorySecurityClassificationDTO `json:"security_classification"`
	AdvisoryRecipient      *organizationDTOs.UserMinimalDTO          `json:"advisory_recipient"`
	OrganizationUnit       *settingsDTOs.OrganizationUnitDTO         `json:"organization_unit"`
	ExecutedBy             *organizationDTOs.UserMinimalDTO          `json:"executed_by"`
	AdvisoryStatus         *PlannedAdvisoryStatusDTO                 `json:"advisory_status"`
	InternalStatus         *PlannedAdvisoryInternalStatusDTO         `json:"internal_status"`
	NotifyDate             *time.Time                                `json:"notify_date"`
	NotifyDateString       *string                                   `json:"notify_date_string"`
	DeadlineDate           *time.Time                                `json:"deadline_date"`
	DeadlineDateString     *string                                   `json:"deadline_date_string"`
	ResolveDate            *time.Time                                `json:"resolve_date"`
	ResolveDateString      *string                                   `json:"resolve_date_string"`
	ReplyDate              *time.Time                                `json:"reply_date"`
	ReplyDateString        *string                                   `json:"reply_date_string"`
	SAEDTicketNo           *string                                   `json:"saed_ticket_no"`
	Attachment             *string                                   `json:"attachment"`
	ReportLink             *string                                   `json:"report_link"`
	Notes                  *string                                   `json:"notes"`
	CreatedByUser          organizationDTOs.UserMinimalDTO           `json:"created_by"`
	UpdatedByUser          organizationDTOs.UserMinimalDTO           `json:"updated_by"`
	CreatedAt              time.Time                                 `json:"created_at"`
	UpdatedAt              *time.Time                                `json:"updated_at"`
	Documents              []*PlannedAdvisoryDocumentDTO             `json:"documents"`
}

type PlannedAdvisoryFormDTO struct {
	CaseNumber               string                        `json:"case_number" binding:"required"`
	TitleEn                  string                        `json:"title_en" binding:"required"`
	TitleAr                  string                        `json:"title_ar" binding:"required"`
	ReceiptDate              string                        `json:"receipt_date" binding:"required"`
	CaseLink                 string                        `json:"case_link" binding:"required"`
	AdvisorySourceID         uint                          `json:"advisory_source_id" binding:"required"`
	CaseClassificationID     uint                          `json:"case_classification_id" binding:"required"`
	AdvisoryTypeID           uint                          `json:"advisory_type_id" binding:"required"`
	ImportanceID             uint                          `json:"importance_id" binding:"required"`
	SecurityClassificationID uint                          `json:"security_classification_id" binding:"required"`
	ThreatType               *string                       `json:"threat_type"`
	ThreatSeverityID         *uint                         `json:"threat_severity_id"`
	LikelihoodID             *uint                         `json:"likelihood_id"`
	Description              *string                       `json:"description"`
	AffectedAssets           *string                       `json:"affected_assets"`
	Tags                     *string                       `json:"tags"`
	IOCs                     *string                       `json:"iocs"`
	AdvisoryRecipientID      uint                          `json:"advisory_recipient_id" binding:"required"`
	OrganizationUnitID       uint                          `json:"organization_unit_id" binding:"required"`
	ExecutedByID             *uint                         `json:"executed_by_id"`
	InternalStatusID         *uint                         `json:"internal_status_id"`
	NotifyDate               *string                       `json:"notify_date"`
	DeadlineDate             *string                       `json:"deadline_date"`
	ResolveDate              *string                       `json:"resolve_date"`
	ReplyDate                *string                       `json:"reply_date"`
	SAEDTicketNo             *string                       `json:"saed_ticket_no"`
	Attachment               *string                       `json:"attachment"`
	ReportLink               *string                       `json:"report_link"`
	Notes                    *string                       `json:"notes"`
	Documents                []*commonDTOs.DocumentFormDTO `json:"documents"`
}

type PlannedAdvisoryParsedDates struct {
	ReceiptDate  time.Time
	NotifyDate   *time.Time
	DeadlineDate *time.Time
	ResolveDate  *time.Time
	ReplyDate    *time.Time
}

type PlannedAdvisoryDocumentDTO struct {
	ID                                  uint   `json:"id"`
	ThreatIntelligencePlannedAdvisoryID uint   `json:"threat_intelligence_planned_advisory_id"`
	URL                                 string `json:"url"`
	Name                                string `json:"name"`
	Size                                int    `json:"size"`
	Token                               string `json:"token"`
	OriginalName                        string `json:"original_name"`
	Type                                string `json:"type"`
}

type PlannedAdvisoryStatusUpdateFormDTO struct {
	AdvisoryStatusID uint `json:"advisory_status_id" binding:"required"`
}
