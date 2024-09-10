package models

type AuditMeetingParticipant struct {
	ID             uint `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditMeetingID uint `json:"audit_meeting_id" binding:"required" gorm:"not null"`
	UserID         uint `json:"user_id" binding:"required" gorm:"not null"`

	// Relationships
	AuditMeeting *AuditMeeting `gorm:"foreignKey:AuditMeetingID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	User         *User         `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

// TableName overrides the default table name used by GORM
func (AuditMeetingParticipant) TableName() string {
	return "audit_meeting_participants"
}
