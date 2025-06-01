package models

import "time"

type UserJobDescription struct {
	ID               uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID           uint      `json:"user_id" gorm:"not null"`
	JobDescriptionID uint      `json:"job_description_id" gorm:"not null"`
	CreatedByID      uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt        time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	User           *User                           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	JobDescription *JobDescription                 `gorm:"foreignKey:JobDescriptionID;constraint:OnDelete:CASCADE"`
	CreatedByUser  *User                           `gorm:"foreignKey:CreatedByID"`
	Competencies   []*UserJobDescriptionCompetency `gorm:"foreignKey:UserJobDescriptionID"`
}

func (UserJobDescription) TableName() string {
	return "user_job_descriptions"
}
