package models

import "time"

type JobDescription struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn                  string     `json:"title_en" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	TitleAr                  string     `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	DescriptionEn            string     `json:"description_en" gorm:"type:nvarchar(max);"`
	DescriptionAr            string     `json:"description_ar" gorm:"type:nvarchar(max);"`
	RolesAndResponsibilities string     `json:"roles_and_responsibilities" gorm:"type:nvarchar(max);"`
	Status                   string     `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID              uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"default:null"`

	Competencies  []*JobDescriptionCompetency `gorm:"foreignKey:JobDescriptionID"`
	CreatedByUser *User                       `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User                       `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by JobDescription to `job_descriptions`
func (JobDescription) TableName() string {
	return "job_descriptions"
}
