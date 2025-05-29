package models

import (
	"time"
)

type StakeholderNeedsAndExpectation struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	StakeholderID      uint       `json:"stakeholder_id" binding:"required" gorm:"not null"`
	Title              string     `json:"title" binding:"required,max=500" gorm:"type:varchar(500);not null"`
	SeverityLevelID    *uint      `json:"severity_level_id"`
	SuccessCriteria    *string    `json:"success_criteria"`
	Challenges         *string    `json:"challenges"`
	Description        *string    `json:"description"`
	EvaluationStatusID *uint      `json:"evaluation_status_id" gorm:"not null"`
	Score              *int       `json:"score" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	EvaluatedAt        *time.Time `json:"evaluated_at" binding:"omitempty" gorm:"type:date,default:null"`
	EvaluatedBYID      *uint      `json:"evaluated_by_id,omitempty" gorm:"default:null"`
	CreatedAt          time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`
	CreatedBYID        uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID        *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`

	// Relationships
	Stakeholder     *Stakeholder                                `gorm:"foreignKey:StakeholderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel   *SeverityLevel                              `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User                                       `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser   *User                                       `gorm:"foreignKey:UpdatedBYID"`
	EvaluatedByUser *User                                       `gorm:"foreignKey:EvaluatedBYID"`
	Evaluations     *[]StakeholderNeedsAndExpectationEvaluation `gorm:"foreignKey:StakeholderNeedsAndExpectationID"`
}

// TableName overrides the table name used by StakeholderNeedsAndExpectation to `stakeholder_needs_and_expectations`
func (StakeholderNeedsAndExpectation) TableName() string {
	return "stakeholder_needs_and_expectations"
}
