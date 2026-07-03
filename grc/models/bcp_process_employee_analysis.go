package models

import (
	"time"
)

type BcpProcessEmployeeAnalysis struct {
	ID            uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BcpProcessID  uint      `json:"bcp_process_id" gorm:"not null"`
	DesignationID uint      `json:"designation_id" gorm:"not null"`
	NormalCount   *int      `json:"normal_count" gorm:"default:null"`
	CriticalCount *int      `json:"critical_count" gorm:"default:null"`
	CreatedBYID   *uint     `json:"created_by_id" gorm:"default:null"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	BcpProcess    *BcpProcess                          `gorm:"foreignKey:BcpProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Designation   *Designation                         `gorm:"foreignKey:DesignationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User                                `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Employees     []BcpProcessEmployeeAnalysisEmployee `gorm:"foreignKey:BcpProcessEmployeeAnalysisID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BcpProcessEmployeeAnalysis) TableName() string {
	return "bcp_process_employee_analysis"
}
