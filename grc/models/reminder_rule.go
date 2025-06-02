package models

import "time"

type ReminderRule struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn            string     `json:"title_en" gorm:"type:nvarchar(100);not null"`
	TitleAr            string     `json:"title_ar" gorm:"type:nvarchar(100);not null"`
	ReminderID         uint       `json:"reminder_id"`
	TriggerConditionID uint       `json:"trigger_condition_id"`
	NotificationMethod string     `json:"notification_method"`
	DayOffset          int        `json:"day_offset"`
	Type               string     `json:"type"`
	CreatedBYID        uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID        *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	Reminder         *Reminder         `gorm:"foreignKey:ReminderID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	TriggerCondition *TriggerCondition `gorm:"foreignKey:TriggerConditionID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser    *User             `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by ReminderRule to `reminder_rules`
func (ReminderRule) TableName() string {
	return "reminder_rules"
}
