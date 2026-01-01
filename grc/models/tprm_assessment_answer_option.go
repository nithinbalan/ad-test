package models

import (
	"time"
)

type TPRMAssessmentAnswerOption struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TPRMAssessmentQuestionID uint       `json:"tprm_assessment_question_id" binding:"required" gorm:"not null"`
	AnswerOptionID           uint       `json:"answer_option_id" binding:"required" gorm:"not null"`
	TPRMAssessmentAnswerID   uint       `json:"tprm_assessment_answer_id" binding:"required" gorm:"not null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime"`
	TPRMComplianceResponseID *uint      `json:"tprm_compliance_response_id"`

	// Relationships
	TPRMAssessmentQuestion *TPRMAssessmentSectionQuestion      `gorm:"foreignKey:TPRMAssessmentQuestionID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	AnswerOption           *TPRMAssessmentQuestionAnswerOption `gorm:"foreignKey:AnswerOptionID"`
	TPRMAssessmentAnswer   *TPRMAssessmentAnswer               `gorm:"foreignKey:TPRMAssessmentAnswerID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	TPRMComplianceResponse *TPRMComplianceResponse             `gorm:"foreignKey:TPRMComplianceResponseID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
}

// TableName overrides the table name used by TPRMAssessmentAnswerOption to `tprm_assessment_answer_options`
func (TPRMAssessmentAnswerOption) TableName() string {
	return "tprm_assessment_answer_options"
}
