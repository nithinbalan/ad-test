package models

import (
	"time"
)

type TestAndExerciseTestResultDocument struct {
	ID                          uint      `json:"id" gorm:"primary_key;auto_increment"`
	TestAndExerciseTestResultID uint      `json:"test_and_exercise_test_result_id"`
	URL                         string    `json:"url"`
	Name                        string    `json:"name" binding:"required"`
	Size                        int       `json:"size"`
	Token                       string    `json:"token"`
	OriginalName                string    `json:"original_name"`
	Type                        string    `json:"type"`
	CreatedByID                 uint      `json:"created_by_id" gorm:"not null"`
	UpdatedByID                 uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                   time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                   time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	TestAndExerciseTestResult *TestAndExerciseTestResult `gorm:"foreignKey:TestAndExerciseTestResultID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser             *User                      `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser             *User                      `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by Test and exercise test result documents to `test_and_exercise_test_result_documents`
func (TestAndExerciseTestResultDocument) TableName() string {
	return "test_and_exercise_test_result_documents"
}
