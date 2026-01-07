package models

import "time"

type IncidentReport struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title  string `json:"title" binding:"required,max=255"`

	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`

	CreatedBYID uint      `json:"created_by_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	SubmittedBYID *uint      `json:"submitted_by_id"`
	SubmittedAt   *time.Time `json:"submitted_at" gorm:"type:datetime"`
	
	//relationships
	CreatedByUser   *User                  `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SubmittedByUser *User                  `gorm:"foreignKey:SubmittedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (IncidentReport) TableName() string {
	return "incident_reports"
}