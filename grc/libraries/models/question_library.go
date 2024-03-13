package models

type QuestionLibrary struct {
	ID              uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	QuestionGroupID uint   `json:"question_group_id" binding:"required" gorm:"not null"`
	QuestionEn      string `json:"question_en" binding:"required,max=255" gorm:"type:varchar(255);not null"`
	QuestionAr      string `json:"question_ar" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	Status          string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`

	// Relationships
	QuestionGroup *QuestionGroup `gorm:"foreignKey:QuestionGroupID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
}

type QuestionLibraryWithGroupTitle struct {
	ID                   uint   `json:"id"`
	QuestionGroupID      uint   `json:"question_group_id"`
	QuestionEn           string `json:"question_en"`
	QuestionAr           string `json:"question_ar"`
	Status               string `json:"status"`
	QuestionGroupTitleAr string `json:"question_group_title_ar"`
	QuestionGroupTitleEn string `json:"question_group_title_en"`
}

// TableName overrides the table name used by QuestionLibrary to `question_library`
func (QuestionLibrary) TableName() string {
	return "question_library"
}
