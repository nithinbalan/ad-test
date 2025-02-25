package models

import "time"

type RiskComment struct {
	ID                 uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID     uint      `json:"risk_register_id" gorm:"not null;index"`
	RiskTimelineTypeID uint      `json:"risk_timeline_type_id" gorm:"not null"`
	Comments           string    `json:"comments" gorm:"type:text;not null"`
	RiskStatusID       uint      `json:"risk_status_id" gorm:"not null;index"`
	CreatedByID        uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt          time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (RiskComment) TableName() string {
	return "risk_comments"
}
