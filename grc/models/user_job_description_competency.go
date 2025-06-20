package models

import "time"

type UserJobDescriptionCompetency struct {
	ID                     uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserJobDescriptionID   uint      `json:"user_job_description_id" gorm:"not null"`
	CompetencyID           uint      `json:"competency_id" gorm:"not null"`
	ExpectedLevelID        *uint     `json:"expected_level_id" gorm:"default:null"`
	ActualLevelID          *uint     `json:"actual_level_id" gorm:"default:null"`
	Comments               *string   `json:"comments" gorm:"type:nvarchar(MAX);default:null"`
	CompetencyTitleEn      *string   `json:"competency_title_en" gorm:"type:nvarchar(255);default:null"`
	CompetencyTitleAr      *string   `json:"competency_title_ar" gorm:"type:nvarchar(255);default:null"`
	CompetencyGroupTitleEn *string   `json:"competency_group_title_en" gorm:"type:nvarchar(255);default:null"`
	CompetencyGroupTitleAr *string   `json:"competency_group_title_ar" gorm:"type:nvarchar(255);default:null"`
	CreatedAt              time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	UserJobDescription *UserJobDescription              `gorm:"foreignKey:UserJobDescriptionID;constraint:OnDelete:CASCADE"`
	Competency         *Competency                      `gorm:"foreignKey:CompetencyID;constraint:OnDelete:CASCADE"`
	ExpectedLevel      *CompetencyLevel                 `gorm:"foreignKey:ExpectedLevelID;constraint:OnDelete:RESTRICT"`
	ActualLevel        *CompetencyLevel                 `gorm:"foreignKey:ActualLevelID;constraint:OnDelete:RESTRICT"`
	Certificates       []*UserJobDescriptionCertificate `gorm:"foreignKey:UserJobDescriptionCompetencyID"`
}

func (UserJobDescriptionCompetency) TableName() string {
	return "user_job_description_competencies"
}
