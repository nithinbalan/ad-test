package models

import "time"

type GovernanceReportMeeting struct {
	ID                 uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	GovernanceReportID uint      `json:"governance_report_id" gorm:"not null"`
	MeetingTitle       string    `json:"meeting_title" gorm:"type:nvarchar(255);not null"`
	MeetingDate        time.Time `json:"meeting_date" gorm:"type:date;not null"`
	ParticipantsCount  int       `json:"participants_count" gorm:"not null;default:0"`
	KeyDecisions       string    `json:"key_decisions,omitempty" gorm:"type:nvarchar(max)"`
	ActionsAssigned    string    `json:"actions_assigned,omitempty" gorm:"type:nvarchar(max)"`
	CreatedByID        uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	GovernanceReport *GovernanceReport `gorm:"foreignKey:GovernanceReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (GovernanceReportMeeting) TableName() string {
	return "governance_report_meetings"
}
