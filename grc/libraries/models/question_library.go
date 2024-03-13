package models

type QuestionLibrary struct {
	ID              uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	QuestionGroupID uint   `json:"question_group_id" gorm:"not null"`
	QuestionEn      string `json:"question_en" binding:"required" gorm:"type:text;not null"`
	QuestionAr      string `json:"question_ar" binding:"required" gorm:"type:text;not null"`

	// Relationships
	QuestionGroup *QuestionGroup `gorm:"foreignKey:QuestionGroupID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
}

type QuestionLibraryWithGroupTitle struct {
	ID                   uint   `json:"id"`
	QuestionGroupID      uint   `json:"question_group_id"`
	QuestionEn           string `json:"question_en"`
	QuestionAr           string `json:"question_ar"`
	QuestionGroupTitleAr string `json:"question_group_title_ar"`
	QuestionGroupTitleEn string `json:"question_group_title_en"`
}

// TableName overrides the table name used by QuestionLibrary to `question_library`
func (QuestionLibrary) TableName() string {
	return "question_library"
}
