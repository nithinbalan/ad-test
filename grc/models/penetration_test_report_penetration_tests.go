package models

import "time"

type PenetrationTestReportPenetrationTest struct {
	ID                                   uint       `json:"id" gorm:"primary_key;auto_increment"`
	PenetrationTestReportID              uint       `json:"penetration_test_report_id" gorm:"not null"`
	PenetrationTestID                    uint       `json:"penetration_test_id" gorm:"not null"`
	ReferenceNumber                      *string    `json:"reference_number" gorm:"type:nvarchar(50)"`
	SystemApplicationName                string     `json:"system_application_name" gorm:"type:nvarchar(255);not null"`
	Scope                                *string    `json:"scope" gorm:"type:nvarchar(max)"`
	Year                                 int        `json:"year" gorm:"not null"`
	Period                               string     `json:"period" gorm:"type:varchar(2);not null"`
	PenetrationTestTypeTitleEn           *string    `json:"penetration_test_type_title_en" gorm:"type:nvarchar(255)"`
	PenetrationTestTypeTitleAr           *string    `json:"penetration_test_type_title_ar" gorm:"type:nvarchar(255)"`
	PenetrationTestPlanTitleEn           *string    `json:"penetration_test_plan_title_en" gorm:"type:nvarchar(255)"`
	PenetrationTestPlanTitleAr           *string    `json:"penetration_test_plan_title_ar" gorm:"type:nvarchar(255)"`
	PenetrationTestClassificationTitleEn *string    `json:"penetration_test_classification_title_en" gorm:"type:nvarchar(255)"`
	PenetrationTestClassificationTitleAr *string    `json:"penetration_test_classification_title_ar" gorm:"type:nvarchar(255)"`
	PenetrationTestExecutorTitleEn       *string    `json:"penetration_test_executor_title_en" gorm:"type:nvarchar(255)"`
	PenetrationTestExecutorTitleAr       *string    `json:"penetration_test_executor_title_ar" gorm:"type:nvarchar(255)"`
	CSResponsibleEmployeeNameEn          *string    `json:"cs_responsible_employee_name_en" gorm:"type:nvarchar(255)"`
	CSResponsibleEmployeeNameAr          *string    `json:"cs_responsible_employee_name_ar" gorm:"type:nvarchar(255)"`
	PenetrationTestStatusTitleEn         *string    `json:"penetration_test_status_title_en" gorm:"type:nvarchar(255)"`
	PenetrationTestStatusTitleAr         *string    `json:"penetration_test_status_title_ar" gorm:"type:nvarchar(255)"`
	PenetrationTestStatusColorCode       *string    `json:"penetration_test_status_color_code" gorm:"type:nvarchar(255)"`
	EnvironmentReadinessTitleEn          *string    `json:"environment_readiness_title_en" gorm:"type:nvarchar(255)"`
	EnvironmentReadinessTitleAr          *string    `json:"environment_readiness_title_ar" gorm:"type:nvarchar(255)"`
	EnvironmentNotifyDate                *time.Time `json:"environment_notify_date" gorm:"type:datetime"`
	ExpectedReadyDate                    *time.Time `json:"expected_ready_date" gorm:"type:datetime"`
	ExecutionStartDate                   *time.Time `json:"execution_start_date" gorm:"type:date"`
	ExecutionEndDate                     *time.Time `json:"execution_end_date" gorm:"type:date"`
	ExecutionStatusTitleEn               *string    `json:"execution_status_title_en" gorm:"type:nvarchar(255)"`
	ExecutionStatusTitleAr               *string    `json:"execution_status_title_ar" gorm:"type:nvarchar(255)"`
	SecurityResultTitleEn                *string    `json:"security_result_title_en" gorm:"type:nvarchar(255)"`
	SecurityResultTitleAr                *string    `json:"security_result_title_ar" gorm:"type:nvarchar(255)"`
	ResultsSharedWithManagement          bool       `json:"results_shared_with_management" gorm:"not null;default:false"`
	ManagementNotifyDate                 *time.Time `json:"management_notify_date" gorm:"type:datetime"`
	ManagementResponseDate               *time.Time `json:"management_response_date" gorm:"type:datetime"`
	RemediationNotes                     *string    `json:"remediation_notes" gorm:"type:nvarchar(max)"`
	ClosureDate                          *time.Time `json:"closure_date" gorm:"type:date"`
	ClosureRemarks                       *string    `json:"closure_remarks" gorm:"type:nvarchar(max)"`
	ClosureOwnerNameEn                   *string    `json:"closure_owner_name_en" gorm:"type:nvarchar(255)"`
	ClosureOwnerNameAr                   *string    `json:"closure_owner_name_ar" gorm:"type:nvarchar(255)"`
	CreatedAt                            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                            *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	PenetrationTestReport *PenetrationTestReport           `gorm:"foreignKey:PenetrationTestReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	PenetrationTest       *PenetrationTest                 `gorm:"foreignKey:PenetrationTestID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Findings              []*PenetrationTestReportFinding  `gorm:"foreignKey:PenetrationTestReportPenetrationTestID"`
	Retests               []*PenetrationTestReportRetest   `gorm:"foreignKey:PenetrationTestReportPenetrationTestID"`
	Documents             []*PenetrationTestReportDocument `gorm:"foreignKey:PenetrationTestReportPenetrationTestID"`
}

func (PenetrationTestReportPenetrationTest) TableName() string {
	return "penetration_test_report_penetration_tests"
}
