package dtos

import (
	incidentDTOs "jethur/grc/incident/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type BusinessImpactAssessmentListDTO struct {
	ID                          uint    `json:"id"`
	ProcessID                   *uint   `json:"process_id"`
	ProcessTitleEn              string  `json:"process_title_en"`
	ProcessTitleAr              string  `json:"process_title_ar"`
	ImpactCategoryID            uint    `json:"impact_category_id"`
	ImpactCategoryTitleEn       string  `json:"impact_category_title_en"`
	ImpactCategoryTitleAr       string  `json:"impact_category_title_ar"`
	ImpactCategoryType          string  `json:"impact_category_type"`
	RiskLevelID                 *uint   `json:"risk_level_id"`
	RiskLevelTitleEn            string  `json:"risk_level_title_en"`
	RiskLevelTitleAr            string  `json:"risk_level_title_ar"`
	RiskLevelColorCode          string  `json:"risk_level_color_code"`
	RiskLevelType               string  `json:"risk_level_type"`
	Cost                        *uint   `json:"cost"`
	Comments                    *string `json:"comments"`
	IsApplicable                bool    `json:"is_applicable"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   string  `json:"updated_at"`
	CreatedUserID               uint    `json:"created_user_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
}

type BusinessImpactAssessmentRecoveryObjectiveListDTO struct {
	ID             uint    `json:"id"`
	ProcessID      uint    `json:"process_id"`
	ProcessTitleEn *string `json:"process_title_en"`
	ProcessTitleAr *string `json:"process_title_ar"`

	ProcessDependencySubTypeID      uint    `json:"process_dependency_sub_type_id"`
	ProcessDependencySubTypeType    *string `json:"process_dependency_sub_type_type"`
	ProcessDependencySubTypeTitleEn *string `json:"process_dependency_sub_type_title_en"`
	ProcessDependencySubTypeTitleAr *string `json:"process_dependency_sub_type_title_ar"`
	ProcessDependencyTypeTitleEn    *string `json:"process_dependency_type_title_en"`
	ProcessDependencyTypeTitleAr    *string `json:"process_dependency_type_title_ar"`
	ProcessDependencyTypeType       *string `json:"process_dependency_type_type"`

	DependentProcessID      *uint   `json:"dependent_process_id"`
	DependentProcessTitleEn *string `json:"dependent_process_title_en"`
	DependentProcessTitleAr *string `json:"dependent_process_title_ar"`

	OrganizationUnitID      *uint   `json:"organization_unit_id"`
	OrganizationUnitTitleEn *string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr *string `json:"organization_unit_title_ar"`

	AssetID      *uint   `json:"asset_id"`
	AssetTitleEn *string `json:"asset_title_en"`
	AssetTitleAr *string `json:"asset_title_ar"`

	Title       string  `json:"title"`
	Description *string `json:"description"`

	Rto               *int    `json:"rto"`
	RtoUnitId         *uint   `json:"rto_unit_id"`
	RtoUnitTitleEn    *string `json:"rto_unit_title_en"`
	RtoUnitTitleAr    *string `json:"rto_unit_title_ar"`
	Mtpd              *int    `json:"mtpd"`
	MtpdUnitId        *uint   `json:"mtpd_unit_id"`
	MtpdUnitTitleEn   *string `json:"mtpd_unit_title_en"`
	MtpdUnitTitleAr   *string `json:"mtpd_unit_title_ar"`
	Rpo               *int    `json:"rpo"`
	RpoUnitId         *uint   `json:"rpo_unit_id"`
	RpoUnitTitleEn    *string `json:"rpo_unit_title_en"`
	RpoUnitTitleAr    *string `json:"rpo_unit_title_ar"`
	ImpactDescription *string `json:"impact_description"`

	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   string  `json:"updated_at"`
	CreatedUserID               uint    `json:"created_user_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
}

type CriticalSystemListDTO struct {
	SystemRegisterID      uint                       `json:"system_register_id"`
	SystemRegisterTitleEn string                     `json:"system_register_title_en"`
	SystemRegisterTitleAr string                     `json:"system_register_title_ar"`
	ServiceDescription    *string                    `json:"service_description"`
	RTO                   *int                       `json:"rto"`
	RTOUnitTitleEn        *string                    `json:"rto_unit_title_en"`
	RTOUnitTitleAr        *string                    `json:"rto_unit_title_ar"`
	MTPD                  *int                       `json:"mtpd"`
	MTPDUnitTitleEn       *string                    `json:"mtpd_unit_title_en"`
	MTPDUnitTitleAr       *string                    `json:"mtpd_unit_title_ar"`
	RPO                   *int                       `json:"rpo"`
	RPOUnitTitleEn        *string                    `json:"rpo_unit_title_en"`
	RPOUnitTitleAr        *string                    `json:"rpo_unit_title_ar"`
	RecoveryManagerPhone  string                     `json:"recovery_manager_phone"`
	RecoveryManagerEmail  string                     `json:"recovery_manager_email"`
	RecoveryLocation      *string                    `json:"recovery_location"`
	VendorRTO             *int                       `json:"vendor_rto"`
	VendorRTOUnitTitleEn  *string                    `json:"vendor_rto_unit_title_en"`
	VendorRTOUnitTitleAr  *string                    `json:"vendor_rto_unit_title_ar"`
	CriticalProcessCount  int                        `json:"critical_process_count"`
	CriticalProcesses     []CriticalSystemProcessDTO `json:"critical_processes"`
}

type CriticalSystemProcessDTO struct {
	ProcessID                    uint    `json:"process_id"`
	ProcessTitleEn               string  `json:"process_title_en"`
	ProcessTitleAr               string  `json:"process_title_ar"`
	ImpactAnalysisLevelID        *uint   `json:"impact_analysis_level_id"`
	ImpactAnalysisLevelTitleEn   *string `json:"impact_analysis_level_title_en"`
	ImpactAnalysisLevelTitleAr   *string `json:"impact_analysis_level_title_ar"`
	ImpactAnalysisLevelColorCode *string `json:"impact_analysis_level_color_code"`
	RtoFromValue                 *int    `json:"rto_from_value"`
	RtoToValue                   *int    `json:"rto_to_value"`
	RtoCondition                 *string `json:"rto_condition"`
	RtoUnitTitleEn               *string `json:"rto_unit_title_en"`
	RtoUnitTitleAr               *string `json:"rto_unit_title_ar"`
	MtpdFromValue                *int    `json:"mtpd_from_value"`
	MtpdToValue                  *int    `json:"mtpd_to_value"`
	MtpdCondition                *string `json:"mtpd_condition"`
	MtpdUnitTitleEn              *string `json:"mtpd_unit_title_en"`
	MtpdUnitTitleAr              *string `json:"mtpd_unit_title_ar"`
	RpoFromValue                 *int    `json:"rpo_from_value"`
	RpoToValue                   *int    `json:"rpo_to_value"`
	RpoCondition                 *string `json:"rpo_condition"`
	RpoUnitTitleEn               *string `json:"rpo_unit_title_en"`
	RpoUnitTitleAr               *string `json:"rpo_unit_title_ar"`
}

type BusinessImpactAssessmentListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Filters     map[string][]string
	QuickFilter []string
}

type BusinessImpactAssessmentRecoveryObjectiveListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Filters map[string][]string
}

type BusinessImpactFormDTO struct {
	ProcessID *uint `json:"process_id"`
}

type BusinessImpactAssessmentFormDTO struct {
	ProcessID                 *uint                         `json:"process_id"`
	BusinessImpactAssessments []BusinessImpactAssessmentDTO `json:"business_impact_assessments"`
	// New fields
	OverallImpactLevelID *uint   `json:"overall_impact_level_id"`
	ImpactJustification  *string `json:"impact_justification"`
}

type BusinessImpactAssessmentDTO struct {
	ImpactCategoryID                 uint                                  `json:"impact_category_id"`
	IsApplicable                     bool                                  `json:"is_applicable"`
	BusinessImpactAssessmentAnalysis []BusinessImpactAssessmentAnalysisDTO `json:"business_impact_assessment_analysis"`
}

type BusinessImpactAssessmentAnalysisDTO struct {
	OutageBucketID        *uint `json:"outage_bucket_id"`
	ImpactAnalysisLevelID *uint `json:"impact_analysis_level_id"`
}

type BusinessImpactAssessmentRecoveryObjectiveFormDTO struct {
	ProcessID          *uint                                          `json:"process_id"`
	Rto                *int                                           `json:"rto"`
	RtoUnitId          *uint                                          `json:"rto_unit_id"`
	Mtpd               *int                                           `json:"mtpd"`
	MtpdUnitId         *uint                                          `json:"mtpd_unit_id"`
	Rpo                *int                                           `json:"rpo"`
	RpoUnitId          *uint                                          `json:"rpo_unit_id"`
	RecoveryObjectives []BusinessImpactAssessmentRecoveryObjectiveDTO `json:"recovery_objectives"`
}

type BusinessImpactAssessmentProcessCriticalityFormDTO struct {
	ProcessID                 *uint   `json:"process_id"`
	ProcessCriticalityLevelID *uint   `json:"process_criticality_level_id"`
	CriticalityJustification  *string `json:"criticality_justification"`
}

type BusinessImpactAssessmentRecoveryObjectiveDTO struct {
	ProcessDependencyID *uint `json:"process_dependency_id"`
	Rto                 *int  `json:"rto"`
	RtoUnitId           *uint `json:"rto_unit_id"`
	Mtpd                *int  `json:"mtpd"`
	MtpdUnitId          *uint `json:"mtpd_unit_id"`
	Rpo                 *int  `json:"rpo"`
	RpoUnitId           *uint `json:"rpo_unit_id"`
}

type BIAMatrixResponseDTO struct {
	BusinessFunctionID                 uint    `json:"business_function_id"`
	BusinessFunctionTitle              string  `json:"business_function_title"`
	BusinessFunctionCriticalityTitleEn *string `json:"business_function_criticality_title_en"`
	BusinessFunctionCriticalityTitleAr *string `json:"business_function_criticality_title_ar"`
	BusinessFunctionCriticalityColor   *string `json:"business_function_criticality_color"`

	ProcessID                 uint    `json:"process_id"`
	ProcessTitleEn            string  `json:"process_title_en"`
	ProcessTitleAr            string  `json:"process_title_ar"`
	ProcessCriticalityTitleEn *string `json:"process_criticality_title_en"`
	ProcessCriticalityTitleAr *string `json:"process_criticality_title_ar"`
	ProcessCriticalityColor   *string `json:"process_criticality_color"`

	ImpactAnalysisLevelID        *uint   `json:"impact_analysis_level_id"`
	ImpactAnalysisLevelTitleEn   *string `json:"impact_analysis_level_title_en"`
	ImpactAnalysisLevelTitleAr   *string `json:"impact_analysis_level_title_ar"`
	ImpactAnalysisLevelColorCode *string `json:"impact_analysis_level_color_code"`
	ImpactAnalysisLevelType      *string `json:"impact_analysis_level_type"`

	ProcessRTOFromValue    *int    `json:"process_rto_from_value"`
	ProcessRTOToValue      *int    `json:"process_rto_to_value"`
	ProcessRTOCondition    *string `json:"process_rto_condition"`
	ProcessRTOUnitID       *uint   `json:"process_rto_unit_id"`
	ProcessRTOUnitTitleEn  *string `json:"process_rto_unit_title_en"`
	ProcessRTOUnitTitleAr  *string `json:"process_rto_unit_title_ar"`
	ProcessMTPDFromValue   *int    `json:"process_mtpd_from_value"`
	ProcessMTPDToValue     *int    `json:"process_mtpd_to_value"`
	ProcessMTPDCondition   *string `json:"process_mtpd_condition"`
	ProcessMTPDUnitID      *uint   `json:"process_mtpd_unit_id"`
	ProcessMTPDUnitTitleEn *string `json:"process_mtpd_unit_title_en"`
	ProcessMTPDUnitTitleAr *string `json:"process_mtpd_unit_title_ar"`
	ProcessRPOFromValue    *int    `json:"process_rpo_from_value"`
	ProcessRPOToValue      *int    `json:"process_rpo_to_value"`
	ProcessRPOCondition    *string `json:"process_rpo_condition"`
	ProcessRPOUnitID       *uint   `json:"process_rpo_unit_id"`
	ProcessRPOUnitTitleEn  *string `json:"process_rpo_unit_title_en"`
	ProcessRPOUnitTitleAr  *string `json:"process_rpo_unit_title_ar"`
	ProcessRSPFromValue    *int    `json:"process_rsp_from_value"`
	ProcessRSPToValue      *int    `json:"process_rsp_to_value"`
	ProcessRSPCondition    *string `json:"process_rsp_condition"`
	ProcessRSPUnitID       *uint   `json:"process_rsp_unit_id"`
	ProcessRSPUnitTitleEn  *string `json:"process_rsp_unit_title_en"`
	ProcessRSPUnitTitleAr  *string `json:"process_rsp_unit_title_ar"`

	RiskLevelID        *uint   `json:"risk_level_id"`
	RiskLevelTitleEn   *string `json:"risk_level_title_en"`
	RiskLevelTitleAr   *string `json:"risk_level_title_ar"`
	RiskLevelColorCode *string `json:"risk_level_color_code"`
}

type BIAMatrixQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	Export      bool
	QuickFilter []string
}

type BIAProcessDetailDTO struct {
	ID                              uint                                                 `json:"id"`
	ReferenceNumber                 string                                               `json:"reference_number"`
	TitleEn                         string                                               `json:"title_en"`
	TitleAr                         string                                               `json:"title_ar"`
	Description                     *string                                              `json:"description"`
	ProcessGroup                    *organizationDTOs.ProcessGroupDTO                    `json:"process_group"`
	ResponsibleUnit                 *settingsDTOs.OrganizationUnitDTO                    `json:"owner_unit"`
	ResponsibleUserID               *uint                                                `json:"owner_id"`
	ResponsibleUser                 *organizationDTOs.UserMinimalDTO                     `json:"owner"`
	CreatedByUser                   organizationDTOs.UserMinimalDTO                      `json:"created_by"`
	UpdatedByUser                   *organizationDTOs.UserMinimalDTO                     `json:"updated_by"`
	LifecycleStatus                 masterDTOs.LifecycleStatusDTO                        `json:"status"`
	Rto                             *int                                                 `json:"rto"`
	RtoUnitId                       *uint                                                `json:"rto_uint_id"`
	RtoUnit                         masterDTOs.TimeUnitDTO                               `json:"rto_unit"`
	Mtpd                            *int                                                 `json:"mtpd"`
	MtpdUnitId                      *uint                                                `json:"mtpd_uint_id"`
	MtpdUnit                        masterDTOs.TimeUnitDTO                               `json:"mtpd_unit"`
	Rpo                             *int                                                 `json:"rpo"`
	RpoUnitId                       *uint                                                `json:"rpo_uint_id"`
	RpoUnit                         masterDTOs.TimeUnitDTO                               `json:"rpo_unit"`
	RiskLevelId                     *uint                                                `json:"risk_level_id"`
	RiskMatrixRating                *organizationDTOs.RiskMatrixRiskRatingDTO            `json:"risk_matrix_risk_rating"`
	RiskMatrixRatingValue           *organizationDTOs.RiskMatrixRiskRatingValueDTO       `json:"risk_matrix_rating_value"`
	LastAssessedAt                  *time.Time                                           `json:"last_assessed_at"`
	LastAssessedByID                *uint                                                `json:"last_assessed_by_id"`
	LastAssessedByUser              *organizationDTOs.UserMinimalDTO                     `json:"last_assessed_by"`
	IsBiaPerformed                  bool                                                 `json:"is_bia_performed"`
	CreatedAt                       time.Time                                            `json:"created_at"`
	UpdatedAt                       *time.Time                                           `json:"updated_at"`
	MarkedAsAuditableItem           bool                                                 `json:"marked_as_auditable_item"`
	RiskLevel                       *organizationDTOs.RiskLevelDTO                       `json:"risk_level"`
	ImpactJustification             *string                                              `json:"impact_justification"`
	CriticalityJustification        *string                                              `json:"criticality_justification"`
	ProcessCriticalityLevel         *organizationDTOs.ProcessCriticalityLevelDTO         `json:"process_criticality_level"`
	BusinessImpactAssessments       []*BusinessImpactAssessmentDetailDTO                 `json:"business_impact_assessments"`
	BusinessImpactAssessmentRatings []*BusinessImpactAssessmentRatingDTO                 `json:"business_impact_assessment_ratings"`
	ProcessDependencies             []*organizationDTOs.ProcessDependencyDetailDTO       `json:"process_dependencies"`
	Risks                           []*ProcessRiskDTO                                    `json:"risks"`
	Incidents                       []*ProcessIncidentDTO                                `json:"incidents"`
	ProcessEmployeeAnalysis         []*organizationDTOs.ProcessEmployeeAnalysisDetailDTO `json:"process_employee_analysis"`
	BCMStaffAllocations             []*BCMStaffAllocationDetailDTO                       `json:"bcm_staff_allocations"`
}

type ProcessRiskDTO struct {
	ID             *uint          `json:"id"`
	ProcessID      uint           `json:"process_id"`
	RiskRegisterID uint           `json:"risk_register_id"`
	Risk           RiskMinimalDTO `json:"risk"`
}

type RiskMinimalDTO struct {
	ID                       uint                                   `json:"id"`
	Title                    string                                 `json:"title"`
	Description              *string                                `json:"description"`
	RiskCategory             *riskDTOs.RiskCategoryDTO              `json:"risk_category"`
	RiskSource               *riskDTOs.RiskSourceDTO                `json:"risk_source"`
	CurrentLikelihood        *riskDTOs.RiskMatrixLikelihoodDTO      `json:"current_likelihood"`
	CurrentImpact            *riskDTOs.RiskMatrixImpactDTO          `json:"current_impact"`
	CurrentRiskRating        *riskDTOs.RiskMatrixRiskRatingDTO      `json:"current_risk_rating"`
	CurrentRiskRatingValue   *riskDTOs.RiskMatrixRiskRatingValueDTO `json:"current_risk_rating_value"`
	CurrentRiskRatingID      *uint                                  `json:"current_risk_rating_id"`
	CurrentRiskRatingValueID *uint                                  `json:"current_risk_rating_value_id"`
}

type ProcessIncidentDTO struct {
	ID                 *uint              `json:"id"`
	ProcessID          uint               `json:"process_id"`
	IncidentRegisterID uint               `json:"incident_register_id"`
	IncidentRegister   IncidentMinimalDTO `json:"incident_register"`
}

type IncidentMinimalDTO struct {
	ID               uint                              `json:"id"`
	Title            string                            `json:"title"`
	Description      *string                           `json:"description"`
	IncidentAt       time.Time                         `json:"incident_at"`
	Location         *string                           `json:"location"`
	ReportedAt       *time.Time                        `json:"reported_at"`
	IncidentCategory *incidentDTOs.IncidentCategoryDTO `json:"incident_category"`
	IncidentType     *incidentDTOs.IncidentTypeDTO     `json:"incident_type"`
}

type BusinessImpactAssessmentDetailDTO struct {
	ImpactCategory                   *masterDTOs.ImpactCategoryDTO               `json:"impact_category"`
	IsApplicable                     bool                                        `json:"is_applicable"`
	BusinessImpactAssessmentAnalysis []BusinessImpactAssessmentAnalysisDetailDTO `json:"business_impact_assessment_analysis"`
}

type BusinessImpactAssessmentAnalysisDetailDTO struct {
	OutageBucketID        *uint                      `json:"outage_bucket_id"`
	ImpactAnalysisLevelID *uint                      `json:"impact_analysis_level_id"`
	OutageBucket          *BCMBIAOutageBucketDTO     `json:"outage_bucket"`
	ImpactAnalysisLevel   *BCMImpactAnalysisLevelDTO `json:"impact_analysis_level"`
}

type BusinessImpactAssessmentRatingFormDTO struct {
	RtoFromValue   *int    `json:"rto_from_value"`
	RtoToValue     *int    `json:"rto_to_value"`
	MtpdFromValue  *int    `json:"mtpd_from_value"`
	MtpdToValue    *int    `json:"mtpd_to_value"`
	RpoFromValue   *int    `json:"rpo_from_value"`
	RpoToValue     *int    `json:"rpo_to_value"`
	RspFromValue   *int    `json:"rsp_from_value"`
	RspToValue     *int    `json:"rsp_to_value"`
	OverrideReason *string `json:"override_reason"`
}

type BusinessImpactAssessmentRatingDTO struct {
	ID                   uint                               `json:"id"`
	ProcessID            uint                               `json:"process_id"`
	OverallImpactLevelID uint                               `json:"overall_impact_level_id"`
	OverallImpactLevel   *BCMImpactAnalysisLevelDTO         `json:"overall_impact_level"`
	BCMBIAImpactLevelID  uint                               `json:"bcm_bia_impact_level_id"`
	BCMBIAImpactLevel    *BCMBIAImpactLevelDetailDTO        `json:"bcm_bia_impact_level"`
	IsEdited             bool                               `json:"is_edited"`
	RtoFromValue         *int                               `json:"rto_from_value"`
	RtoToValue           *int                               `json:"rto_to_value"`
	RtoCondition         *string                            `json:"rto_condition"`
	RtoUnit              *BusinessContinuityDurationUnitDTO `json:"rto_unit"`
	MtpdFromValue        *int                               `json:"mtpd_from_value"`
	MtpdToValue          *int                               `json:"mtpd_to_value"`
	MtpdCondition        *string                            `json:"mtpd_condition"`
	MtpdUnit             *BusinessContinuityDurationUnitDTO `json:"mtpd_unit"`
	RpoFromValue         *int                               `json:"rpo_from_value"`
	RpoToValue           *int                               `json:"rpo_to_value"`
	RpoCondition         *string                            `json:"rpo_condition"`
	RpoUnit              *BusinessContinuityDurationUnitDTO `json:"rpo_unit"`
	RspFromValue         *int                               `json:"rsp_from_value"`
	RspToValue           *int                               `json:"rsp_to_value"`
	RspCondition         *string                            `json:"rsp_condition"`
	RspUnit              *BusinessContinuityDurationUnitDTO `json:"rsp_unit"`
	OverrideReason       *string                            `json:"override_reason"`
	CreatedBYID          *uint                              `json:"created_by_id"`
	CreatedAt            time.Time                          `json:"created_at"`
}
