package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type BusinessContinuityPlanListDTO struct {
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

type BusinessContinuityPlanFormDTO struct {
	BusinessFunctionID uint                          `json:"business_function_id" binding:"required"`
	Version            *int                          `json:"version"`
	Purpose            *string                       `json:"purpose"`
	Scope              *string                       `json:"scope"`
	Objectives         *string                       `json:"objectives"`
	Assumptions        *string                       `json:"assumptions"`
	ResponsibleUserID  *uint                         `json:"responsible_user_id" binding:"required"`
	Documents          []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type BusinessContinuityPlanListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type BusinessContinuityPlanListResponseDTO struct {
	Items      []*BusinessContinuityPlanListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}

type BusinessContinuityPlanDetailDTO struct {
	ID                 uint                             `json:"id"`
	ReferenceNumber    string                           `json:"reference_number"`
	BusinessFunctionID uint                             `json:"business_function_id"`
	BusinessFunction   BusinessFunctionDetailDTO        `json:"business_function"`
	Version            int                              `json:"version"`
	IsLatestVersion    *bool                            `json:"is_latest_version"`
	Purpose            *string                          `json:"purpose"`
	Scope              *string                          `json:"scope"`
	Objectives         *string                          `json:"objectives"`
	Assumptions        *string                          `json:"assumptions"`
	ResponsibleUserID  *uint                            `json:"responsible_user_id"`
	ResponsibleUser    *organizationDTOs.UserMinimalDTO `json:"responsible_user"`
	ResponsibleUnit    settingsDTOs.OrganizationUnitDTO `json:"responsible_unit"`

	BcpStatus                  BcpStatusDTO                               `json:"bcp_status"`
	Rto                        *int                                       `json:"rto"`
	RtoUnitId                  *uint                                      `json:"rto_unit_id"`
	RtoUnit                    *BusinessContinuityDurationUnitDTO         `json:"rto_unit"`
	RtoCondition               *string                                    `json:"rto_condition"`
	Mtpd                       *int                                       `json:"mtpd"`
	MtpdUnitId                 *uint                                      `json:"mtpd_unit_id"`
	MtpdUnit                   *BusinessContinuityDurationUnitDTO         `json:"mtpd_unit"`
	MtpdCondition              *string                                    `json:"mtpd_condition"`
	Rpo                        *int                                       `json:"rpo"`
	RpoUnitId                  *uint                                      `json:"rpo_unit_id"`
	RpoUnit                    *BusinessContinuityDurationUnitDTO         `json:"rpo_unit"`
	RpoCondition               *string                                    `json:"rpo_condition"`
	Processes                  []*BcpProcessDTO                           `json:"processes"`
	Documents                  []*BcpDocumentDTO                          `json:"documents" gorm:"-"`
	BcpRecoveryPlans           []*BcpRecoveryPlanDetailDTO                `json:"bcp_recovery_plans"`
	BcpCommunicationPlans      []*BcpCommunicationPlanDetailDTO           `json:"bcp_communication_plans"`
	EquipmentTools             []organizationDTOs.ProcessEquipmentToolDTO `json:"equipment_tools"`
	OutageBuckets              []BCMBIAOutageBucketDTO                    `json:"outage_buckets"`
	IsTestAndExerciseCompleted bool                                       `json:"is_test_and_exercise_completed"`
	CreatedAt                  time.Time                                  `json:"created_at"`
	UpdatedAt                  *time.Time                                 `json:"updated_at"`
	CreatedByUser              organizationDTOs.UserMinimalDTO            `json:"created_by"`
	UpdatedByUser              *organizationDTOs.UserMinimalDTO           `json:"updated_by"`
}

type BcpProcessDTO struct {
	ID                       *uint                              `json:"id"`
	BusinessContinuityPlanID uint                               `json:"business_continuity_plan_id"`
	ProcessID                uint                               `json:"process_id"`
	ImpactAnalysisLevelID    *uint                              `json:"impact_analysis_level_id"`
	ImpactAnalysisLevel      *BCMImpactAnalysisLevelDTO         `json:"impact_analysis_level"`
	RtoFromValue             *int                               `json:"rto_from_value"`
	RtoToValue               *int                               `json:"rto_to_value"`
	RtoUnitId                *uint                              `json:"rto_unit_id"`
	RtoUnit                  *BusinessContinuityDurationUnitDTO `json:"rto_unit"`
	RtoCondition             *string                            `json:"rto_condition"`
	MtpdFromValue            *int                               `json:"mtpd_from_value"`
	MtpdToValue              *int                               `json:"mtpd_to_value"`
	MtpdUnitId               *uint                              `json:"mtpd_unit_id"`
	MtpdUnit                 *BusinessContinuityDurationUnitDTO `json:"mtpd_unit"`
	MtpdCondition            *string                            `json:"mtpd_condition"`
	RpoFromValue             *int                               `json:"rpo_from_value"`
	RpoToValue               *int                               `json:"rpo_to_value"`
	RpoUnitId                *uint                              `json:"rpo_unit_id"`
	RpoUnit                  *BusinessContinuityDurationUnitDTO `json:"rpo_unit"`
	RpoCondition             *string                            `json:"rpo_condition"`
	RspFromValue             *int                               `json:"rsp_from_value"`
	RspToValue               *int                               `json:"rsp_to_value"`
	RspUnitID                *uint                              `json:"rsp_unit_id"`
	RspUnit                  *BusinessContinuityDurationUnitDTO `json:"rsp_unit"`
	RspCondition             *string                            `json:"rsp_condition"`
	Process                  organizationDTOs.ProcessDetailDTO  `json:"process"`
	ExecutionRecords         []*ExecutionRecordDTO              `json:"execution_records"`
	ManualProcedureForms     []*ManualProcedureDTO              `json:"manual_procedure_forms"`
	SystemRegisters          []*SystemRegisterDTO               `json:"system_register"`
	InternalDependencies     []*InternalDependencyDTO           `json:"internal_dependencies"`
	ExternalDependencies     []*ExternalDependencyDTO           `json:"external_dependencies"`
	EmployeeAnalyses         []*EmployeeAnalysisDTO             `json:"employee_analyses"`
	StaffAllocations         []*StaffAllocationDTO              `json:"staff_allocations"`
}

type BcpDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}

type BcpTestSummaryDTO struct {
	TestScore        int `json:"test_score"`
	ExecutedTests    int `json:"executed_tests"`
	SuccessfulTests  int `json:"successful_tests"`
	FailedTests      int `json:"failed_tests"`
	PartiallySuccess int `json:"partially_success"`
	NotTestedTests   int `json:"not_tested_tests"`
}

// Detail Dto
type ExecutionRecordDTO struct {
	ID                 *uint                                         `json:"id"`
	BCMExecutionRecord organizationDTOs.BCMExecutionRecordMinimalDTO `json:"bcm_execution_record"`
}

type ManualProcedureDTO struct {
	ID                     *uint                                             `json:"id"`
	BCMManualProcedureForm organizationDTOs.BCMManualProcedureFormMinimalDTO `json:"bcm_manual_procedure_form"`
}

type SystemRegisterDTO struct {
	ID             *uint                                     `json:"id"`
	SystemRegister organizationDTOs.SystemRegisterMinimalDTO `json:"system_register"`
}

type InternalDependencyDTO struct {
	ID                       uint                                       `json:"id"`
	SupportingUnitID         uint                                       `json:"supporting_unit_id"`
	SupportingUnit           *organizationDTOs.OrganizationUnit1DTO     `json:"supporting_unit"`
	DependencyTierID         uint                                       `json:"dependency_tier_id"`
	DependencyTier           *organizationDTOs.ProcessDependencyTierDTO `json:"dependency_tier"`
	DependencyPercentageFrom *int                                       `json:"dependency_percentage_from"`
	DependencyPercentageTo   *int                                       `json:"dependency_percentage_to"`
	RequiredSupport          *string                                    `json:"required_support"`
	Notes                    *string                                    `json:"notes"`
}

type ExternalDependencyDTO struct {
	ID                       uint                                                     `json:"id"`
	EntityName               string                                                   `json:"entity_name"`
	EntityTypeID             *uint                                                    `json:"entity_type_id"`
	EntityType               *organizationDTOs.ProcessExternalDependencyEntityTypeDTO `json:"entity_type"`
	NeededSupportDescription *string                                                  `json:"needed_support_description"`
	DependencyTierID         *uint                                                    `json:"dependency_tier_id"`
	DependencyTier           *organizationDTOs.ProcessDependencyTierDTO               `json:"dependency_tier"`
	SLARegister              *organizationDTOs.SLADetailDTO                           `json:"sla_register"`
	VendorRTO                *int                                                     `json:"vendor_rto"`
	VendorRTOUnit            *BusinessContinuityDurationUnitDTO                       `json:"vendor_rto_unit"`
	InternalStakeholder      *string                                                  `json:"internal_stakeholder"`
	ContactNumber            *string                                                  `json:"contact_number"`
	ContactEmail             *string                                                  `json:"contact_email"`
}

type EmployeeAnalysisDTO struct {
	ID            uint                             `json:"id"`
	DesignationID uint                             `json:"designation_id"`
	Designation   *organizationDTOs.DesignationDTO `json:"designation"`
	NormalCount   *int                             `json:"normal_count"`
	CriticalCount *int                             `json:"critical_count"`
	Employees     []EmployeeAnalysisEmployeeDTO    `json:"employees"`
}

type EmployeeAnalysisEmployeeDTO struct {
	ID                        uint                              `json:"id"`
	ProcessEmployeeAnalysisID uint                              `json:"process_employee_analysis_id"`
	UserID                    uint                              `json:"user_id"`
	User                      organizationDTOs.UserMinimalDTO   `json:"user"`
	Tools                     []EmployeeAnalysisEmployeeToolDTO `json:"equipment_tools"`
}

type EmployeeAnalysisEmployeeToolDTO struct {
	ID            uint                                     `json:"id"`
	EquipmentTool organizationDTOs.ProcessEquipmentToolDTO `json:"equipment_tool"`
	IsHas         bool                                     `json:"is_has"`
	IsNeeded      bool                                     `json:"is_needed"`
}

type StaffAllocationDTO struct {
	ID                         uint                       `json:"id"`
	BcpProcessEmployeeAnalysis *EmployeeAnalysisDTO       `json:"bcp_process_employee_analysis"`
	OutageBuckets              []StaffAllocationBucketDTO `json:"outage_buckets"`
}

type StaffAllocationBucketDTO struct {
	ID             uint                     `json:"id"`
	OutageBucketID uint                     `json:"outage_bucket_id"`
	OutageBucket   BCMBIAOutageBucketDTO    `json:"outage_bucket"`
	IsSkipped      bool                     `json:"is_skipped"`
	Users          []StaffAllocationUserDTO `json:"users"`
}

type StaffAllocationUserDTO struct {
	ID       uint                             `json:"id"`
	UserID   uint                             `json:"user_id"`
	User     *organizationDTOs.UserMinimalDTO `json:"user"`
	IsRemote bool                             `json:"is_remote"`
}
