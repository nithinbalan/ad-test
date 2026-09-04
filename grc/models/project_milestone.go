package models

import (
	"time"
)

type ProjectMilestone struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProjectID                uint       `json:"project_id" gorm:"not null"`
	TitleEn                  string     `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr                  string     `json:"title_ar" gorm:"type:nvarchar(255);not null"`
	DescriptionEn            *string    `json:"description_en" gorm:"type:nvarchar(max);null"`
	DescriptionAr            *string    `json:"description_ar" gorm:"type:nvarchar(max);null"`
	PlannedDate              *time.Time `json:"planned_date" gorm:"type:date;null"`
	ActualDate               *time.Time `json:"actual_date" gorm:"type:date;null"`
	CreatedBYID              uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;null"`
	ProjectMilestoneStatusID uint       `json:"project_milestone_status_id" gorm:"not null"`
	ProjectPhaseID           *uint      `json:"project_phase_id" gorm:"null"`

	// Relationships
	Project                *Project                `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ProjectMilestoneStatus *ProjectMilestoneStatus `gorm:"foreignKey:ProjectMilestoneStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ProjectPhase           *ProjectPhase           `gorm:"foreignKey:ProjectPhaseID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (ProjectMilestone) TableName() string {
	return "project_milestones"
}
