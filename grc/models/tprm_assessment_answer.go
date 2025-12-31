package models

import (
	"time"
)

type TPRMAssessmentAnswer struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TPRMAssessmentQuestionID uint       `json:"tprm_assessment_question_id" binding:"required" gorm:"not null"`
	RatingScale              *int       `json:"rating_scale"`
	AnswerText               *string    `json:"answer_text"`
	AnswerDate               *time.Time `json:"answer_date"`
	NotApplicable            *bool      `json:"not_applicable"`
	Comment                  string     `json:"comment"`

	// Relationships
	TPRMAssessmentQuestion *TPRMAssessmentSectionQuestion  `gorm:"foreignKey:TPRMAssessmentQuestionID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	AnswerOptions          []*TPRMAssessmentAnswerOption   `gorm:"foreignKey:TPRMAssessmentAnswerID"`
	Documents              []*TPRMAssessmentAnswerDocument `gorm:"foreignKey:TPRMAssessmentAnswerID"`
}

// TableName overrides the table name used by TPRMAssessmentAnswer to `tprm_assessment_answers`
func (TPRMAssessmentAnswer) TableName() string {
	return "tprm_assessment_answers"
}
