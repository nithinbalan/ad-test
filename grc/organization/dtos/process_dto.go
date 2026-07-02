package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	settingsDTOs "jethur/grc/settings/dtos"

	"time"
)

// ProcessListDTO defines the data transfer object for listing processes
type ProcessListDTO struct {
	ID                               uint                 `json:"id"`
	ProcessGroupTitleEn              string               `json:"process_group_title_en"`
	ProcessGroupTitleAr              string               `json:"process_group_title_ar"`
	OwnerUnitTitleEn                 *string              `json:"owner_unit_title_en"`
	OwnerUnitTitleAr                 *string              `json:"owner_unit_title_ar"`
	OwnerID                          *uint                `json:"owner_id"`
	OwnerFirstNameEn                 *string              `json:"owner_first_name_en"`
	OwnerFirstNameAr                 *string              `json:"owner_first_name_ar"`
	OwnerLastNameEn                  *string              `json:"owner_last_name_en"`
	OwnerLastNameAr                  *string              `json:"owner_last_name_ar"`
	OwnerInitial                     string               `json:"owner_initial"`
	OwnerColor                       string               `json:"owner_color"`
	OwnerEmail                       string               `json:"owner_email"`
	OwnerProfileImagePath            string               `json:"owner_profile_image_path"`
	LifecycleStatusTitleEn           string               `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr           string               `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode         string               `json:"lifecycle_status_color_code"`
	ReferenceNumber                  string               `json:"reference_number"`
	TitleEn                          string               `json:"title_en"`
	TitleAr                          string               `json:"title_ar"`
	Description                      *string              `json:"description"`
	MarkedAsAuditableItem            bool                 `json:"marked_as_auditable_item"`
	Rto                              *uint                `json:"rto"`
	RtoUnitId                        *uint                `json:"rto_uint_id"`
	RtoUnitTitleEn                   *string              `json:"rto_unit_title_en"`
	RtoUnitTitleAr                   *string              `json:"rto_unit_title_ar"`
	RtoCondition                     *string              `json:"rto_condition"`
	Mtpd                             *uint                `json:"mtpd"`
	MtpdUnitId                       *uint                `json:"mtpd_uint_id"`
	MtpdUnitTitleEn                  *string              `json:"mtpd_unit_title_en"`
	MtpdUnitTitleAr                  *string              `json:"mtpd_unit_title_ar"`
	MtpdCondition                    *string              `json:"mtpd_condition"`
	Rpo                              *uint                `json:"rpo"`
	RpoUnitId                        *uint                `json:"rpo_uint_id"`
	RpoUnitTitleEn                   *string              `json:"rpo_unit_title_en"`
	RpoUnitTitleAr                   *string              `json:"rpo_unit_title_ar"`
	RpoCondition                     *string              `json:"rpo_condition"`
	ImpactAnalysisLevelId            *uint                `json:"impact_analysis_level_id"`
	ImpactAnalysisLevelTitleEn       *string              `json:"impact_analysis_level_title_en"`
	ImpactAnalysisLevelTitleAr       *string              `json:"impact_analysis_level_title_ar"`
	ImpactAnalysisLevelColorCode     *string              `json:"impact_analysis_level_color_code"`
	ImpactAnalysisLevelType          *string              `json:"impact_analysis_level_type"`
	RiskLevelId                      *uint                `json:"risk_level_id"`
	RiskLevelTitleEn                 *string              `json:"risk_level_title_en"`
	RiskLevelTitleAr                 *string              `json:"risk_level_title_ar"`
	RiskLevelColorCode               *string              `json:"risk_level_color_code"`
	RiskLevelType                    *string              `json:"risk_level_type"`
	LastAssessedAt                   *time.Time           `json:"last_assessed_at"`
	LastAssessedUserID               *uint                `json:"last_assessed_user_id"`
	LastAssessedUserFirstNameEn      *string              `json:"last_assessed_user_first_name_en"`
	LastAssessedUserFirstNameAr      *string              `json:"last_assessed_user_first_name_ar"`
	LastAssessedUserLastNameEn       *string              `json:"last_assessed_user_last_name_en"`
	LastAssessedUserLastNameAr       *string              `json:"last_assessed_user_last_name_ar"`
	LastAssessedUserInitial          *string              `json:"last_assessed_user_initial"`
	LastAssessedUserColor            *string              `json:"last_assessed_user_color"`
	LastAssessedUserEmail            *string              `json:"last_assessed_user_email"`
	LastAssessedUserProfileImagePath *string              `json:"last_assessed_user_profile_image_path"`
	IsBiaPerformed                   bool                 `json:"is_bia_performed"`
	ProcessCriticalityLevelID        uint                 `json:"process_criticality_level_id"`
	ProcessCriticalityLevelTitleEn   string               `json:"process_criticality_level_title_en"`
	ProcessCriticalityLevelTitleAr   string               `json:"process_criticality_level_title_ar"`
	CreatedUserID                    uint                 `json:"created_user_id"`
	CreatedUserFirstNameEn           string               `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr           string               `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn            string               `json:"created_user_last_name_en"`
	CreatedUserLastNameAr            string               `json:"created_user_last_name_ar"`
	CreatedUserInitial               string               `json:"created_user_initial"`
	CreatedUserColor                 string               `json:"created_user_color"`
	CreatedUserEmail                 string               `json:"created_user_email"`
	CreatedUserProfileImagePath      string               `json:"created_user_profile_image_path"`
	UpdatedUserID                    *uint                `json:"updated_user_id"`
	UpdatedUserFirstNameEn           *string              `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr           *string              `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn            *string              `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr            *string              `json:"updated_user_last_name_ar"`
	UpdatedUserInitial               string               `json:"updated_user_initial"`
	UpdatedUserColor                 string               `json:"updated_user_color"`
	UpdatedUserEmail                 string               `json:"updated_user_email"`
	UpdatedUserProfileImagePath      string               `json:"updated_user_profile_image_path"`
	CreatedAt                        time.Time            `json:"created_at"`
	UpdatedAt                        *time.Time           `json:"updated_at"`
	RiskRatingTitleEn                *string              `json:"risk_rating_title_en"`
	RiskRatingTitleAr                *string              `json:"risk_rating_title_ar"`
	RiskRatingColorCode              *string              `json:"risk_rating_color_code"`
	RiskRatingValue                  *uint                `json:"risk_rating_value"`
	Activities                       []ProcessActivityDTO `json:"activities" gorm:"-"`

	// New BCM fields
	DivisionID                            *uint                            `json:"division_id"`
	DivisionTitleEn                       *string                          `json:"division_title_en"`
	DivisionTitleAr                       *string                          `json:"division_title_ar"`
	IsAutomatic                           bool                             `json:"is_automatic"`
	IsManual                              bool                             `json:"is_manual"`
	ProcessTypeID                         *uint                            `json:"process_type_id"`
	ProcessTypeTitleEn                    *string                          `json:"process_type_title_en"`
	ProcessTypeTitleAr                    *string                          `json:"process_type_title_ar"`
	LocationID                            *uint                            `json:"location_id"`
	LocationTitleEn                       *string                          `json:"location_title_en"`
	LocationTitleAr                       *string                          `json:"location_title_ar"`
	AlternativeLocation                   *string                          `json:"alternative_location"`
	RushHour                              *string                          `json:"rush_hour"`
	ReviewFrequencyID                     *uint                            `json:"review_frequency_id"`
	ReviewFrequencyTitleEn                *string                          `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr                *string                          `json:"review_frequency_title_ar"`
	IsCanBeDoneManually                   bool                             `json:"is_can_be_done_manually"`
	TimeToPerformProcedureManually        *string                          `json:"time_to_perform_procedure_manually"`
	StepsToPerformProcedureManually       *string                          `json:"steps_to_perform_procedure_manually"`
	RequirementToPerformProcedureManually *string                          `json:"requirement_to_perform_procedure_manually"`
	ExecutionRecords                      []*ProcessExecutionRecordDTO     `json:"execution_records" gorm:"-"`
	ManualProcedureForms                  []*ProcessManualProcedureFormDTO `json:"manual_procedure_forms" gorm:"-"`
	SystemRegister                        *ProcessSystemRegisterDTO        `json:"system_register" gorm:"-"`
}

