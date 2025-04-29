package models

type ThirdPartyQuestionType struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	Type    string `json:"type" binding:"required,max=45,oneof=yes-no dropdown-single dropdown-multiple rating-scale single-text text-area file-upload date-picker" gorm:"type:varchar(45);not null"`
	Status  string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active'"`
}

// TableName overrides the table name used by ThirdPartyQuestionType to `question_types`
func (ThirdPartyQuestionType) TableName() string {
	return "third_party_question_types"
}
