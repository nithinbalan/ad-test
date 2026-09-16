package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type PenetrationTestReportRequestDTO struct {
	Years                            []int   `json:"years" binding:"omitempty,dive,min=1,max=9999"`
	PenetrationTestExecutorIDs       []uint  `json:"penetration_test_executor_ids" binding:"omitempty,dive,min=1"`
	PenetrationTestTypeIDs           []uint  `json:"penetration_test_type_ids" binding:"omitempty,dive,min=1"`
	PenetrationTestClassificationIDs []uint  `json:"penetration_test_classification_ids" binding:"omitempty,dive,min=1"`
	PenetrationTestPlanIDs           []uint  `json:"penetration_test_plan_ids" binding:"omitempty,dive,min=1"`
	CSResponsibleEmployeeIDs         []uint  `json:"cs_responsible_employee_ids" binding:"omitempty,dive,min=1"`
	ResponsibleUnitIDs               []uint  `json:"responsible_unit_ids" binding:"omitempty,dive,min=1"`
	PenetrationTestStatusIDs         []uint  `json:"penetration_test_status_ids" binding:"omitempty,dive,min=1"`
	Type                             string  `json:"type" binding:"required,oneof=individual date-range year filtered"`
	PenetrationTestIDs               *[]uint `json:"penetration_test_ids"`
	Title                            string  `json:"title" binding:"required"`
	StartDate                        *string `json:"start_date"`
	EndDate                          *string `json:"end_date"`
}

type PenetrationTestReportQueryParamsDTO struct {
	Years                            []int
	PenetrationTestExecutorIDs       []uint
	PenetrationTestTypeIDs           []uint
	PenetrationTestClassificationIDs []uint
	PenetrationTestPlanIDs           []uint
	CSResponsibleEmployeeIDs         []uint
	ResponsibleUnitIDs               []uint
	PenetrationTestStatusIDs         []uint
	PenetrationTestIDs               []uint
	StartDate                        *time.Time
	EndDate                          *time.Time
}

type PenetrationTestReportParsedDates struct {
	StartDate *time.Time
	EndDate   *time.Time
}

type PenetrationTestReportListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type PenetrationTestReportListDTO struct {
	ID                          uint    `json:"id"`
	Title                       string  `json:"title"`
	StartDate                   *string `json:"start_date"`
	EndDate                     *string `json:"end_date"`
	CreatedAt                   string  `json:"created_at"`
	CreatedUserID               uint    `json:"created_user_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath *string `json:"created_user_profile_image_path"`
}

type PenetrationTestReportListResponseDTO struct {
	Items      []*PenetrationTestReportListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO    `json:"pagination"`
}

type PenetrationTestReportUserDTO struct {
	ID               uint    `json:"id" mapstructure:"ID"`
	FirstNameEn      string  `json:"first_name_en" mapstructure:"FirstNameEn"`
	LastNameEn       string  `json:"last_name_en" mapstructure:"LastNameEn"`
	FirstNameAr      string  `json:"first_name_ar" mapstructure:"FirstNameAr"`
	LastNameAr       string  `json:"last_name_ar" mapstructure:"LastNameAr"`
	ProfileImagePath *string `json:"profile_image_path" mapstructure:"ProfileImagePath"`
	Email            string  `json:"email" mapstructure:"Email"`
	Initial          string  `json:"initial" mapstructure:"Initial"`
	Color            string  `json:"color" mapstructure:"Color"`
	MobileNumber     *string `json:"mobile_number" mapstructure:"MobileNumber"`
}

type PenetrationTestReportSectionDTO struct {
	ID            uint    `json:"id" mapstructure:"ID"`
	Type          string  `json:"type" mapstructure:"Type"`
	TitleEn       string  `json:"title_en" mapstructure:"TitleEn"`
	TitleAr       string  `json:"title_ar" mapstructure:"TitleAr"`
	DescriptionEn *string `json:"description_en" mapstructure:"DescriptionEn"`
	DescriptionAr *string `json:"description_ar" mapstructure:"DescriptionAr"`
	OrderNum      int     `json:"order_num" mapstructure:"OrderNum"`
}

type PenetrationTestReportSectionFormDTO struct {
	TitleEn       string  `json:"title_en" binding:"required"`
	TitleAr       string  `json:"title_ar" binding:"required"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
}

