package models

import (
	"time"
)

type TaeRecoveryActionPlanComment struct {
	ID                                  uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	TestAndExerciseRecoveryActionPlanID uint      `json:"test_and_exercise_recovery_action_plan_id" gorm:"not null"`
	Comment                             string    `json:"comment" binding:"required"`
	CreatedBYID                         uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                           time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	TaeRecoveryActionPlan *TestAndExerciseRecoveryActionPlan      `gorm:"foreignKey:TestAndExerciseRecoveryActionPlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser         *User                                   `gorm:"foreignKey:CreatedBYID"`
	Documents             []*TaeRecoveryActionPlanCommentDocument `gorm:"foreignKey:TaeRecoveryActionPlanCommentID"`
}

func (TaeRecoveryActionPlanComment) TableName() string {
	return "tae_recovery_action_plan_comments"
}
