package models

import (
	"time"
)

type SLACheckAnswer struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	SLACheckQuestionID uint       `json:"sla_check_question_id" binding:"required" gorm:"not null"`
	ComplianceStatusID uint       `json:"compliance_status_id" binding:"required" gorm:"not null"`
	Comment            string     `json:"comment" binding:"required" gorm:"type:text;not null"`
	CreatedAt          time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`
	CreatedByID        uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBy          *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`

	// Relationships
	SLACheckQuestion *SLACheckQuestion         `gorm:"foreignKey:SLACheckQuestionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceStatus *ComplianceStatus         `gorm:"foreignKey:ComplianceStatusID"`
	CreatedByUser    *User                     `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser    *User                     `gorm:"foreignKey:UpdatedBy"`
	Documents        []*SLACheckAnswerDocument `gorm:"foreignKey:SLACheckAnswerID"`
}

// TableName overrides the table name used by SLACheckAnswer to `sla_check_answers`
func (SLACheckAnswer) TableName() string {
	return "sla_check_answers"
}
