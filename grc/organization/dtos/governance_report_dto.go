package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type GovernanceReportListDTO struct {
	ID                          uint   `json:"id"`
	Title                       string `json:"title"`
	StartDate                   string `json:"start_date"`
	EndDate                     string `json:"end_date"`
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

type GovernanceReportListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type GovernanceReportListResponseDTO struct {
	Items      []*GovernanceReportListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type GovernanceReportDetailDTO struct {
	ID                     uint                                 `json:"id"`
	Title                  string                               `json:"title"`
	StartDate              string                               `json:"start_date"`
	EndDate                string                               `json:"end_date"`
	CreatedUserFirstNameEn string                               `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string                               `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string                               `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  string                               `json:"created_user_last_name_ar"`
	CreatedAt              string                               `json:"created_at"`
	Sections               []GovernanceReportSectionResponseDTO `json:"sections"`
	Documents              []GovernanceDocumentResponseDTO      `json:"documents"`
	Meetings               []GovernanceReportMeetingResponseDTO `json:"meetings"`
}

type GovernanceReportSectionResponseDTO struct {
	ID            uint   `json:"id"`
	Type          string `json:"type"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
}

type GovernanceDocumentResponseDTO struct {
	ID                     uint   `json:"id"`
	ReferenceNumber        string `json:"reference_number"`
	DocumentTitleEn        string `json:"document_title_en"`
	DocumentTitleAr        string `json:"document_title_ar"`
	ResponsibleUserNameEn  string `json:"responsible_user_name_en"`
	ResponsibleUserNameAr  string `json:"responsible_user_name_ar"`
	ResponsibleUnitTitleEn string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr string `json:"responsible_unit_title_ar"`
	TotalEmployees         int    `json:"total_employees"`
	AcknowledgedCount      int    `json:"acknowledged_count"`
	PendingCount           int    `json:"pending_count"`
	CreatedByID            uint   `json:"created_by_id"`
	CreatedAt              string `json:"created_at"`
}

type GovernanceReportMeetingResponseDTO struct {
	ID                uint      `json:"id"`
	MeetingTitle      string    `json:"meeting_title"`
	MeetingDate       time.Time `json:"meeting_date"`
	ParticipantsCount int       `json:"participants_count"`
	KeyDecisions      string    `json:"key_decisions"`
	ActionsAssigned   string    `json:"actions_assigned"`
	CreatedByID       uint      `json:"created_by_id"`
	CreatedAt         time.Time `json:"created_at"`
}

type GovernanceReportFormDTO struct {
	Title               string `json:"title" binding:"required"`
	StartDate           string `json:"start_date" binding:"required"`
	EndDate             string `json:"end_date" binding:"required"`
	OrganizationUnitIDs []uint `json:"organization_unit_ids"`
}

type GovernanceReportIdentifiedGovernanceDTO struct {
	ID                        uint      `json:"id"`
	ReferenceNumber           string    `json:"reference_number"`
	Title                     string    `json:"title"`
	Description               *string   `json:"description"`
	ReportedAt                time.Time `json:"reported_at"`
	GovernanceCategoryTitleEn string    `json:"governance_category_title_en"`
	GovernanceCategoryTitleAr string    `json:"governance_category_title_ar"`
	GovernanceSourceTitleEn   *string   `json:"governance_source_title_en"`
	GovernanceSourceTitleAr   *string   `json:"governance_source_title_ar"`
	ResponsibleUnitTitleEn    string    `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr    string    `json:"responsible_unit_title_ar"`

	InherentGovernanceRatingTitleEn   *string `json:"inherent_governance_rating_title_en"`
	InherentGovernanceRatingTitleAr   *string `json:"inherent_governance_rating_title_ar"`
	InherentGovernanceRatingColorCode *string `json:"inherent_governance_rating_color_code"`

	InherentLikelihoodTitleEn *string `json:"inherent_likelihood_title_en"`
	InherentLikelihoodTitleAr *string `json:"inherent_likelihood_title_ar"`

	InherentImpactTitleEn *string `json:"inherent_impact_title_en"`
	InherentImpactTitleAr *string `json:"inherent_impact_title_ar"`

	ResidualGovernanceRatingTitleEn   *string `json:"residual_governance_rating_title_en"`
	ResidualGovernanceRatingTitleAr   *string `json:"residual_governance_rating_title_ar"`
	ResidualGovernanceRatingColorCode *string `json:"residual_governance_rating_color_code"`

	ResidualLikelihoodTitleEn *string `json:"residual_likelihood_title_en"`
	ResidualLikelihoodTitleAr *string `json:"residual_likelihood_title_ar"`

	ResidualImpactTitleEn *string `json:"residual_impact_title_en"`
	ResidualImpactTitleAr *string `json:"residual_impact_title_ar"`

	ExistingControls *string `json:"existing_controls"`
	ProposedControls *string `json:"proposed_controls"`

	GovernanceStatusTitleEn   string `json:"governance_status_title_en"`
	GovernanceStatusTitleAr   string `json:"governance_status_title_ar"`
	GovernanceStatusColorCode string `json:"governance_status_color_code"`
}

type GovernanceReportActionPlanDTO struct {
	ID                               uint      `json:"id"`
	GovernanceReportID               uint      `json:"governance_report_id"`
	GovernanceReportIdentificationID uint      `json:"governance_report_identification_id"`
	ActionPlanTitle                  string    `json:"action_plan_title"`
	ActionPlanDescription            *string   `json:"action_plan_description"`
	SeverityLevelTitleEn             string    `json:"severity_level_title_en"`
	SeverityLevelTitleAr             string    `json:"severity_level_title_ar"`
	ResponsibleUserNameEn            string    `json:"responsible_user_name_en"`
	ResponsibleUserNameAr            string    `json:"responsible_user_name_ar"`
	TargetDate                       time.Time `json:"target_date"`
	ActionPlanStatusTitleEn          string    `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr          string    `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode        string    `json:"action_plan_status_color_code"`
	ProposedControls                 *string   `json:"proposed_controls"`
}

type GovernanceReportKRIDTO struct {
	Title           string  `json:"title"`
	Description     *string `json:"description"`
	LowThreshold    int     `json:"low_threshold"`
	MediumThreshold int     `json:"medium_threshold"`
	UnitTitleEn     string  `json:"unit_title_en"`
	UnitTitleAr     string  `json:"unit_title_ar"`
	CurrentValue    int     `json:"current_value"`
}

type GovernanceReportDocumentDTO struct {
	ID                               uint    `json:"id"`
	GovernanceReportID               uint    `json:"governance_report_id"`
	GovernanceReportIdentificationID uint    `json:"governance_report_identification_id"`
	URL                              string  `json:"url"`
	Name                             *string `json:"name"`
	Size                             int     `json:"size"`
	Token                            string  `json:"token"`
	OriginalName                     string  `json:"original_name"`
	Type                             string  `json:"type"`
}

type GovernanceReportSectionFormDTO struct {
	TitleEn       string  `json:"title_en" binding:"required"`
	TitleAr       string  `json:"title_ar" binding:"required"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
}

type GovernanceReportParsedDates struct {
	StartDate *time.Time
	EndDate   *time.Time
}

type GovernanceReportDocumentFormDTO struct {
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"required"`
}

type DocumentAcknowledgmentSummaryDTO struct {
	TotalUsers             int64   `json:"total_users"`
	AcknowledgedUsers      int64   `json:"acknowledged_users"`
	NotAcknowledgedUsers   int64   `json:"not_acknowledged_users"`
	AcknowledgedPercent    float64 `json:"acknowledged_percent"`
	NotAcknowledgedPercent float64 `json:"not_acknowledged_percent"`
}
