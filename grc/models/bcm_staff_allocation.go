package models

import (
	"time"
)

type BCMStaffAllocation struct {
	ID                        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProcessID                 uint      `json:"process_id" gorm:"not null"`
	ProcessEmployeeAnalysisID *uint     `json:"process_employee_analysis_id" gorm:"default:null"`
	CreatedBYID               uint      `json:"created_by_id"`
	CreatedAt                 time.Time `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	// Relationships
	Process                 *Process                         `gorm:"foreignKey:ProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ProcessEmployeeAnalysis *ProcessEmployeeAnalysis         `gorm:"foreignKey:ProcessEmployeeAnalysisID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser           *User                            `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	OutageBuckets           []BCMStaffAllocationOutageBucket `gorm:"foreignKey:BCMStaffAllocationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BCMStaffAllocation) TableName() string {
	return "bcm_staff_allocations"
}
