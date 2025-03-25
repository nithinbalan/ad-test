package models

type SwotPestelQuestionGroup struct {
	ID       uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn  string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr  string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Order    int    `json:"order" gorm:"not null" binding:"required"`
	Status   string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
	IsSwot   bool   `json:"is_swot" gorm:"type:tinyint(1);null;default:0"`
	IsPestel bool   `json:"is_pestel" gorm:"type:tinyint(1);null;default:0"`
}

// TableName overrides the table name used by SwotQuestionGroup to `swot_question_groups`
func (SwotPestelQuestionGroup) TableName() string {
	return "swot_pestel_question_groups"
}
