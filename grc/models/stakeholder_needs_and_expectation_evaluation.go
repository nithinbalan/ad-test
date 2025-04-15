package models

import (
	"time"
)

type StakeholderNeedsAndExpectationEvaluation struct {
	ID                               uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	StakeholderNeedsAndExpectationID uint       `json:"risk_action_plan_id" gorm:"not null"`
	EvaluationStatusID               uint       `json:"evaluation_status_id" gorm:"not null"`
	Score                            *int       `json:"score" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	Comment                          *string    `json:"comment" binding:"required" gorm:"type:nvarchar(max);not null"`
	EvaluatedAt                      *time.Time `json:"evaluated_at" binding:"omitempty" gorm:"type:date,default:null"`
	EvaluatedBYID                    uint       `json:"evaluated_by_id"`
	CreatedAt                        time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                        *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID                      uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                      *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	StakeholderNeedsAndExpectation *StakeholderNeedsAndExpectation `gorm:"foreignKey:StakeholderNeedsAndExpectationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	EvaluationStatus               *EvaluationStatus               `gorm:"foreignKey:EvaluationStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                  *User                           `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser                  *User                           `gorm:"foreignKey:UpdatedBYID"`
	EvaluatedByUser                *User                           `gorm:"foreignKey:EvaluatedBYID"`
}

// TableName overrides the table name used by StakeholderNeedsAndExpectationEvaluation to `stakeholder_needs_and_expectation_evaluations`
func (StakeholderNeedsAndExpectationEvaluation) TableName() string {
	return "stakeholder_needs_and_expectation_evaluations"
}
