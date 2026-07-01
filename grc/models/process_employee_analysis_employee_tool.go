package models

import (
	"time"
)

type ProcessEmployeeAnalysisEmployeeTool struct {
	ID                                uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProcessEmployeeAnalysisEmployeeID uint      `json:"process_employee_analysis_employee_id" gorm:"not null"`
	EquipmentToolID                   uint      `json:"equipment_tool_id" gorm:"not null"`
	IsHas                             bool      `json:"is_has" gorm:"type:bit;not null;default:0"`
	IsNeeded                          bool      `json:"is_needed" gorm:"type:bit;not null;default:0"`
	CreatedBYID                       *uint     `json:"created_by_id" gorm:"default:null"`
	CreatedAt                         time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ProcessEmployeeAnalysisEmployee *ProcessEmployeeAnalysisEmployee `gorm:"foreignKey:ProcessEmployeeAnalysisEmployeeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	EquipmentTool                   *ProcessEquipmentTool            `gorm:"foreignKey:EquipmentToolID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser                   *User                            `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProcessEmployeeAnalysisEmployeeTool) TableName() string {
	return "process_employee_analysis_employee_tools"
}
