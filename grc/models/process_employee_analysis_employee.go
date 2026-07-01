package models

import (
	"time"
)

type ProcessEmployeeAnalysisEmployee struct {
	ID                        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProcessEmployeeAnalysisID uint      `json:"process_employee_analysis_id" gorm:"not null"`
	UserID                    uint      `json:"user_id" gorm:"not null"`
	CreatedBYID               *uint     `json:"created_by_id" gorm:"default:null"`
	CreatedAt                 time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ProcessEmployeeAnalysis *ProcessEmployeeAnalysis              `gorm:"foreignKey:ProcessEmployeeAnalysisID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	User                    *User                                 `gorm:"foreignKey:UserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser           *User                                 `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Tools                   []ProcessEmployeeAnalysisEmployeeTool `gorm:"foreignKey:ProcessEmployeeAnalysisEmployeeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProcessEmployeeAnalysisEmployee) TableName() string {
	return "process_employee_analysis_employees"
}
