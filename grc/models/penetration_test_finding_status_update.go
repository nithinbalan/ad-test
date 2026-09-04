package models

import (
	"time"
)

type PenetrationTestFindingStatusUpdate struct {
	ID                       uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	PenetrationTestFindingID uint       `json:"penetration_test_finding_id" gorm:"not null"`
	ActionPlanStatusID       uint       `json:"action_plan_status_id" gorm:"not null"`
	Percentage               *uint      `json:"percentage" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	Comment                  string     `json:"comment" binding:"required" gorm:"type:nvarchar(max);not null"`
	CreatedBYID              uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID              *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt                *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	PenetrationTestFinding *PenetrationTestFinding                       `json:"penetration_test_finding,omitempty" gorm:"foreignKey:PenetrationTestFindingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus       *ActionPlanStatus                             `json:"action_plan_status,omitempty" gorm:"foreignKey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                                         `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser          *User                                         `json:"updated_by_user,omitempty" gorm:"foreignKey:UpdatedBYID"`
	Documents              []*PenetrationTestFindingStatusUpdateDocument `json:"documents,omitempty" gorm:"foreignKey:PenetrationTestFindingStatusUpdateID"`
}

// TableName overrides the table name used by PenetrationTestFindingStatusUpdate to `penetration_test_finding_status_updates`
func (PenetrationTestFindingStatusUpdate) TableName() string {
	return "penetration_test_finding_status_updates"
}
