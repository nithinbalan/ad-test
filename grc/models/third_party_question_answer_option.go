package models

import "time"

type ThirdPartyQuestionLibraryAnswerOption struct {
	ID                          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ThirdPartyQuestionLibraryID uint       `json:"third_party_question_library_id" gorm:"not null" binding:"required"`
	AnswerOptionTitleEn         *string    `json:"answer_option_title_en" gorm:"type:nvarchar(500);default:null"`
	AnswerOptionTitleAr         *string    `json:"answer_option_title_ar" gorm:"type:nvarchar(500);default:null"`
	RiskLevelID                 *uint      `json:"risk_level_id" gorm:"default:null"`
	RangeFrom                   *int       `json:"range_from" gorm:"default:null"`
	RangeTo                     *int       `json:"range_to" gorm:"default:null"`
	CreatedBYID                 *uint      `json:"created_by_id" gorm:"default:null"`
	UpdatedBYID                 *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                   *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	QuestionLibrary *ThirdPartyQuestionLibrary `gorm:"foreignKey:ThirdPartyQuestionLibraryID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	RiskLevel       *RiskLevel                 `gorm:"foreignKey:RiskLevelID;constraint:OnUpdate:NO ACTION,OnDelete:SET NULL"`
	CreatedByUser   *User                      `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser   *User                      `gorm:"foreignKey:UpdatedBYID"`
}

func (ThirdPartyQuestionLibraryAnswerOption) TableName() string {
	return "third_party_question_library_answer_options"
}
