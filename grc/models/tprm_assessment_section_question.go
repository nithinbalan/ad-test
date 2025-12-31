package models

import (
	"time"
)

type TPRMAssessmentSectionQuestion struct {
	ID                          uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ThirdPartyAssessmentID      uint       `json:"third_party_assessment_id" binding:"required" gorm:"not null"`
	TPRMAssessmentSectionID     uint       `json:"tprm_assessment_section_id" binding:"required" gorm:"not null"`
	ThirdPartyQuestionLibraryID uint       `json:"third_party_question_library_id" binding:"required" gorm:"not null"`
	QuestionTitleEn             string     `json:"question_title_en" binding:"required,max=1000" gorm:"type:nvarchar(1000);not null"`
	QuestionTitleAr             string     `json:"question_title_ar" binding:"required,max=1000" gorm:"type:nvarchar(1000);not null"`
	Order                       int        `json:"order"`
	CreatedBYID                 uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                 *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                   *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	ThirdPartyAssessment *ThirdPartyAssessment                 `gorm:"foreignKey:ThirdPartyAssessmentID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	ThirdParty           *ThirdPartyAssessment                 `gorm:"foreignKey:ThirdPartyAssessmentID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	AnswerOptions        []*TPRMAssessmentQuestionAnswerOption `gorm:"foreignKey:TPRMAssessmentQuestionID"`
	QuestionLibrary      *ThirdPartyQuestionLibrary            `gorm:"foreignKey:ThirdPartyQuestionLibraryID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	Answer               *TPRMAssessmentAnswer                 `gorm:"foreignKey:TPRMAssessmentQuestionID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	CreatedByUser        *User                                 `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser        *User                                 `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by TprmAssessmentQuestion to `tprm_assessment_questions`
func (TPRMAssessmentSectionQuestion) TableName() string {
	return "tprm_assessment_section_questions"
}
