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
	Title                           string     `json:"title" binding:"required"`
	Description                     *string    `json:"description,omitempty"`
	TargetDate                      *time.Time `json:"target_date,omitempty"`
	URLParam                        *string    `json:"url_param,omitempty"`
	QueryParam                      *string    `json:"query_param,omitempty"`
	SeverityLevelID                 *uint      `json:"severity_level_id,omitempty"`
	ResponsibleUserID               uint       `json:"responsible_user_id" binding:"required"`
	TaskStatusID                    *uint      `json:"task_status_id"`
}

type AssociatedControlListDTO struct {
	ID                                   uint       `json:"id"`
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
}
