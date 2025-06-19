package models

import "time"

type JobDescriptionCompetencyCertification struct {
	ID                         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	JobDescriptionCompetencyID uint      `json:"job_description_competency_id" gorm:"not null"`
	Title                      string    `json:"title" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	IsMandatory                bool      `json:"is_mandatory" gorm:"type:bit;not null;default:0"`
	CreatedByID                uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                  time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	JobDescriptionCompetency *JobDescriptionCompetency `gorm:"foreignKey:JobDescriptionCompetencyID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser            *User                     `gorm:"foreignKey:CreatedByID"`
}

func (JobDescriptionCompetencyCertification) TableName() string {
	return "job_description_competency_certifications"
}
