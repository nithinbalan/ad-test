package models

import "time"

type ThirdPartyQuestionLibrary struct {
	ID                           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ReferenceNumber              string     `json:"reference_number" gorm:"type:varchar(50);not null;unique" binding:"required"`
	QuestionTitleEn              string     `json:"question_title_en" gorm:"type:nvarchar(1000);not null" binding:"required"`
	QuestionTitleAr              string     `json:"question_title_ar" gorm:"type:nvarchar(1000);not null" binding:"required"`
	ThirdPartyQuestionCategoryID uint       `json:"third_party_question_category_id" gorm:"not null" binding:"required"`
	Mandatory                    bool       `json:"mandatory" gorm:"not null"`
	ThirdPartyQuestionTypeID     uint       `json:"third_party_question_type_id" gorm:"not null" binding:"required"`
	CreatedBYID                  *uint      `json:"created_by_id" gorm:"default:null"`
	UpdatedBYID                  *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                    time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                    *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	QuestionCategory *ThirdPartyQuestionCategory              `gorm:"foreignKey:ThirdPartyQuestionCategoryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	QuestionType     *ThirdPartyQuestionType                  `gorm:"foreignKey:ThirdPartyQuestionTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser    *User                                    `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser    *User                                    `gorm:"foreignKey:UpdatedBYID"`
	AnswerOptions    []*ThirdPartyQuestionLibraryAnswerOption `gorm:"foreignKey:ThirdPartyQuestionLibraryID"`
}

func (ThirdPartyQuestionLibrary) TableName() string {
	return "third_party_question_libraries"
}