// ProcessListQueryParamsDTO defines the filters and search parameters for querying processes
type ProcessListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
	DateFrom    string
	DateTo      string
}

// ProcessListResponseDTO encapsulates the paginated response for a list of processes
type ProcessListResponseDTO struct {
	Items      []*ProcessListDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ProcessMinimalListResponseDTO struct {
	Items      []*ProcessMinimalDTO         `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ProcessActivityDTO struct {
	ID                uint       `json:"id"`
	TitleEn           string     `json:"title_en"`
	TitleAr           string     `json:"title_ar"`
	DescriptionEn     *string    `json:"description_en"`
	DescriptionAr     *string    `json:"description_ar"`
	Sequence          int        `json:"sequence"`
	EstimatedDuration *string    `json:"estimated_duration"`
	ExpectedOutput    string     `json:"expected_output"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
}

// ProcessFormDTO defines the structure for creating or updating a process
type ProcessFormDTO struct {
	ProcessGroupID   *uint   `json:"process_group_id"`
	OwnerUnitID      uint    `json:"owner_unit_id" binding:"required"`
	OwnerID          *uint   `json:"owner_id"`
	TitleEn          string  `json:"title_en" binding:"required"`
	TitleAr          string  `json:"title_ar" binding:"required"`
	Description      *string `json:"description" binding:"omitempty"`
	ControlIds       []*uint `json:"control_ids"`
	AssetRegisterIds []*uint `json:"asset_register_ids"`
	DocumentIds      []*uint `json:"document_ids"`

	// New BCM form fields
	DivisionID                            *uint   `json:"division_id"`
	ExecutionAutomation                   *string `json:"execution_automation"` // "automatic" or "manual"
	ProcessTypeID                         *uint   `json:"process_type_id"`
	LocationID                            *uint   `json:"location_id"`
	AlternativeLocation                   *string `json:"alternative_location"`
	RushHour                              *string `json:"rush_hour"`
	ReviewFrequencyID                     *uint   `json:"review_frequency_id"`
	IsCanBeDoneManually                   *bool   `json:"is_can_be_done_manually"`
	TimeToPerformProcedureManually        *string `json:"time_to_perform_procedure_manually"`
	StepsToPerformProcedureManually       *string `json:"steps_to_perform_procedure_manually"`
	RequirementToPerformProcedureManually *string `json:"requirement_to_perform_procedure_manually"`
	ExecutionRecordIDs                    []*uint `json:"execution_record_ids"`
	ManualProcedureFormIDs                []*uint `json:"manual_procedure_form_ids"`
	SystemRegisterID                      *uint   `json:"system_register_id"`
}

type ProcessDetailDTO struct {
	ID                       uint                              `json:"id"`
	ReferenceNumber          string                            `json:"reference_number"`
	TitleEn                  string                            `json:"title_en"`
	TitleAr                  string                            `json:"title_ar"`
	Description              *string                           `json:"description"`
	ProcessGroup             *ProcessGroupDTO                  `json:"process_group"`
	ResponsibleUnit          *settingsDTOs.OrganizationUnitDTO `json:"owner_unit"`
	ResponsibleUserID        *uint                             `json:"owner_id"`
	ResponsibleUser          *UserMinimalDTO                   `json:"owner"`
	CreatedByUser            UserMinimalDTO                    `json:"created_by"`
	UpdatedByUser            *UserMinimalDTO                   `json:"updated_by"`
	LifecycleStatus          masterDTOs.LifecycleStatusDTO     `json:"status"`
	Rto                      *int                              `json:"rto"`
	RtoUnitId                *uint                             `json:"rto_uint_id"`
	RtoCondition             *string                           `json:"rto_condition"`
	RtoUnit                  BusinessContinuityDurationUnitDTO `json:"rto_unit"`
	Mtpd                     *int                              `json:"mtpd"`
	MtpdUnitId               *uint                             `json:"mtpd_uint_id"`
	MtpdCondition            *string                           `json:"mtpd_condition"`
	MtpdUnit                 BusinessContinuityDurationUnitDTO `json:"mtpd_unit"`
	Rpo                      *int                              `json:"rpo"`
	RpoUnitId                *uint                             `json:"rpo_uint_id"`
	RpoCondition             *string                           `json:"rpo_condition"`
	RpoUnit                  BusinessContinuityDurationUnitDTO `json:"rpo_unit"`
	RiskLevelId              *uint                             `json:"risk_level_id"`
	ImpactAnalysisLevelId    *uint                             `json:"impact_analysis_level_id"`
	ImpactAnalysisLevel      *BCMImpactAnalysisLevelDTO        `json:"impact_analysis_level"`
	LastAssessedAt           *time.Time                        `json:"last_assessed_at"`
	LastAssessedByID         *uint                             `json:"last_assessed_by_id"`
	LastAssessedByUser       *UserMinimalDTO                   `json:"last_assessed_by"`
	IsBiaPerformed           bool                              `json:"is_bia_performed"`
	CreatedAt                time.Time                         `json:"created_at"`
	UpdatedAt                *time.Time                        `json:"updated_at"`
	Controls                 []*ProcessControlDTO              `json:"controls"`
	Assets                   []*ProcessAssetDTO                `json:"assets"`
	Documents                []*ProcessAssociatedDocumentDTO   `json:"documents"`
	Risks                    []*ProcessRiskDTO                 `json:"risks"`
	MarkedAsAuditableItem    bool                              `json:"marked_as_auditable_item"`
	RiskMatrixRating         *RiskMatrixRiskRatingDTO          `json:"risk_matrix_risk_rating"`
	RiskMatrixRatingValue    *RiskMatrixRiskRatingValueDTO     `json:"risk_matrix_rating_value"`
	AuditableItem            *AuditableItemMinimalDetailDTO    `json:"auditable_item"`
	RiskLevel                *RiskLevelDTO                     `json:"risk_level"`
	ImpactJustification      *string                           `json:"impact_justification"`
	CriticalityJustification *string                           `json:"criticality_justification"`
	ProcessCriticalityLevel  *ProcessCriticalityLevelDTO       `json:"process_criticality_level"`
	ProcessActivities        []ProcessActivityDetailDTO        `json:"process_activities"`

	// New BCM fields
	DivisionID                            *uint                              `json:"division_id"`
	Division                              *settingsDTOs.OrganizationUnitDTO  `json:"division"`
	IsAutomatic                           bool                               `json:"is_automatic"`
	IsManual                              bool                               `json:"is_manual"`
	ProcessTypeID                         *uint                              `json:"process_type_id"`
	ProcessType                           *ProcessTypeDTO                    `json:"process_type"`
	LocationID                            *uint                              `json:"location_id"`
	Location                              *masterDTOs.LocationDTO            `json:"location"`
	AlternativeLocation                   *string                            `json:"alternative_location"`
	RushHour                              *string                            `json:"rush_hour"`
	ReviewFrequencyID                     *uint                              `json:"review_frequency_id"`
	ReviewFrequency                       *masterDTOs.ReviewFrequencyDTO     `json:"review_frequency"`
	IsCanBeDoneManually                   bool                               `json:"is_can_be_done_manually"`
	TimeToPerformProcedureManually        *string                            `json:"time_to_perform_procedure_manually"`
	StepsToPerformProcedureManually       *string                            `json:"steps_to_perform_procedure_manually"`
	RequirementToPerformProcedureManually *string                            `json:"requirement_to_perform_procedure_manually"`
	ProcessExecutionRecords               []*ProcessExecutionRecordDetailDTO `json:"process_execution_records"`
	ProcessManualProcedureForms           []*ProcessManualProcedureDetailDTO `json:"process_manual_procedure_forms"`
	ProcessSystemRegister                 *ProcessSystemRegisterDetailDTO    `json:"process_system_register"`
}
type AuditableItemMinimalDetailDTO struct {
	ID                  uint                             `json:"id"`
	Type                string                           `json:"type"`
	TitleEn             *string                          `json:"title_en"`
	TitleAr             *string                          `json:"title_ar"`
	Comment             *string                          `json:"comment"`
	ResponsibleUnitID   uint                             `json:"responsible_unit_id"`
	ResponsibleUserID   uint                             `json:"responsible_user_id"`
	CreatedByID         uint                             `json:"created_by_id"`
	UpdatedByID         *uint                            `json:"updated_by_id"`
	AuditorRiskRatingID uint                             `json:"auditor_risk_rating_id"`
	AuditorRiskRating   *AuditorRiskRatingMinimalDTO     `json:"auditor_risk_rating"`
	ResponsibleUnit     settingsDTOs.OrganizationUnitDTO `json:"responsible_unit"`
	CreatedAt           time.Time                        `json:"created_at"`
	UpdatedAt           *time.Time                       `json:"updated_at"`
	ResponsibleUser     *UserMinimalDTO                  `json:"responsible_user"`
	CreatedByUser       *UserMinimalDTO                  `json:"created_by"`
	UpdatedByUser       *UserMinimalDTO                  `json:"updated_by"`
}
type AuditorRiskRatingMinimalDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	Status    string `json:"status"`
}
type RiskMatrixRiskRatingDTO struct {
	ID           uint   `json:"id"`
	RiskMatrixID uint   `json:"risk_matrix_id"`
	TitleEn      string `json:"title_en"`
	TitleAr      string `json:"title_ar"`
	ColorCode    string `json:"color_code"`
	Type         string `json:"type"`
	Status       string `json:"status"`
}

