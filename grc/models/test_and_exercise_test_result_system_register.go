package models

import (
	"time"
)

type TestAndExerciseTestResultSystemRegister struct {
	ID                            uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TestAndExerciseTestResultID   uint       `json:"test_and_exercise_test_result_id" binding:"required" gorm:"not null;index"`
	SystemRegisterID              uint       `json:"system_register_id" binding:"required" gorm:"not null"`
	TestCompletionDate            *time.Time `json:"test_completion_date"`
	CommunicationFeedbackOptionID *uint      `json:"communication_feedback_option_id"`
	RecoveryFeedbackOptionID      *uint      `json:"recovery_feedback_option_id"`
	TestResultOptionID            *uint      `json:"test_result_option_id"`
	CreatedBYID                   uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt                     time.Time  `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	TestAndExerciseTestResult   *TestAndExerciseTestResult `gorm:"foreignKey:TestAndExerciseTestResultID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SystemRegister              *SystemRegister            `gorm:"foreignKey:SystemRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CommunicationFeedbackOption *FeedbackOption            `gorm:"foreignKey:CommunicationFeedbackOptionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RecoveryFeedbackOption      *FeedbackOption            `gorm:"foreignKey:RecoveryFeedbackOptionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	TestResultOption            *TestResultOption          `gorm:"foreignKey:TestResultOptionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser               *User                      `gorm:"foreignKey:CreatedBYID"`
}

func (TestAndExerciseTestResultSystemRegister) TableName() string {
	return "test_and_exercise_test_result_system_registers"
}
