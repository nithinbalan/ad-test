package models

import "time"

type MeetingStrategicObjective struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	MeetingID   uint      `json:"strategic_objective_id" gorm:"not null"`
	ProjecctID  uint      `json:"project_id" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`
	CreatedBYID uint      `json:"created_by_id" gorm:"not null"`

	// Relationships
	StrategicObjective *StrategicObjective `gorm:"foreignKey:StrategicObjectiveID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Meeting            *Meeting            `gorm:"foreignKey:MeetingID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBYID"`
}

// TableName overrides the default table name
func (MeetingStrategicObjective) TableName() string {
	return "meeting_strategic_objectives"
}