type ProcessControlDTO struct {
	ID        *uint            `json:"id"`
	ProcessID uint             `json:"process_id"`
	ControlID uint             `json:"control_id"`
	Control   ControlDetailDTO `json:"control"`
}

type ProcessRiskDTO struct {
	ID             *uint          `json:"id"`
	ProcessID      uint           `json:"process_id"`
	RiskRegisterID uint           `json:"risk_register_id"`
	Risk           RiskMinimalDTO `json:"risk"`
}

type RiskMinimalDTO struct {
	ID                       uint                          `json:"id"`
	Title                    string                        `json:"title"`
	Description              *string                       `json:"description"`
	CurrentRiskRatingID      *uint                         `json:"current_risk_rating_id"`
	CurrentRiskRatingValueID *uint                         `json:"current_risk_rating_value_id"`
	CurrentRiskRatingValue   *RiskMatrixRiskRatingValueDTO `json:"current_risk_rating_value"`
}

type RiskMatrixRiskRatingValueDTO struct {
	ID                     uint `json:"id"`
	Value                  int  `json:"value"` // This is used to find the maximum rating
	RiskMatrixRiskRatingID uint `json:"risk_matrix_risk_rating_id"`
}

type ProcessAssociatedDocumentDTO struct {
	ID         *uint              `json:"id"`
	ProcessID  uint               `json:"process_id"`
	DocumentID uint               `json:"document_id"`
	Document   DocumentDetailsDTO `json:"document"`
}

