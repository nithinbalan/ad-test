package models

import (
	"time"
)

type Project struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn           string     `json:"title_en" gorm:"type:varchar(100);not null"`
	TitleAr           string     `json:"title_ar" gorm:"type:varchar(100);not null"`
	DescriptionEn     *string    `json:"description_en" gorm:"type:text;null"`
	DescriptionAr     *string    `json:"description_ar" gorm:"type:text;null"`
	ResponsibleUserID uint       `json:"responsible_user_id" gorm:"not null"`
	StartDate         time.Time  `json:"start_date" gorm:"type:date;not null"`
	EndDate           time.Time  `json:"end_date" gorm:"type:date;not null"`
	ProjectTypeID     uint       `json:"project_type_id" gorm:"not null"`
	ProjectStatusID   uint       `json:"project_status_id" gorm:"not null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;null"`
	CreatedBYID       uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID       *uint      `json:"updated_by_id" gorm:"null"`

	// Relationships
	ResponsibleUser *User          `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ProjectType     *ProjectType   `gorm:"foreignKey:ProjectTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ProjectStatus   *ProjectStatus `gorm:"foreignKey:ProjectStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User          `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	UpdatedByUser   *User          `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName overrides the table name used by Project to `projects`
func (Project) TableName() string {
	return "projects"
}