type PenetrationTestReportDetailsDTO struct {
	ID               uint                                       `json:"id" mapstructure:"ID"`
	Title            string                                     `json:"title" mapstructure:"Title"`
	StartDate        *time.Time                                 `json:"start_date" mapstructure:"StartDate"`
	EndDate          *time.Time                                 `json:"end_date" mapstructure:"EndDate"`
	CreatedAt        time.Time                                  `json:"created_at" mapstructure:"CreatedAt"`
	CreatedByUser    *PenetrationTestReportUserDTO              `json:"created_by" mapstructure:"CreatedByUser"`
	Sections         []*PenetrationTestReportSectionDTO         `json:"sections" mapstructure:"Sections"`
	PenetrationTests []*PenetrationTestReportPenetrationTestDTO `json:"penetration_tests" mapstructure:"PenetrationTests"`
}

type PenetrationTestReportPenetrationTestDTO struct {
	ID                                   uint                                `json:"id" mapstructure:"ID"`
	PenetrationTestReportID              uint                                `json:"penetration_test_report_id" mapstructure:"PenetrationTestReportID"`
	PenetrationTestID                    uint                                `json:"penetration_test_id" mapstructure:"PenetrationTestID"`
	ReferenceNumber                      *string                             `json:"reference_number" mapstructure:"ReferenceNumber"`
	SystemApplicationName                string                              `json:"system_application_name" mapstructure:"SystemApplicationName"`
	Scope                                *string                             `json:"scope" mapstructure:"Scope"`
	Year                                 int                                 `json:"year" mapstructure:"Year"`
	Period                               string                              `json:"period" mapstructure:"Period"`
	PenetrationTestTypeTitleEn           *string                             `json:"penetration_test_type_title_en" mapstructure:"PenetrationTestTypeTitleEn"`
	PenetrationTestTypeTitleAr           *string                             `json:"penetration_test_type_title_ar" mapstructure:"PenetrationTestTypeTitleAr"`
	PenetrationTestPlanTitleEn           *string                             `json:"penetration_test_plan_title_en" mapstructure:"PenetrationTestPlanTitleEn"`
	PenetrationTestPlanTitleAr           *string                             `json:"penetration_test_plan_title_ar" mapstructure:"PenetrationTestPlanTitleAr"`
	PenetrationTestClassificationTitleEn *string                             `json:"penetration_test_classification_title_en" mapstructure:"PenetrationTestClassificationTitleEn"`
	PenetrationTestClassificationTitleAr *string                             `json:"penetration_test_classification_title_ar" mapstructure:"PenetrationTestClassificationTitleAr"`
	PenetrationTestExecutorTitleEn       *string                             `json:"penetration_test_executor_title_en" mapstructure:"PenetrationTestExecutorTitleEn"`
	PenetrationTestExecutorTitleAr       *string                             `json:"penetration_test_executor_title_ar" mapstructure:"PenetrationTestExecutorTitleAr"`
	CSResponsibleEmployeeNameEn          *string                             `json:"cs_responsible_employee_name_en" mapstructure:"CSResponsibleEmployeeNameEn"`
	CSResponsibleEmployeeNameAr          *string                             `json:"cs_responsible_employee_name_ar" mapstructure:"CSResponsibleEmployeeNameAr"`
	PenetrationTestStatusTitleEn         *string                             `json:"penetration_test_status_title_en" mapstructure:"PenetrationTestStatusTitleEn"`
	PenetrationTestStatusTitleAr         *string                             `json:"penetration_test_status_title_ar" mapstructure:"PenetrationTestStatusTitleAr"`
	PenetrationTestStatusColorCode       *string                             `json:"penetration_test_status_color_code" mapstructure:"PenetrationTestStatusColorCode"`
	EnvironmentReadinessTitleEn          *string                             `json:"environment_readiness_title_en" mapstructure:"EnvironmentReadinessTitleEn"`
	EnvironmentReadinessTitleAr          *string                             `json:"environment_readiness_title_ar" mapstructure:"EnvironmentReadinessTitleAr"`
	EnvironmentNotifyDate                *time.Time                          `json:"environment_notify_date" mapstructure:"EnvironmentNotifyDate"`
	ExpectedReadyDate                    *time.Time                          `json:"expected_ready_date" mapstructure:"ExpectedReadyDate"`
	ExecutionStartDate                   *time.Time                          `json:"execution_start_date" mapstructure:"ExecutionStartDate"`
	ExecutionEndDate                     *time.Time                          `json:"execution_end_date" mapstructure:"ExecutionEndDate"`
	ExecutionStatusTitleEn               *string                             `json:"execution_status_title_en" mapstructure:"ExecutionStatusTitleEn"`
	ExecutionStatusTitleAr               *string                             `json:"execution_status_title_ar" mapstructure:"ExecutionStatusTitleAr"`
	SecurityResultTitleEn                *string                             `json:"security_result_title_en" mapstructure:"SecurityResultTitleEn"`
	SecurityResultTitleAr                *string                             `json:"security_result_title_ar" mapstructure:"SecurityResultTitleAr"`
	ResultsSharedWithManagement          bool                                `json:"results_shared_with_management" mapstructure:"ResultsSharedWithManagement"`
	ManagementNotifyDate                 *time.Time                          `json:"management_notify_date" mapstructure:"ManagementNotifyDate"`
	ManagementResponseDate               *time.Time                          `json:"management_response_date" mapstructure:"ManagementResponseDate"`
	RemediationNotes                     *string                             `json:"remediation_notes" mapstructure:"RemediationNotes"`
	ClosureDate                          *time.Time                          `json:"closure_date" mapstructure:"ClosureDate"`
	ClosureRemarks                       *string                             `json:"closure_remarks" mapstructure:"ClosureRemarks"`
	ClosureOwnerNameEn                   *string                             `json:"closure_owner_name_en" mapstructure:"ClosureOwnerNameEn"`
	ClosureOwnerNameAr                   *string                             `json:"closure_owner_name_ar" mapstructure:"ClosureOwnerNameAr"`
	CreatedAt                            *time.Time                          `json:"created_at" mapstructure:"CreatedAt"`
	UpdatedAt                            *time.Time                          `json:"updated_at" mapstructure:"UpdatedAt"`
	Findings                             []*PenetrationTestReportFindingDTO  `json:"findings" mapstructure:"Findings"`
	Retests                              []*PenetrationTestReportRetestDTO   `json:"retests" mapstructure:"Retests"`
	Documents                            []*PenetrationTestReportDocumentDTO `json:"documents" mapstructure:"Documents"`

	PenetrationTest *PenetrationTestReportSourceDTO  `json:"penetration_test,omitempty" mapstructure:"PenetrationTest"`
	Closure         *PenetrationTestClosureDetailDTO `json:"closure,omitempty" mapstructure:"-"`
}

