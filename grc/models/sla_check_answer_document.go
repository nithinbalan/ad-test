package models

import (
	"time"
)

type SLACheckAnswerDocument struct {
	ID               uint      `json:"id" gorm:"primary_key;auto_increment"`
	SLACheckAnswerID uint      `json:"sla_check_answer_id"`
	URL              string    `json:"url"`
	Name             string    `json:"name" binding:"required"`
	Size             int       `json:"size"`
	Token            string    `json:"token"`
	Type             string    `json:"type"`
	Icon             string    `json:"icon" binding:"required"`
	CreatedBy        uint      `json:"created_by" gorm:"not null"`
	UpdatedBy        uint      `json:"updated_by" gorm:"default:null"`
	CreatedAt        time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	SLACheckAnswer *SLACheckAnswer `gorm:"foreignKey:SLACheckAnswerID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser  *User           `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser  *User           `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by SLACheckAnswer to `sla_check_answers`
func (SLACheckAnswerDocument) TableName() string {
	return "sla_check_answer_documents"
}
