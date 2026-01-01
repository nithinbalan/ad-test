package models

import (
	"time"
)

type QuestionSetSection struct {
	ID                           uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	QuestionSetID                uint      `json:"question_set_id" binding:"required" gorm:"not null"`
	ThirdPartyQuestionCategoryID uint      `json:"third_party_question_category_id" binding:"required" gorm:"not null"`
	TitleEn                      string    `json:"title_en" binding:"required,max=1000" gorm:"type:nvarchar(1000);not null"`
	TitleAr                      string    `json:"title_ar" binding:"required,max=1000" gorm:"type:nvarchar(1000);not null"`
	Order                        int       `json:"order"`
	CreatedBYID                  uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                    time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	QuestionSet                *QuestionSet                 `gorm:"foreignKey:QuestionSetID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	ThirdPartyQuestionCategory *ThirdPartyQuestionCategory  `gorm:"foreignKey:ThirdPartyQuestionCategoryID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	Questions                  []QuestionSetSectionQuestion `gorm:"foreignKey:QuestionSetSectionID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	CreatedByUser              *User                        `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by QuestionSetQuestion to `question_set_questions`
func (QuestionSetSection) TableName() string {
	return "question_set_sections"
}