type PenetrationTestReportFindingDTO struct {
	ID                                     uint       `json:"id" mapstructure:"ID"`
	PenetrationTestReportID                uint       `json:"penetration_test_report_id" mapstructure:"PenetrationTestReportID"`
	PenetrationTestReportPenetrationTestID uint       `json:"penetration_test_report_penetration_test_id" mapstructure:"PenetrationTestReportPenetrationTestID"`
	PenetrationTestFindingID               uint       `json:"penetration_test_finding_id" mapstructure:"PenetrationTestFindingID"`
	ReferenceNumber                        *string    `json:"reference_number" mapstructure:"ReferenceNumber"`
	Title                                  string     `json:"title" mapstructure:"Title"`
	Description                            *string    `json:"description" mapstructure:"Description"`
	Recommendation                         *string    `json:"recommendation" mapstructure:"Recommendation"`
	SeverityLevelTitleEn                   *string    `json:"severity_level_title_en" mapstructure:"SeverityLevelTitleEn"`
	SeverityLevelTitleAr                   *string    `json:"severity_level_title_ar" mapstructure:"SeverityLevelTitleAr"`
	SeverityLevelColorCode                 *string    `json:"severity_level_color_code" mapstructure:"SeverityLevelColorCode"`
	OwnerNameEn                            *string    `json:"owner_name_en" mapstructure:"OwnerNameEn"`
	OwnerNameAr                            *string    `json:"owner_name_ar" mapstructure:"OwnerNameAr"`
	ActionPlanStatusTitleEn                *string    `json:"action_plan_status_title_en" mapstructure:"ActionPlanStatusTitleEn"`
	ActionPlanStatusTitleAr                *string    `json:"action_plan_status_title_ar" mapstructure:"ActionPlanStatusTitleAr"`
	ActionPlanStatusColorCode              *string    `json:"action_plan_status_color_code" mapstructure:"ActionPlanStatusColorCode"`
	Percentage                             *int       `json:"percentage" mapstructure:"Percentage"`
	TargetRemediationDate                  *time.Time `json:"target_remediation_date" mapstructure:"TargetRemediationDate"`
	CreatedAt                              *time.Time `json:"created_at" mapstructure:"CreatedAt"`
	UpdatedAt                              *time.Time `json:"updated_at" mapstructure:"UpdatedAt"`
}

