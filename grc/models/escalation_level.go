package models

import "time"

type EscalationLevel struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	EscalationID       uint       `json:"escalation_id"`
	NotificationMethod string     `json:"notification_method"`
	DayOffset          int        `json:"day_offset"`
	EscalationTo       string     `json:"escalation_to"`
	DesignationID      *uint      `json:"designation_id"`
	CreatedBYID        uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID        *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	Escalation    *Escalation  `gorm:"foreignKey:EscalationID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Designation   *Designation `gorm:"foreignKey:DesignationID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser *User        `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by EscalationLevel to `escalation_levels`
func (EscalationLevel) TableName() string {
	return "escalation_levels"
}
