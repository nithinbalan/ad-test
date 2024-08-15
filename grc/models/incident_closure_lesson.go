package models

import (
	"time"
)

type IncidentClosureLesson struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentClosureID uint       `json:"incident_closure_id" binding:"required" gorm:"not null"`
	Lessons           string     `json:"lessons" gorm:"type:nvarchar(500);not null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	IncidentClosure *IncidentClosure `gorm:"foreignKey:IncidentClosureID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by IncidentClosureLesson to `incident_closure_lessons`
func (IncidentClosureLesson) TableName() string {
	return "incident_closure_lessons"
}
