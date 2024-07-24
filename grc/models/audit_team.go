package models

import (
	"time"
)

type AuditTeam struct {
	ID          uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	TitleEn     string     `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr     string     `json:"title_ar" gorm:"type:nvarchar(255);not null"`
	TeamLeadID  uint       `json:"team_lead_id" gorm:"not null"`
	CreatedBYID uint       `json:"created_by_id"`
	UpdatedBYID *uint      `json:"updated_by_id"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	TeamLead      *User              `gorm:"foreignKey:TeamLeadID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Users         []*AuditTeamMember `gorm:"foreignKey:AuditTeamID"`
	CreatedByUser *User              `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User              `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by AuditTeam to `audit_teams`
func (AuditTeam) TableName() string {
	return "audit_teams"
}
