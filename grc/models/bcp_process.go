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
	Rto                      *int       `json:"rto"`
	RtoUnitId                *uint      `json:"rto_unit_id"`
	RtoCondition             *string    `json:"rto_condition" gorm:"type:varchar(20);null"`
	Mtpd                     *int       `json:"mtpd"`
	MtpdUnitId               *uint      `json:"mtpd_unit_id"`
	MtpdCondition            *string    `json:"mtpd_condition" gorm:"type:varchar(20);null"`
	Rpo                      *int       `json:"rpo"`
	RpoUnitId                *uint      `json:"rpo_unit_id"`
	RpoCondition             *string    `json:"rpo_condition" gorm:"type:varchar(20);null"`
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
