package models

import "time"

type UserJobDescriptionCertificate struct {
	ID                                      uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserJobDescriptionCompetencyID          uint      `json:"user_job_description_competency_id" gorm:"not null"`
	JobDescriptionCompetencyCertificationID uint      `json:"job_description_competency_certification_id" gorm:"not null"`
	CertificateTitle                        string    `json:"certificate_title" gorm:"type:nvarchar(255);not null"`
	URL                                     *string   `json:"url" gorm:"type:varchar(1000);default:null"`
	Name                                    *string   `json:"name" gorm:"type:nvarchar(255);default:null"`
	Size                                    *int      `json:"size" gorm:"default:null"`
	Token                                   *string   `json:"token" gorm:"type:varchar(255);default:null"`
	OriginalName                            *string   `json:"original_name" gorm:"type:nvarchar(255);default:null"`
	Type                                    *string   `json:"type" gorm:"type:varchar(255);default:null"`
	CreatedByID                             *uint     `json:"created_by_id" gorm:"default:null"`
	CreatedAt                               time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	UserJobDescriptionCompetency *UserJobDescriptionCompetency          `gorm:"foreignKey:UserJobDescriptionCompetencyID;constraint:OnDelete:CASCADE"`
	JobDescriptionCertification  *JobDescriptionCompetencyCertification `gorm:"foreignKey:JobDescriptionCompetencyCertificationID;constraint:OnDelete:CASCADE"`
	CreatedByUser                *User                                  `gorm:"foreignKey:CreatedByID"`
}

func (UserJobDescriptionCertificate) TableName() string {
	return "user_job_description_certificates"
}
