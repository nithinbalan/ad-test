package models

import "time"

type JobDescriptionCompetency struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	JobDescriptionID  uint       `json:"job_description_id" gorm:"not null"`
	CompetencyID      uint       `json:"competency_id" gorm:"not null"`
	CompetencyLevelID uint       `json:"competency_level_id" gorm:"not null"`
	CreatedByID       uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID       *uint      `json:"updated_by_id"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	JobDescription  *JobDescription                          `gorm:"foreignKey:JobDescriptionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Competency      *Competency                              `gorm:"foreignKey:CompetencyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CompetencyLevel *CompetencyLevel                         `gorm:"foreignKey:CompetencyLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User                                    `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser   *User                                    `gorm:"foreignKey:UpdatedByID"`
	Certifications  []*JobDescriptionCompetencyCertification `gorm:"foreignKey:JobDescriptionCompetencyID"`
}

func (JobDescriptionCompetency) TableName() string {
	return "job_description_competencies"
}
