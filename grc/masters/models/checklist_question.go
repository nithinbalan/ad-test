package models

type ChecklistQuestion struct {
	ID                uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ChecklistID       uint   `json:"checklist_id" binding:"required" gorm:"not null"`
	QuestionLibraryID uint   `json:"question_library_id,omitempty" binding:"omitempty" gorm:"default:null"`
	TitleEn           string `json:"title_en" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	TitleAr           string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Status            string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null"`

	// Relationships
	ChecklistID       *Checklist             `gorm:"foreignKey:ChecklistID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	QuestionLibraryID *model.QuestionLibrary `gorm:"foreignKey:QuestionLibraryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

// TableName overrides the table name used by ChecklistQuestion to `checklist_questions`
func (ChecklistQuestion) TableName() string {
	return "checklist_questions"
}
