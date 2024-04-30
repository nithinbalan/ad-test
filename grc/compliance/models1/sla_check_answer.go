package models

import (
	"time"

	"jethur/grc/models"
)

type SLACheckAnswer struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	SLACheckQuestionID uint       `json:"sla_check_question_id" binding:"required" gorm:"not null"`
	ComplianceStatusID uint       `json:"compliance_status_id" binding:"required" gorm:"not null"`
	Comment            string     `json:"comment" binding:"required" gorm:"type:text;not null"`
	CreatedAt          time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`
	CreatedBy          uint       `json:"created_by" binding:"required" gorm:"not null"`
	UpdatedBy          *uint      `json:"updated_by,omitempty" gorm:"default:null"`

	// Relationships
	SLACheckQuestion *SLACheckQuestion `gorm:"foreignkey:SLACheckQuestionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceStatus *ComplianceStatus `gorm:"foreignKey:ComplianceStatusID"`
	CreatedByUser    *models.User      `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser    *models.User      `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by SLACheckAnswer to `sla_check_answers`
func (SLACheckAnswer) TableName() string {
	return "sla_check_answers"
}
