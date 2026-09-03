package models

import (
	"time"
)

type TestResultPerformanceMeasure struct {
	ID             uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TestResultID   uint       `json:"test_result_id" binding:"required"`
	Criterion      *string    `json:"criterion" gorm:"type:nvarchar(max)"`
	Target         *string    `json:"target" gorm:"type:nvarchar(max)"`
	AchievedResult *string    `json:"achieved_result" gorm:"type:nvarchar(max)"`
	Remarks        *string    `json:"remarks" gorm:"type:nvarchar(max)"`
	CreatedBYID    uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID    *uint      `json:"updated_by_id"`
	CreatedAt      time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt      *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	TestResult    *TestAndExerciseTestResult `gorm:"foreignKey:TestResultID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User                      `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User                      `gorm:"foreignKey:UpdatedBYID"`
}

func (TestResultPerformanceMeasure) TableName() string {
	return "test_result_performance_measures"
}
