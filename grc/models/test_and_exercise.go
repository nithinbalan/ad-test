package models

import (
	"time"
)

type TestAndExercise struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber          string     `json:"reference_number" binding:"omitempty,max=45" gorm:"size:45"`
	BusinessContinuityPlanID uint       `json:"business_continuity_id" binding:"required"`
	BcpTestAndExerciseID     uint       `json:"bcp_test_and_exercise_id" binding:"required"`
	ResponsibleUserID        *uint      `json:"responsible_user_id" gorm:"not null" binding:"required"`
	ResponsibleUnitID        *uint      `json:"responsible_unit_id" gorm:"not null" binding:"required"`
	TestingDate              *time.Time `json:"testing_date"`
	Description              *string    `json:"description" gorm:"type:nvarchar(max)"`
	TestAndExerciseStatusID  uint       `json:"test_and_exercise_status_id" binding:"required" gorm:"not null"`
	CreatedBYID              uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	BusinessContinuityPlan                  *BusinessContinuityPlan                   `gorm:"foreignKey:BusinessContinuityPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser                         *User                                     `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit                         *OrganizationUnit                         `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	BcpTestAndExercise                      *BcpTestAndExercise                       `gorm:"foreignKey:BcpTestAndExerciseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	TestAndExerciseStatus                   *TestAndExerciseStatus                    `gorm:"foreignKey:TestAndExerciseStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	TestAndExerciseProcesses                []*TestAndExerciseProcess                 `gorm:"foreignKey:TestAndExerciseID"`
	TestAndExerciseRecoveryActionPlans      []*TestAndExerciseRecoveryActionPlan      `gorm:"foreignKey:TestAndExerciseID"`
	TestAndExerciseCommunicationActionPlans []*TestAndExerciseCommunicationActionPlan `gorm:"foreignKey:TestAndExerciseID"`
	TestAndExerciseTestResult               *TestAndExerciseTestResult                `gorm:"foreignKey:TestAndExerciseID"`
	CreatedByUser                           *User                                     `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser                           *User                                     `gorm:"foreignKey:UpdatedBYID"`
}

func (TestAndExercise) TableName() string {
	return "test_and_exercises"
}
