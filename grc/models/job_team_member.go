package models

import (
	"time"
)

type JobTeamMember struct {
	ID            uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	JobTeamID     uint      `json:"job_team_id" gorm:"not null"`
	DesignationID uint      `json:"designation_id" gorm:"not null"`
	CreatedBYID   uint      `json:"created_by_id"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	JobTeam       *JobTeam     `gorm:"foreignKey:JobTeamID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Designation   *Designation `gorm:"foreignKey:DesignationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedBYID"`
}

// TableName overrides the table name used by JobTeamMember to `job_team_members`
func (JobTeamMember) TableName() string {
	return "job_team_members"
}
