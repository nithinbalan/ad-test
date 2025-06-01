package models

import "time"

type UserJobDescriptionCompetency struct {
	ID                         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserJobDescriptionID       uint      `json:"user_job_description_id" gorm:"not null"`
	JobDescriptionCompetencyID uint      `json:"job_description_competency_id" gorm:"not null"`
	ActualLevelID              uint      `json:"actual_level_id" gorm:"not null"`
	Comments                   *string   `json:"comments" gorm:"type:nvarchar(MAX);default:null"`
	CreatedAt                  time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	UserJobDescription       *UserJobDescription              `gorm:"foreignKey:UserJobDescriptionID;constraint:OnDelete:CASCADE"`
	JobDescriptionCompetency *JobDescriptionCompetency        `gorm:"foreignKey:JobDescriptionCompetencyID;constraint:OnDelete:CASCADE"`
	ActualLevel              *CompetencyLevel                 `gorm:"foreignKey:ActualLevelID;constraint:OnDelete:RESTRICT"`
	Certificates             []*UserJobDescriptionCertificate `gorm:"foreignKey:UserJobDescriptionCompetencyID"`
}

func (UserJobDescriptionCompetency) TableName() string {
	return "user_job_description_competencies"
}
