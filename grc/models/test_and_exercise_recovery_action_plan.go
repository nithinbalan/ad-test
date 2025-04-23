package models

import (
	"time"
)

type TestAndExerciseRecoveryActionPlan struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TestAndExerciseID  uint       `json:"bcp_test_and_exercise_id" binding:"required"`
	BcpPhaseID         *uint      `json:"bcp_phase_id"`
	ProcessID          uint       `json:"process_id" binding:"required"`
	SeverityLevelID    *uint      `json:"severity_level_id"`
	Title              string     `json:"title"`
	ResponsibleUserID  uint       `json:"responsible_user_id" gorm:"not null" binding:"required"`
	ResponsibleUnitID  uint       `json:"responsible_unit_id" gorm:"not null" binding:"required"`
	TargetDate         *time.Time `json:"target_date"`
	ActionPlanStatusId *uint      `json:"action_plan_status_id"`
	IsRecovery         bool       `json:"is_recovery"`
	Percentage         *int       `json:"percentage"`
	ReviewUserID       *uint      `json:"review_user_id" gorm:"default:null"`
	CreatedBYID        uint       `json:"created_by_id" binding:"required"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID        *uint      `json:"updated_by_id"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	TestAndExercise         *TestAndExercise                     `gorm:"foreignKey:TestAndExerciseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	BcpPhase                *BcpPhase                            `gorm:"foreignKey:BcpPhaseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Process                 *Process                             `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel           *SeverityLevel                       `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus        *ActionPlanStatus                    `gorm:"foreignKey:ActionPlanStatusId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser         *User                                `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit         *OrganizationUnit                    `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewUser              *User                                `gorm:"foreignKey:ReviewUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"` // New foreign key relationship
	ActionPlanStatusUpdates []*TaeRecoveryActionPlanStatusUpdate `gorm:"foreignKey:TestAndExerciseRecoveryActionPlanID"`
	CreatedByUser           *User                                `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser           *User                                `gorm:"foreignKey:UpdatedBYID"`
}

func (TestAndExerciseRecoveryActionPlan) TableName() string {
	return "test_and_exercise_recovery_action_plans"
}