type PenetrationTestReportRetestDTO struct {
	ID                                     uint                                     `json:"id" mapstructure:"ID"`
	PenetrationTestReportID                uint                                     `json:"penetration_test_report_id" mapstructure:"PenetrationTestReportID"`
	PenetrationTestReportPenetrationTestID uint                                     `json:"penetration_test_report_penetration_test_id" mapstructure:"PenetrationTestReportPenetrationTestID"`
	PenetrationTestRetestID                uint                                     `json:"penetration_test_retest_id" mapstructure:"PenetrationTestRetestID"`
	ReferenceNumber                        *string                                  `json:"reference_number" mapstructure:"ReferenceNumber"`
	StartDate                              time.Time                                `json:"start_date" mapstructure:"StartDate"`
	EndDate                                *time.Time                               `json:"end_date" mapstructure:"EndDate"`
	ExecutorTitleEn                        *string                                  `json:"executor_title_en" mapstructure:"ExecutorTitleEn"`
	ExecutorTitleAr                        *string                                  `json:"executor_title_ar" mapstructure:"ExecutorTitleAr"`
	StatusTitleEn                          *string                                  `json:"status_title_en" mapstructure:"StatusTitleEn"`
	StatusTitleAr                          *string                                  `json:"status_title_ar" mapstructure:"StatusTitleAr"`
	StatusColorCode                        *string                                  `json:"status_color_code" mapstructure:"StatusColorCode"`
	Notes                                  *string                                  `json:"notes" mapstructure:"Notes"`
	CreatedAt                              *time.Time                               `json:"created_at" mapstructure:"CreatedAt"`
	UpdatedAt                              *time.Time                               `json:"updated_at" mapstructure:"UpdatedAt"`
	RetestFindings                         []*PenetrationTestReportRetestFindingDTO `json:"retest_findings" mapstructure:"RetestFindings"`
}

type PenetrationTestReportRetestFindingDTO struct {
	ID                             uint                                `json:"id" mapstructure:"ID"`
	PenetrationTestReportRetestID  uint                                `json:"penetration_test_report_retest_id" mapstructure:"PenetrationTestReportRetestID"`
	PenetrationTestReportFindingID uint                                `json:"penetration_test_report_finding_id" mapstructure:"PenetrationTestReportFindingID"`
	PenetrationTestRetestFindingID uint                                `json:"penetration_test_retest_finding_id" mapstructure:"PenetrationTestRetestFindingID"`
	Comment                        *string                             `json:"comment" mapstructure:"Comment"`
	CommentAt                      *time.Time                          `json:"comment_at" mapstructure:"CommentAt"`
	ResultTitleEn                  *string                             `json:"result_title_en" mapstructure:"ResultTitleEn"`
	ResultTitleAr                  *string                             `json:"result_title_ar" mapstructure:"ResultTitleAr"`
	ResultAt                       *time.Time                          `json:"result_at" mapstructure:"ResultAt"`
	Documents                      []*PenetrationTestReportDocumentDTO `json:"documents" mapstructure:"Documents"`
}

