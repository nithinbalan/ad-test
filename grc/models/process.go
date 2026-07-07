package models

import (
	"time"
)

type Process struct {
	ID                          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ProcessGroupID              *uint      `json:"process_group_id" gorm:"default:null"`
	ResponsibleUnitID           uint       `json:"owner_unit_id" gorm:"not null" binding:"required"`
	ResponsibleUserID           *uint      `json:"owner_id" gorm:"default:null"`
	LifecycleStatusID           *uint      `json:"lifecycle_status_id" gorm:"default:null"`
	ReferenceNumber             string     `json:"reference_number" gorm:"type:varchar(45);not null;unique" binding:"required"`
	TitleEn                     string     `json:"title_en" gorm:"type:nvarchar(500);not null" binding:"required"`
	TitleAr                     string     `json:"title_ar" gorm:"type:nvarchar(500);not null" binding:"required"`
	Description                 *string    `json:"description" gorm:"type:nvarchar(MAX);default:null" binding:"omitempty"`
	RtoFromValue                *int       `json:"rto_from_value" gorm:"type:int;null"`
	RtoToValue                  *int       `json:"rto_to_value" gorm:"type:int;null"`
	RtoUnitId                   *uint      `json:"rto_uint_id" gorm:"default:null"`
	RtoCondition                *string    `json:"rto_condition" gorm:"type:varchar(20);null"`
	MtpdFromValue               *int       `json:"mtpd_from_value" gorm:"type:int;null"`
	MtpdToValue                 *int       `json:"mtpd_to_value" gorm:"type:int;null"`
	MtpdUnitId                  *uint      `json:"mtpd_uint_id" gorm:"default:null"`
	MtpdCondition               *string    `json:"mtpd_condition" gorm:"type:varchar(20);null"`
	RpoFromValue                *int       `json:"rpo_from_value" gorm:"type:int;null"`
	RpoToValue                  *int       `json:"rpo_to_value" gorm:"type:int;null"`
	RpoUnitId                   *uint      `json:"rpo_uint_id" gorm:"default:null"`
	RpoCondition                *string    `json:"rpo_condition" gorm:"type:varchar(20);null"`
	RspFromValue                *int       `json:"rsp_from_value" gorm:"type:int;null"`
	RspToValue                  *int       `json:"rsp_to_value" gorm:"type:int;null"`
	RspCondition                *string    `json:"rsp_condition" gorm:"type:varchar(20);null"`
	RspUnitID                   *uint      `json:"rsp_unit_id" gorm:"null"`
	ImpactAnalysisLevelID       *uint      `json:"impact_analysis_level_id" gorm:"default:null"`
	RiskLevelId                 *uint      `json:"risk_level_id" gorm:"default:null"`
	ImpactJustification         *string    `json:"impact_justification" gorm:"default:null"`
	ProcessCriticalityLevelId   *uint      `json:"process_criticality_level_id" gorm:"default:null"`
	CriticalityJustification    *string    `json:"criticality_justification" gorm:"default:null"`
	IsBiaPerformed              bool       `json:"is_bia_performed" gorm:"default:false"`
	LastAssessedAt              *time.Time `json:"last_assessed_at" gorm:"type:datetime; null" binding:"-"`
	LastAssessedByID            *uint      `json:"last_assessed_by_id" gorm:"default:null;column:last_assessed_by_id" binding:"omitempty"`
	CreatedBYID                 uint       `json:"created_by_id" gorm:"not null;column:created_by_id" binding:"required"`
	UpdatedBYID                 *uint      `json:"updated_by_id" gorm:"default:null;column:updated_by_id" binding:"omitempty"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"type:datetime;not null" binding:"-"`
	UpdatedAt                   *time.Time `json:"updated_at" gorm:"type:datetime;default:null" binding:"omitempty"`
	MarkedAsAuditableItem       bool       `json:"marked_as_auditable_item" gorm:"default:false" binding:"-"`
	RiskMatrixRiskRatingID      *uint      `json:"risk_matrix_risk_rating_id" gorm:"default:null"`
	RiskMatrixRiskRatingValueID *uint      `json:"risk_matrix_risk_rating_value_id" gorm:"default:null"`

	// New BCM fields (added via migration)
	DivisionID                            *uint   `json:"division_id" gorm:"default:null"`
	IsAutomatic                           bool    `json:"is_automatic" gorm:"not null;default:0"`
	IsManual                              bool    `json:"is_manual" gorm:"not null;default:0"`
	ProcessTypeID                         *uint   `json:"process_type_id" gorm:"default:null"`
	LocationID                            *uint   `json:"location_id" gorm:"default:null"`
	AlternativeLocation                   *string `json:"alternative_location" gorm:"type:nvarchar(MAX);default:null"`
	RushHour                              *string `json:"rush_hour" gorm:"type:varchar(100);default:null"`
	ReviewFrequencyID                     *uint   `json:"review_frequency_id" gorm:"default:null"`
	IsCanBeDoneManually                   bool    `json:"is_can_be_done_manually" gorm:"not null;default:0"`
	TimeToPerformProcedureManually        *string `json:"time_to_perform_procedure_manually" gorm:"type:varchar(100);default:null"`
	StepsToPerformProcedureManually       *string `json:"steps_to_perform_procedure_manually" gorm:"type:nvarchar(MAX);default:null"`
	RequirementToPerformProcedureManually *string `json:"requirement_to_perform_procedure_manually" gorm:"type:nvarchar(MAX);default:null"`

	// Relationships
	ProcessGroup          *ProcessGroup              `gorm:"foreignKey:ProcessGroupID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ResponsibleUnit       *OrganizationUnit          `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser       *User                      `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	LifecycleStatus       *LifecycleStatus           `gorm:"foreignKey:LifecycleStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskMatrixRating      *RiskMatrixRiskRating      `gorm:"foreignKey:RiskMatrixRiskRatingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixRatingValue *RiskMatrixRiskRatingValue `gorm:"foreignKey:RiskMatrixRiskRatingValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// New BCM relationships
	Division        *OrganizationUnit `gorm:"foreignKey:DivisionID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ProcessType     *BCMProcessType   `gorm:"foreignKey:ProcessTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Location        *Location         `gorm:"foreignKey:LocationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ReviewFrequency *ReviewFrequency  `gorm:"foreignKey:ReviewFrequencyID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`

	CreatedByUser                   *User                             `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser                   *User                             `gorm:"foreignKey:UpdatedBYID"`
	LastAssessedByUser              *User                             `gorm:"foreignKey:LastAssessedByID"`
	Controls                        []*ProcessControl                 `gorm:"foreignKey:ProcessID"`
	Assets                          []*ProcessAsset                   `gorm:"foreignKey:ProcessID"`
	Documents                       []*ProcessAssociatedDocument      `gorm:"foreignKey:ProcessID"`
	Risks                           []*RiskProcess                    `gorm:"foreignKey:ProcessID"`
	Incidents                       []*IncidentProcess                `gorm:"foreignKey:ProcessID"`
	AuditableItems                  []*AuditableItem                  `gorm:"foreignKey:ProcessID"`
	BusinessImpactAssessments       []*BusinessImpactAssessment       `gorm:"foreignKey:ProcessID"`
	ProcessDependencies             []*ProcessDependency              `gorm:"foreignKey:ProcessID"`
	ProcessActivities               []*ProcessActivity                `gorm:"foreignKey:ProcessID"`
	RtoUnit                         *BusinessContinuityDurationUnit   `gorm:"foreignKey:RtoUnitId"`
	MtpdUnit                        *BusinessContinuityDurationUnit   `gorm:"foreignKey:MtpdUnitId"`
	RpoUnit                         *BusinessContinuityDurationUnit   `gorm:"foreignKey:RpoUnitId"`
	RspUnit                         *BusinessContinuityDurationUnit   `gorm:"foreignKey:RspUnitID"`
	RiskLevel                       *RiskLevel                        `gorm:"foreignKey:RiskLevelId"`
	ProcessCriticalityLevel         *ProcessCriticalityLevel          `gorm:"foreignKey:ProcessCriticalityLevelId"`
	BusinessImpactAssessmentRatings []*BusinessImpactAssessmentRating `gorm:"foreignKey:ProcessID"`
	ImpactAnalysisLevel             *BCMImpactAnalysisLevel           `gorm:"foreignKey:ImpactAnalysisLevelID"`

	// Process BCM associations
	ProcessExecutionRecords     []*ProcessExecutionRecord     `gorm:"foreignKey:ProcessID"`
	ProcessManualProcedureForms []*ProcessManualProcedureForm `gorm:"foreignKey:ProcessID"`
	ProcessSystemRegister       *ProcessSystemRegister        `gorm:"foreignKey:ProcessID"`
	ProcessEmployeeAnalysis     []*ProcessEmployeeAnalysis    `gorm:"foreignKey:ProcessID"`
	BCMStaffAllocations         []*BCMStaffAllocation         `gorm:"foreignKey:ProcessID"`
	InternalDependencies        []*ProcessInternalDependency  `gorm:"foreignKey:ProcessID;references:ID"`
	ExternalDependencies        []*ProcessExternalDependency  `gorm:"foreignKey:ProcessID;references:ID"`
}

func (Process) TableName() string {
	return "processes"
}
