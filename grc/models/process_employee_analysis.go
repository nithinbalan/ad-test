package models

import (
	"time"
)

type ProcessEmployeeAnalysis struct {
	ID            uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ProcessID     uint       `json:"process_id" gorm:"not null"`
	DesignationID uint       `json:"designation_id" gorm:"not null"`
	NormalCount   *int       `json:"normal_count" gorm:"default:null"`
	CriticalCount *int       `json:"critical_count" gorm:"default:null"`
	CreatedBYID   *uint      `json:"created_by_id" gorm:"default:null"`
	UpdatedBYID   *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	Process       *Process                          `gorm:"foreignKey:ProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Designation   *Designation                      `gorm:"foreignKey:DesignationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User                             `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser *User                             `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Employees     []ProcessEmployeeAnalysisEmployee `gorm:"foreignKey:ProcessEmployeeAnalysisID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProcessEmployeeAnalysis) TableName() string {
	return "process_employee_analysis"
}