type PenetrationTestReportDocumentDTO struct {
	ID                                          uint       `json:"id" mapstructure:"ID"`
	PenetrationTestReportPenetrationTestID      uint       `json:"penetration_test_report_penetration_test_id" mapstructure:"PenetrationTestReportPenetrationTestID"`
	PenetrationTestSupportingDocumentDocumentID *uint      `json:"penetration_test_supporting_document_document_id" mapstructure:"PenetrationTestSupportingDocumentDocumentID"`
	PenetrationTestRetestFindingDocumentID      *uint      `json:"penetration_test_retest_finding_document_id" mapstructure:"PenetrationTestRetestFindingDocumentID"`
	PenetrationTestReportRetestFindingID        *uint      `json:"penetration_test_report_retest_finding_id" mapstructure:"PenetrationTestReportRetestFindingID"`
	URL                                         *string    `json:"url" mapstructure:"URL"`
	Name                                        *string    `json:"name" mapstructure:"Name"`
	Size                                        *int       `json:"size" mapstructure:"Size"`
	Token                                       *string    `json:"token" mapstructure:"Token"`
	OriginalName                                *string    `json:"original_name" mapstructure:"OriginalName"`
	Type                                        *string    `json:"type" mapstructure:"Type"`
	CreatedByID                                 *uint      `json:"created_by_id" mapstructure:"CreatedBYID"`
	UpdatedByID                                 *uint      `json:"updated_by_id" mapstructure:"UpdatedBYID"`
	CreatedAt                                   time.Time  `json:"created_at" mapstructure:"CreatedAt"`
	UpdatedAt                                   *time.Time `json:"updated_at" mapstructure:"UpdatedAt"`
}

type PenetrationTestReportSourceDTO struct {
	Closure *PenetrationTestClosureDetailDTO `json:"closure,omitempty" mapstructure:"Closure"`
}

type PenetrationTestReportRetestFindingSourceDTO struct {
	ID                       uint       `json:"id"`
	PenetrationTestFindingID uint       `json:"penetration_test_finding_id"`
	Comment                  *string    `json:"comment"`
	CommentAt                *time.Time `json:"comment_at"`
	ResultTitleEn            *string    `json:"result_title_en"`
	ResultTitleAr            *string    `json:"result_title_ar"`
	ResultAt                 *time.Time `json:"result_at"`
}

type PenetrationTestReportDocumentSourceDTO struct {
	ID                                  uint       `json:"id"`
	PenetrationTestSupportingDocumentID *uint      `json:"penetration_test_supporting_document_id"`
	PenetrationTestRetestFindingID      *uint      `json:"penetration_test_retest_finding_id"`
	URL                                 *string    `json:"url"`
	Name                                *string    `json:"name"`
	Size                                *int       `json:"size"`
	Token                               *string    `json:"token"`
	OriginalName                        *string    `json:"original_name"`
	Type                                *string    `json:"type"`
	CreatedByID                         *uint      `json:"created_by_id"`
	UpdatedByID                         *uint      `json:"updated_by_id"`
	CreatedAt                           time.Time  `json:"created_at"`
	UpdatedAt                           *time.Time `json:"updated_at"`
}

