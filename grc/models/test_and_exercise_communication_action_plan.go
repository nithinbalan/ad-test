package models

import (
	"time"
)

type TestAndExerciseCommunicationActionPlan struct {
	ID                uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TestAndExerciseID uint      `json:"bcp_test_and_exercise_id" binding:"required"`
	BcpPhaseID        uint      `json:"bcp_phase_id" binding:"required"`
	ProcessID         uint      `json:"process_id" binding:"required"`
	Title             string    `json:"title"`
	ResponsibleUserID uint      `json:"responsible_user_id" gorm:"not null" binding:"required"`
	ResponsibleUnitID uint      `json:"responsible_unit_id" gorm:"not null" binding:"required"`
	Recipient         *string   `json:"recipient"`
	Frequency         *string   `json:"frequency"`
	Channel           *string   `json:"channel"`
	CreatedBYID       uint      `json:"created_by_id" binding:"required"`
	CreatedAt         time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	TestAndExercise *TestAndExercise  `gorm:"foreignKey:TestAndExerciseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Process         *Process          `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	BcpPhase        *BcpPhase         `gorm:"foreignKey:BcpPhaseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser *User             `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit *OrganizationUnit `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User             `gorm:"foreignKey:CreatedBYID"`
}

func (TestAndExerciseCommunicationActionPlan) TableName() string {
	return "test_and_exercise_communication_action_plans"
}
