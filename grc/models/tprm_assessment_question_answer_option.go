package models

import "time"

type TPRMAssessmentQuestionAnswerOption struct {
	ID                          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TPRMAssessmentQuestionID    uint       `json:"tprm_question_id" gorm:"not null" binding:"required"`
	ThirdPartyQuestionLibraryID uint       `json:"third_party_question_library_id" gorm:"not null" binding:"required"`
	AnswerOptionTitleEn         *string    `json:"answer_option_title_en" gorm:"type:nvarchar(500);default:null"`
	AnswerOptionTitleAr         *string    `json:"answer_option_title_ar" gorm:"type:nvarchar(500);default:null"`
	RiskLevelID                 *uint      `json:"risk_level_id" gorm:"default:null"`
	RangeFrom                   *int       `json:"range_from" gorm:"default:null"`
	RangeTo                     *int       `json:"range_to" gorm:"default:null"`
	Score                       *int       `json:"score"`
	TPRMComplianceResponseID    *uint      `json:"tprm_compliance_response_id"`
	CreatedBYID                 *uint      `json:"created_by_id" gorm:"default:null"`
	UpdatedBYID                 *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                   *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	QuestionLibrary        *ThirdPartyQuestionLibrary `gorm:"foreignKey:ThirdPartyQuestionLibraryID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	RiskLevel              *RiskLevel                 `gorm:"foreignKey:RiskLevelID;constraint:OnUpdate:NO ACTION,OnDelete:SET NULL"`
	TPRMComplianceResponse *TPRMComplianceResponse    `gorm:"foreignKey:TPRMComplianceResponseID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	CreatedByUser          *User                      `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser          *User                      `gorm:"foreignKey:UpdatedBYID"`
}

func (TPRMAssessmentQuestionAnswerOption) TableName() string {
	return "tprm_assessment_question_answer_options"
}