type PenetrationTestReportPenetrationTestSourceDTO struct {
	ID                                   uint       `json:"id"`
	ReferenceNumber                      *string    `json:"reference_number"`
	SystemApplicationName                string     `json:"system_application_name"`
	Scope                                *string    `json:"scope"`
	Year                                 int        `json:"year"`
	Period                               string     `json:"period"`
	PenetrationTestTypeTitleEn           *string    `json:"penetration_test_type_title_en"`
	PenetrationTestTypeTitleAr           *string    `json:"penetration_test_type_title_ar"`
	PenetrationTestPlanTitleEn           *string    `json:"penetration_test_plan_title_en"`
	PenetrationTestPlanTitleAr           *string    `json:"penetration_test_plan_title_ar"`
	PenetrationTestClassificationTitleEn *string    `json:"penetration_test_classification_title_en"`
	PenetrationTestClassificationTitleAr *string    `json:"penetration_test_classification_title_ar"`
	PenetrationTestExecutorTitleEn       *string    `json:"penetration_test_executor_title_en"`
	PenetrationTestExecutorTitleAr       *string    `json:"penetration_test_executor_title_ar"`
	CSResponsibleEmployeeNameEn          *string    `json:"cs_responsible_employee_name_en"`
	CSResponsibleEmployeeNameAr          *string    `json:"cs_responsible_employee_name_ar"`
	PenetrationTestStatusTitleEn         *string    `json:"penetration_test_status_title_en"`
	PenetrationTestStatusTitleAr         *string    `json:"penetration_test_status_title_ar"`
	PenetrationTestStatusColorCode       *string    `json:"penetration_test_status_color_code"`
	EnvironmentReadinessTitleEn          *string    `json:"environment_readiness_title_en"`
	EnvironmentReadinessTitleAr          *string    `json:"environment_readiness_title_ar"`
	EnvironmentNotifyDate                *time.Time `json:"environment_notify_date"`
	ExpectedReadyDate                    *time.Time `json:"expected_ready_date"`
	ExecutionStartDate                   *time.Time `json:"execution_start_date"`
	ExecutionEndDate                     *time.Time `json:"execution_end_date"`
	ExecutionStatusTitleEn               *string    `json:"execution_status_title_en"`
	ExecutionStatusTitleAr               *string    `json:"execution_status_title_ar"`
	SecurityResultTitleEn                *string    `json:"security_result_title_en"`
	SecurityResultTitleAr                *string    `json:"security_result_title_ar"`
	ResultsSharedWithManagement          bool       `json:"results_shared_with_management"`
	ManagementNotifyDate                 *time.Time `json:"management_notify_date"`
	ManagementResponseDate               *time.Time `json:"management_response_date"`
	RemediationNotes                     *string    `json:"remediation_notes"`
	ClosureDate                          *time.Time `json:"closure_date"`
	ClosureRemarks                       *string    `json:"closure_remarks"`
	ClosureOwnerNameEn                   *string    `json:"closure_owner_name_en"`
	ClosureOwnerNameAr                   *string    `json:"closure_owner_name_ar"`
	CreatedAt                            time.Time  `json:"created_at"`
	UpdatedAt                            *time.Time `json:"updated_at"`
}

type PenetrationTestReportRetestSourceDTO struct {
	ID              uint       `json:"id"`
	ReferenceNumber *string    `json:"reference_number"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	ExecutorTitleEn *string    `json:"executor_title_en"`
	ExecutorTitleAr *string    `json:"executor_title_ar"`
	StatusTitleEn   *string    `json:"status_title_en"`
	StatusTitleAr   *string    `json:"status_title_ar"`
	StatusColorCode *string    `json:"status_color_code"`
	Notes           *string    `json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type PenetrationTestReportFindingSourceDTO struct {
	ID                        uint       `json:"id"`
	ReferenceNumber           *string    `json:"reference_number"`
	Title                     string     `json:"title"`
	Description               *string    `json:"description"`
	Recommendation            *string    `json:"recommendation"`
	SeverityLevelTitleEn      *string    `json:"severity_level_title_en"`
	SeverityLevelTitleAr      *string    `json:"severity_level_title_ar"`
	SeverityLevelColorCode    *string    `json:"severity_level_color_code"`
	OwnerNameEn               *string    `json:"owner_name_en"`
	OwnerNameAr               *string    `json:"owner_name_ar"`
	ActionPlanStatusTitleEn   *string    `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr   *string    `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode *string    `json:"action_plan_status_color_code"`
	Percentage                *int       `json:"percentage"`
	TargetRemediationDate     *time.Time `json:"target_remediation_date"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 *time.Time `json:"updated_at"`
}

type PenetrationTestReportPDFDetailDTO struct {
	ReferenceNumber        *string                                    `json:"reference_number"`
	ID                     uint                                       `json:"id"`
	Title                  string                                     `json:"title"`
	ReportedAt             string                                     `json:"reported_at"`
	GeneratedAt            string                                     `json:"generated_at"`
	StartDate              *time.Time                                 `json:"start_date"`
	EndDate                *time.Time                                 `json:"end_date"`
	CreatedUserFirstNameEn string                                     `json:"created_user_first_name_en"`
	CreatedUserLastNameEn  string                                     `json:"created_user_last_name_en"`
	CreatedUserFirstNameAr string                                     `json:"created_user_first_name_ar"`
	CreatedUserLastNameAr  string                                     `json:"created_user_last_name_ar"`
	Sections               []*PenetrationTestReportSectionDTO         `json:"sections"`
	PenetrationTests       []*PenetrationTestReportPenetrationTestDTO `json:"penetration_tests"`
}