type DocumentDetailsDTO struct {
	ID            int     `json:"id"`
	TitleEn       string  `json:"title_en"`
	TitleAr       string  `json:"title_ar"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
	IsLocked      bool    `json:"is_locked"`
	Acknowledged  *bool   `json:"acknowledged,omitempty"`
	Rto           *int    `json:"rto"`
	RtoUnitId     *uint   `json:"rto_uint_id"`
	Mtpd          *int    `json:"mtpd"`
	MtpdUnitId    *uint   `json:"mtpd_uint_id"`
	Rpo           *int    `json:"rpo"`
	RpoUnitId     *uint   `json:"rpo_uint_id"`
}

type ProcessAssetDTO struct {
	ID              *uint           `json:"id"`
	ProcessID       uint            `json:"process_id"`
	AssetRegisterID uint            `json:"asset_register_id"`
	Asset           AssetMinimalDTO `json:"asset"`
}

type AssetMinimalDTO struct {
	ID          uint    `json:"id"`
	TitleEn     string  `json:"title_en"`
	TitleAr     string  `json:"title_ar"`
	Description *string `json:"description"`
	Rto         *int    `json:"rto"`
	RtoUnitId   *uint   `json:"rto_uint_id"`
	Mtpd        *int    `json:"mtpd"`
	MtpdUnitId  *uint   `json:"mtpd_uint_id"`
	Rpo         *int    `json:"rpo"`
	RpoUnitId   *uint   `json:"rpo_uint_id"`
}

type RiskLevelDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	Type      string `json:"type"`
	ColorCode string `json:"color_code"`
}

// ProcessTypeDTO inline minimal DTO for process type (avoids circular import with bcm package)
type ProcessTypeDTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

// ProcessExecutionRecordDTO represents a linked execution record in process detail
type ProcessExecutionRecordDTO struct {
	ID                   uint   `json:"id"`
	ProcessID            uint   `json:"process_id"`
	BCMExecutionRecordID uint   `json:"bcm_execution_record_id"`
	TitleEn              string `json:"title_en"`
	TitleAr              string `json:"title_ar"`
}

// ProcessManualProcedureFormDTO represents a linked manual procedure form in process detail
type ProcessManualProcedureFormDTO struct {
	ID                       uint    `json:"id"`
	ProcessID                uint    `json:"process_id"`
	BCMManualProcedureFormID uint    `json:"bcm_manual_procedure_form_id"`
	TitleEn                  string  `json:"title_en"`
	TitleAr                  string  `json:"title_ar"`
	Code                     *string `json:"code"`
}

// ProcessSystemRegisterDTO represents a linked system register in process detail
type ProcessSystemRegisterDTO struct {
	ID               uint   `json:"id"`
	ProcessID        uint   `json:"process_id"`
	SystemRegisterID uint   `json:"system_register_id"`
	TitleEn          string `json:"title_en"`
	TitleAr          string `json:"title_ar"`
}

type ProcessExecutionRecordDetailDTO struct {
	ID                   *uint                        `json:"id"`
	ProcessID            uint                         `json:"process_id"`
	BCMExecutionRecordID uint                         `json:"bcm_execution_record_id"`
	BCMExecutionRecord   BCMExecutionRecordMinimalDTO `json:"bcm_execution_record"`
}

type ProcessManualProcedureDetailDTO struct {
	ID                       *uint                            `json:"id"`
	ProcessID                uint                             `json:"process_id"`
	BCMManualProcedureFormID uint                             `json:"bcm_manual_procedure_form_id"`
	BCMManualProcedureForm   BCMManualProcedureFormMinimalDTO `json:"bcm_manual_procedure_form"`
}

type ProcessSystemRegisterDetailDTO struct {
	ID               *uint                    `json:"id"`
	ProcessID        uint                     `json:"process_id"`
	SystemRegisterID uint                     `json:"system_register_id"`
	SystemRegister   SystemRegisterMinimalDTO `json:"system_register"`
}

type BCMExecutionRecordMinimalDTO struct {
	ID                uint                       `json:"id"`
	TitleEn           string                     `json:"title_en"`
	TitleAr           string                     `json:"title_ar"`
	OwnerEntity       OrganizationUnitDTO        `json:"owner_entity"`
	RecordType        *BCMExecutionRecordTypeDTO `json:"record_type"`
	IsBackUpAvailable *string                    `json:"is_back_up_available"`
	Location          *masterDTOs.LocationDTO    `json:"location"`
	RecoveryPlan      *string                    `json:"recovery_plan"`
}

type BCMManualProcedureFormMinimalDTO struct {
	ID                  uint                    `json:"id"`
	TitleEn             string                  `json:"title_en"`
	TitleAr             string                  `json:"title_ar"`
	Code                *string                 `json:"code"`
	OwnerEntity         OrganizationUnitDTO     `json:"owner_entity"`
	Location            *masterDTOs.LocationDTO `json:"location"`
	InternalStakeholder *string                 `json:"internal_stakeholder"`
	Notes               *string                 `json:"notes"`
}

type SystemRegisterMinimalDTO struct {
	ID                   uint                               `json:"id"`
	TitleEn              string                             `json:"title_en"`
	TitleAr              string                             `json:"title_ar"`
	ServiceDescription   *string                            `json:"service_description"`
	ModulesComponents    *string                            `json:"modules_components"`
	RTO                  *int                               `json:"rto"`
	RTOUnit              *BusinessContinuityDurationUnitDTO `json:"rto_unit"`
	MTPD                 *int                               `json:"mtpd"`
	MTPDUnit             *BusinessContinuityDurationUnitDTO `json:"mtpd_unit"`
	RPO                  *int                               `json:"rpo"`
	RPOUnit              *BusinessContinuityDurationUnitDTO `json:"rpo_unit"`
	RecoveryManager      *DesignationDTO                    `json:"recovery_manager"`
	RecoveryManagerPhone string                             `json:"recovery_manager_phone"`
	RecoveryManagerEmail string                             `json:"recovery_manager_email"`
	RecoveryLocation     *string                            `json:"recovery_location"`
}

type BusinessContinuityDurationUnitDTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Minutes int    `json:"minutes"`
	Type    string `json:"type"`
	Status  string `json:"status"`
}

type BCMExecutionRecordTypeDTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type BCMImpactAnalysisLevelDTO struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	Score     int    `json:"score"`
	ColorCode string `json:"color_code"`
	Status    string `json:"status"`
}
