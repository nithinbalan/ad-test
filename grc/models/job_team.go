package models

import (
	"time"
)

type JobTeam struct {
	ID                    uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	TitleEn               string     `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr               string     `json:"title_ar" gorm:"type:nvarchar(255);not null"`
	TeamLeadDesignationID *uint      `json:"team_lead_designation_id" gorm:"default:null"`
	TeamLeadUserID        *uint      `json:"team_lead_user_id" gorm:"default:null"`
	Scope                 *string    `json:"scope" gorm:"type:nvarchar(max)"`
	Objective             *string    `json:"objective" gorm:"type:nvarchar(max)"`
	CreatedBYID           uint       `json:"created_by_id"`
	UpdatedBYID           *uint      `json:"updated_by_id"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	TeamLeadDesignation *Designation     `gorm:"foreignKey:TeamLeadDesignationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	TeamLeadUser        *User            `gorm:"foreignKey:TeamLeadUserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	JobTeamMembers      []*JobTeamMember `gorm:"foreignKey:JobTeamID"`
	CreatedByUser       *User            `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser       *User            `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by JobTeam to `job_teams`
func (JobTeam) TableName() string {
	return "job_teams"
}
