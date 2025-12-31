package models

import "time"

type TPRMAssessmentSectionQuestionComment struct {
	ID                                     uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TprmAssessmentSectionQuestionID        *uint      `json:"tprm_assessment_section_question_id"`
	TprmAssessmentSectionQuestionCommentID *uint      `json:"tprm_assessment_section_question_comment_id"`
	ThirdPartyAssessmentID                 uint       `json:"third_party_assessment_id"`
	Comment                                string     `json:"comment" binding:"required" gorm:"type:nvarchar(max);not null"`
	Level                                  int        `json:"level" gorm:"not null"`
	CreatedBYID                            uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                            *uint      `json:"updated_by_id"`
	CreatedAt                              time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt                              *time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Question             *TPRMAssessmentSectionQuestion          `gorm:"foreignKey:TprmAssessmentSectionQuestionID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	ParentComment        *TPRMAssessmentSectionQuestionComment   `gorm:"foreignKey:TprmAssessmentSectionQuestionCommentID"`
	ChildComments        []*TPRMAssessmentSectionQuestionComment `gorm:"foreignKey:TprmAssessmentSectionQuestionCommentID"`
	ThirdPartyAssessment *ThirdPartyAssessment                   `gorm:"foreignKey:ThirdPartyAssessmentID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	CreatedByUser        *User                                   `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser        *User                                   `gorm:"foreignKey:UpdatedBYID"`
}

func (TPRMAssessmentSectionQuestionComment) TableName() string {
	return "tprm_assessment_section_question_comments"
}
