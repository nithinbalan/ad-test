package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type PenetrationTestListDTO struct {
	ID                                    uint    `json:"id"`
	ReferenceNumber                       string  `json:"reference_number"`
	SystemApplicationName                 string  `json:"system_application_name"`
	PenetrationTestTypeID                 uint    `json:"penetration_test_type_id"`
	PenetrationTestTypeTitleEn            string  `json:"penetration_test_type_title_en"`
	PenetrationTestTypeTitleAr            string  `json:"penetration_test_type_title_ar"`
	PenetrationTestPlanID                 uint    `json:"penetration_test_plan_id"`
	PenetrationTestPlanTitleEn            string  `json:"penetration_test_plan_title_en"`
	PenetrationTestPlanTitleAr            string  `json:"penetration_test_plan_title_ar"`
	Year                                  int     `json:"year"`
	Period                                string  `json:"period"`
	PenetrationTestStatusID               uint    `json:"penetration_test_status_id"`
	StatusTitleEn                         string  `json:"status_title_en"`
	StatusTitleAr                         string  `json:"status_title_ar"`
	StatusColorCode                       string  `json:"status_color_code"`
	StatusType                            string  `json:"status_type"`
	PenetrationTestClassificationID       uint    `json:"penetration_test_classification_id"`
	PenetrationTestClassificationTitleEn  string  `json:"penetration_test_classification_title_en"`
	PenetrationTestClassificationTitleAr  string  `json:"penetration_test_classification_title_ar"`
	PenetrationTestExecutorID             uint    `json:"penetration_test_executor_id"`
	PenetrationTestExecutorTitleEn        string  `json:"penetration_test_executor_title_en"`
	PenetrationTestExecutorTitleAr        string  `json:"penetration_test_executor_title_ar"`
	CSResponsibleEmployeeID               uint    `json:"cs_responsible_employee_id"`
	CSResponsibleEmployeeFirstNameEn      string  `json:"cs_responsible_employee_first_name_en"`
	CSResponsibleEmployeeFirstNameAr      string  `json:"cs_responsible_employee_first_name_ar"`
	CSResponsibleEmployeeLastNameEn       string  `json:"cs_responsible_employee_last_name_en"`
	CSResponsibleEmployeeLastNameAr       string  `json:"cs_responsible_employee_last_name_ar"`
	CSResponsibleEmployeeInitial          string  `json:"cs_responsible_employee_initial"`
	CSResponsibleEmployeeColor            string  `json:"cs_responsible_employee_color"`
	CSResponsibleEmployeeEmail            string  `json:"cs_responsible_employee_email"`
	CSResponsibleEmployeeProfileImagePath string  `json:"cs_responsible_employee_profile_image_path"`
	EnvironmentReadinessID                *uint   `json:"environment_readiness_id"`
	TestStartDate                         string  `json:"test_start_date"`
	TestEndDate                           string  `json:"test_end_date"`
	ExecutionStatusID                     *uint   `json:"execution_status_id"`
	ExecutionStatusTitleEn                *string `json:"execution_status_title_en"`
	ExecutionStatusTitleAr                *string `json:"execution_status_title_ar"`
	SecurityResultID                      *uint   `json:"security_result_id"`
	SecurityResultTitleEn                 *string `json:"security_result_title_en"`
	SecurityResultTitleAr                 *string `json:"security_result_title_ar"`
	ResultsSharedWithManagement           bool    `json:"results_shared_with_management"`
	CriticalCount                         int     `json:"critical_count"`
	HighCount                             int     `json:"high_count"`
	MediumCount                           int     `json:"medium_count"`
	LowCount                              int     `json:"low_count"`
	TotalCount                            int     `json:"total_count"`
}

type PenetrationTestFormDTO struct {
	SystemApplicationName           string  `json:"system_application_name" binding:"required,max=255"`
	PenetrationTestTypeID           uint    `json:"penetration_test_type_id" binding:"required"`
	PenetrationTestPlanID           uint    `json:"penetration_test_plan_id" binding:"required"`
	Year                            int     `json:"year" binding:"required"`
	Period                          string  `json:"period" binding:"required,oneof=Q1 Q2 Q3 Q4"`
	PenetrationTestClassificationID uint    `json:"penetration_test_classification_id" binding:"required"`
	PenetrationTestExecutorID       *uint   `json:"penetration_test_executor_id"`
	CSResponsibleEmployeeID         uint    `json:"cs_responsible_employee_id" binding:"required"`
	Scope                           *string `json:"scope"`
}

type PenetrationTestStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
}

type PenetrationTestDetailDTO struct {
	ID                            uint                                    `json:"id"`
	ReferenceNumber               string                                  `json:"reference_number"`
	SystemApplicationName         string                                  `json:"system_application_name"`
	PenetrationTestType           PenetrationTestTypeDTO                  `json:"penetration_test_type"`
	PenetrationTestPlan           PenetrationTestPlanDTO                  `json:"penetration_test_plan"`
	PenetrationTestClassification PenetrationTestClassificationDTO        `json:"penetration_test_classification"`
	PenetrationTestExecutor       PenetrationTestExecutorDTO              `json:"penetration_test_executor"`
	CSResponsibleEmployee         organizationDTOs.UserMinimalDTO         `json:"cs_responsible_employee"`
	Year                          int                                     `json:"year"`
	Period                        string                                  `json:"period"`
	Scope                         string                                  `json:"scope"`
	PenetrationTestStatus         *PenetrationTestStatusDTO               `json:"penetration_test_status"`
	EnvironmentNotifyDate         *time.Time                              `json:"environment_notify_date"`
	ExpectedReadyDate             *time.Time                              `json:"expected_ready_date"`
	EnvironmentReadiness          *PenetrationTestEnvironmentReadinessDTO `json:"environment_readiness"`
	TestStartDate                 *time.Time                              `json:"test_start_date"`
	TestEndDate                   *time.Time                              `json:"test_end_date"`
	ExecutionStatus               *PenetrationTestExecutionStatusDTO      `json:"execution_status"`
	SecurityResult                *PenetrationTestSecurityResultDTO       `json:"security_result"`
	ResultsSharedWithManagement   bool                                    `json:"results_shared_with_management"`
	ManagementNotifyDate          *time.Time                              `json:"management_notify_date"`
	ManagementResponseDate        *time.Time                              `json:"management_response_date"`
	RemediationNotes              *string                                 `json:"remediation_notes"`
	CriticalCount                 int                                     `json:"critical_count"`
	HighCount                     int                                     `json:"high_count"`
	MediumCount                   int                                     `json:"medium_count"`
	LowCount                      int                                     `json:"low_count"`
	TotalCount                    int                                     `json:"total_count"`
	CreatedAt                     time.Time                               `json:"created_at"`
	UpdatedAt                     *time.Time                              `json:"updated_at"`
	CreatedByUser                 organizationDTOs.UserMinimalDTO         `json:"created_by_user"`
	UpdatedByUser                 *organizationDTOs.UserMinimalDTO        `json:"updated_by_user"`
}

type PenetrationTestSeverityCountDTO struct {
	ID         uint    `json:"id"`
	TitleEn    string  `json:"title_en"`
	TitleAr    string  `json:"title_ar"`
	ColorCode  string  `json:"color_code"`
	Type       string  `json:"type"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type PenetrationTestSeverityBreakdownDTO struct {
	TotalFindings  int64                             `json:"total_findings"`
	SeverityLevels []PenetrationTestSeverityCountDTO `json:"severity_levels"`
}

type PenetrationTestListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type PenetrationTestListResponseDTO struct {
	Items      []*PenetrationTestListDTO    `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
