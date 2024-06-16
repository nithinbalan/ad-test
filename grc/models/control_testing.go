package models

import (
	"time"
)

type ControlTest struct {
	ID                            uint       `json:"id" gorm:"primaryKey;autoIncrement" binding:"-"`
	ControlID                     uint       `json:"control_id" gorm:"not null" binding:"required"`
	ControlImplementationStatusID uint       `json:"control_implementation_status_id" gorm:"not null" binding:"required"`
	ControlEfficiencyStatusID     uint       `json:"control_efficiency_status_id" gorm:"not null" binding:"required"`
	Title                         string     `json:"title" gorm:"type:nvarchar(500);not null" binding:"required"`
	Description                   *string    `json:"description,omitempty" gorm:"type:nvarchar(MAX);default:null" binding:"omitempty"`
	Score                         *uint      `json:"score,omitempty" gorm:"default:null" binding:"omitempty"`
	Comments                      *string    `json:"comments,omitempty" gorm:"type:nvarchar(MAX);default:null" binding:"omitempty"`
	TestedSource                  string     `json:"tested_source" gorm:"type:varchar(45);not null" binding:"required"`
	TestedAt                      *time.Time `json:"tested_at,omitempty" gorm:"type:date;default:null" binding:"omitempty"`
	TestedBy                      uint       `json:"tested_by" gorm:"not null" binding:"required"`
	UpdatedAt                     *time.Time `json:"updated_at,omitempty" gorm:"type:datetime" binding:"omitempty"`
	UpdatedByID                   *uint      `json:"updated_by_id,omitempty" gorm:"default:null" binding:"omitempty"`

	// Foreign key relationships
	Control                     *Control                     `json:"control" gorm:"foreignKey:ControlID"`
	ControlImplementationStatus *ControlImplementationStatus `json:"control_implementation_status" gorm:"foreignKey:ControlImplementationStatusID"`
	ControlEfficiencyStatus     *ControlEfficiencyStatus     `json:"control_efficiency_status" gorm:"foreignKey:ControlEfficiencyStatusID"`
	TestedByUser                *User                        `json:"tested_by_user" gorm:"foreignKey:TestedBy"`
	UpdatedByUser               *User                        `json:"updated_by_id_user" gorm:"foreignKey:UpdatedByID"`
	Documents                   []*ControlTestDocument       `json:"documents" gorm:"foreignKey:ControlTestID"`
}

// TableName overrides the table name used by User to `control_tests`
func (ControlTest) TableName() string {
	return "control_tests"
}
