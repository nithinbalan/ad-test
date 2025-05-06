package models

import (
	"time"
)

type QuestionSet struct {
	ID          uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn     string     `json:"title_en" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	TitleAr     string     `json:"title_ar" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	Description *string    `json:"description" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	CreatedByUser *User                 `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser *User                 `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Questions     []QuestionSetQuestion `gorm:"foreignKey:QuestionSetID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
}

// TableName overrides the table name used by QuestionSet to `question_sets`
func (QuestionSet) TableName() string {
	return "question_sets"
}
