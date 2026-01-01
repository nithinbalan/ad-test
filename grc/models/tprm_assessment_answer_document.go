package models

import "time"

type TPRMAssessmentAnswerDocument struct {
	ID                       uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TPRMAssessmentAnswerID   uint       `json:"assessment_plan_answer_id" gorm:"not null"`
	TPRMAssessmentQuestionID uint       `json:"tprm_assessment_section_question_id" gorm:"not null"`
	URL                      string     `json:"url" gorm:"type:varchar(1000)"`
	Name                     *string    `json:"name" gorm:"type:nvarchar(255)"`
	Size                     *int       `json:"size"`
	Token                    *string    `json:"token" gorm:"type:varchar(255)"`
	OriginalName             *string    `json:"original_name" gorm:"type:nvarchar(255)"`
	Type                     *string    `json:"type" gorm:"type:varchar(255)"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime"`

	TPRMAssessmentAnswer   *TPRMAssessmentAnswer          `gorm:"foreignKey:TPRMAssessmentAnswerID;constraint:OnDelete:CASCADE,OnUpdate:NO ACTION"`
	TPRMAssessmentQuestion *TPRMAssessmentSectionQuestion `gorm:"foreignKey:TPRMAssessmentQuestionID;constraint:OnDelete:CASCADE,OnUpdate:NO ACTION"`
}

func (TPRMAssessmentAnswerDocument) TableName() string {
	return "tprm_assessment_answer_documents"
}
