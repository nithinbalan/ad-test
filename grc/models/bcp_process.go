package models

import (
	"time"
)

type BcpProcess struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BusinessContinuityPlanID uint       `json:"incident_register_id" binding:"required"`
	ProcessID                uint       `json:"process_id" binding:"required"`
	RiskLevelID              *uint      `json:"risk_level_id"`
	ImpactAnalysisLevelID    *uint      `json:"impact_analysis_level_id" gorm:"default:null"`
	RtoFromValue             *int       `json:"rto_from_value" gorm:"type:int;null"`
	RtoToValue               *int       `json:"rto_to_value" gorm:"type:int;null"`
	RtoUnitId                *uint      `json:"rto_unit_id" gorm:"default:null"`
	RtoCondition             *string    `json:"rto_condition" gorm:"type:varchar(20);null"`
	MtpdFromValue            *int       `json:"mtpd_from_value" gorm:"type:int;null"`
	MtpdToValue              *int       `json:"mtpd_to_value" gorm:"type:int;null"`
	MtpdUnitId               *uint      `json:"mtpd_unit_id" gorm:"default:null"`
	MtpdCondition            *string    `json:"mtpd_condition" gorm:"type:varchar(20);null"`
	RpoFromValue             *int       `json:"rpo_from_value" gorm:"type:int;null"`
	RpoToValue               *int       `json:"rpo_to_value" gorm:"type:int;null"`
	RpoUnitId                *uint      `json:"rpo_unit_id" gorm:"default:null"`
	RpoCondition             *string    `json:"rpo_condition" gorm:"type:varchar(20);null"`
	RspFromValue             *int       `json:"rsp_from_value" gorm:"type:int;null"`
	RspToValue               *int       `json:"rsp_to_value" gorm:"type:int;null"`
	RspCondition             *string    `json:"rsp_condition" gorm:"type:varchar(20);null"`
	RspUnitID                *uint      `json:"rsp_unit_id" gorm:"null"`
	CreatedBYID              uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	BusinessContinuityPlan *BusinessContinuityPlan         `gorm:"foreignKey:BusinessContinuityPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Process                *Process                        `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskLevel              *RiskLevel                      `gorm:"foreignKey:RiskLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RtoUnit                *BusinessContinuityDurationUnit `gorm:"foreignKey:RtoUnitId"`
	MtpdUnit               *BusinessContinuityDurationUnit `gorm:"foreignKey:MtpdUnitId"`
	RpoUnit                *BusinessContinuityDurationUnit `gorm:"foreignKey:RpoUnitId"`
	RspUnit                *BusinessContinuityDurationUnit `gorm:"foreignKey:RspUnitID"`
	CreatedByUser          *User                           `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser          *User                           `gorm:"foreignKey:UpdatedBYID"`
	ImpactAnalysisLevel    *BCMImpactAnalysisLevel         `gorm:"foreignKey:ImpactAnalysisLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// BCP Process associations
	ExecutionRecords     []*BcpProcessExecutionRecord     `gorm:"foreignKey:BcpProcessID;references:ID"`
	ManualProcedureForms []*BcpProcessManualProcedureForm `gorm:"foreignKey:BcpProcessID;references:ID"`
	SystemRegisters      []*BcpProcessSystemRegister      `gorm:"foreignKey:BcpProcessID;references:ID"`
	InternalDependencies []*BcpProcessInternalDependency  `gorm:"foreignKey:BcpProcessID;references:ID"`
	ExternalDependencies []*BcpProcessExternalDependency  `gorm:"foreignKey:BcpProcessID;references:ID"`
	EmployeeAnalyses     []*BcpProcessEmployeeAnalysis    `gorm:"foreignKey:BcpProcessID;references:ID"`
	StaffAllocations     []*BcpProcessStaffAllocation     `gorm:"foreignKey:BcpProcessID;references:ID"`
}

func (BcpProcess) TableName() string {
	return "bcp_processes"
}
