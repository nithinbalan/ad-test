package models

import (
	"time"
)

type AuditTeamMember struct {
	ID          uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditTeamID uint      `json:"audit_team_id" gorm:"not null"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	CreatedBYID uint      `json:"created_by_id"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditTeam     *AuditTeam `gorm:"foreignKey:AuditTeamID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User          *User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User      `gorm:"foreignKey:CreatedBYID"`
}

// TableName overrides the table name used by AuditTeamMember to `audit_team_members`
func (AuditTeamMember) TableName() string {
	return "audit_team_members"
}
