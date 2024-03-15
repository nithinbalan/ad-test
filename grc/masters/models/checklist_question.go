package models

import (
	libraryModel "jethur/grc/libraries/models"
)

type ChecklistQuestion struct {
	ID                uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ChecklistID       uint   `json:"checklist_id" binding:"required" gorm:"not null"`
	QuestionLibraryID uint   `json:"question_library_id,omitempty" binding:"omitempty" gorm:"default:null"`
	QuestionEn        string `json:"question_en" binding:"required,max=500" gorm:"type:varchar(500);not null"`
	QuestionAr        string `json:"question_ar" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`

	// Relationships
	Checklist       *Checklist                    `gorm:"foreignKey:ChecklistID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	QuestionLibrary *libraryModel.QuestionLibrary `gorm:"foreignKey:QuestionLibraryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type ChecklistQuestionWithChecklistTitle struct {
	ID                uint   `json:"id"`
	ChecklistID       uint   `json:"checklist_id"`
	QuestionLibraryID string `json:"question_library_id"`
	QuestionEn        string `json:"question_en"`
	QuestionAr        string `json:"question_ar"`
	ChecklistTitleAr  string `json:"checklist_title_ar"`
	ChecklistTitleEn  string `json:"checklist_title_en"`
}

// TableName overrides the table name used by ChecklistQuestion to `checklist_questions`
func (ChecklistQuestion) TableName() string {
	return "checklist_questions"
}
