package models

import (
	"time"
)

type IncidentClosure struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentRegisterID uint       `json:"incident_register_id" binding:"required" gorm:"not null"`
	CompletionDate     time.Time  `json:"completion_date" binding:"required" gorm:"type:date;not null"`
	ResolutionSummary  *string    `json:"resolution_summary" gorm:"type:nvarchar(max);null"`
	Notes              *string    `json:"notes" gorm:"type:nvarchar(max);null"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime"`
	CreatedByID        uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBy          *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	IncidentRegister *IncidentRegister        `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Lessons          []*IncidentClosureLesson `gorm:"foreignKey:IncidentClosureID"`
	CreatedByUser    *User                    `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser    *User                    `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by IncidentClosure to `incident_closures`
func (IncidentClosure) TableName() string {
	return "incident_closures"
}
