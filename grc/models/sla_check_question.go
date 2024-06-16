package models

import (
	"time"
)

type SLACheckQuestion struct {
	ID                      uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	SLACheckQuestionGroupID uint       `json:"sla_check_question_group_id" binding:"required" gorm:"not null"`
	TitleEn                 string     `json:"title_en" binding:"required,max=500" gorm:"type:varchar(500);not null"`
	TitleAr                 string     `json:"title_ar" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	CreatedAt               time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt               *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`
	CreatedByID             uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBy               *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`

	// Relationships
	Answer                *SLACheckAnswer        `json:"answer,omitempty" gorm:"foreignKey:SLACheckQuestionID"`
	SLACheckQuestionGroup *SLACheckQuestionGroup `gorm:"foreignKey:SLACheckQuestionGroupID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser         *User                  `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by SLACheckQuestion to `sla_check_questions`
func (SLACheckQuestion) TableName() string {
	return "sla_check_questions"
}
