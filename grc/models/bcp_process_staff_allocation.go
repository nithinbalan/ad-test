package models

import (
	"time"
)

type BcpProcessStaffAllocation struct {
	ID                           uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BcpProcessID                 uint      `json:"bcp_process_id" gorm:"not null"`
	BcpProcessEmployeeAnalysisID *uint     `json:"bcp_process_employee_analysis_id" gorm:"default:null"`
	CreatedBYID                  *uint     `json:"created_by_id" gorm:"default:null"`
	CreatedAt                    time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	BcpProcess                 *BcpProcess                        `gorm:"foreignKey:BcpProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	BcpProcessEmployeeAnalysis *BcpProcessEmployeeAnalysis        `gorm:"foreignKey:BcpProcessEmployeeAnalysisID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser              *User                              `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	OutageBuckets              []BcpProcessAllocationOutageBucket `gorm:"foreignKey:BcpProcessAllocationID;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BcpProcessStaffAllocation) TableName() string {
	return "bcp_process_staff_allocations"
}
