package models

import (
	"time"
)

type TestAndExerciseTestResult struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TestAndExerciseID uint       `json:"test_and_exercise_id" binding:"required"`
	TestObservations  *string    `json:"test_observations"`
	Recommendations   *string    `json:"recommendations"`
	Findings          *string    `json:"findings"`
	CreatedBYID       uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID       *uint      `json:"updated_by_id"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	TestAndExercise           *TestAndExercise                          `gorm:"foreignKey:TestAndExerciseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Documents                 []TestAndExerciseTestResultDocument       `gorm:"foreignKey:TestAndExerciseTestResultID"`
	TestResultSystemRegisters []TestAndExerciseTestResultSystemRegister `gorm:"foreignKey:TestAndExerciseTestResultID"`
	CreatedByUser             *User                                     `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser             *User                                     `gorm:"foreignKey:UpdatedBYID"`
}

func (TestAndExerciseTestResult) TableName() string {
	return "test_and_exercise_test_results"
}
