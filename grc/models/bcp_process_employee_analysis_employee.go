package models

import (
	"time"
)

type BcpProcessEmployeeAnalysisEmployee struct {
	ID                           uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BcpProcessEmployeeAnalysisID uint      `json:"bcp_process_employee_analysis_id" gorm:"not null"`
	UserID                       uint      `json:"user_id" gorm:"not null"`
	CreatedBYID                  *uint     `json:"created_by_id" gorm:"default:null"`
	CreatedAt                    time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	BcpProcessEmployeeAnalysis *BcpProcessEmployeeAnalysis              `gorm:"foreignKey:BcpProcessEmployeeAnalysisID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	User                       *User                                    `gorm:"foreignKey:UserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser              *User                                    `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Tools                      []BcpProcessEmployeeAnalysisEmployeeTool `gorm:"foreignKey:BcpProcessEmployeeAnalysisEmployeeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BcpProcessEmployeeAnalysisEmployee) TableName() string {
	return "bcp_process_employee_analysis_employees"
}
