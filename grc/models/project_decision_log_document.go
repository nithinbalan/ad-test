package models

import (
	"time"
)

type ProjectDecisionLogDocument struct {
	ID                   uint       `json:"id" gorm:"primary_key;auto_increment"`
	ProjectDecisionLogID uint       `json:"project_decision_log_id" gorm:"not null"`
	URL                  string     `json:"url" gorm:"type:varchar(255)"`
	Name                 string     `json:"name" binding:"required" gorm:"type:varchar(100);not null"`
	Size                 int        `json:"size" gorm:"not null"`
	Token                string     `json:"token" gorm:"type:varchar(100);not null"`
	OriginalName         *string    `json:"original_name" gorm:"type:nvarchar(255)"`
	Type                 string     `json:"type" gorm:"type:varchar(100);not null"`
	CreatedByID          uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID          *uint      `json:"updated_by_id" gorm:"null"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime;null"`

	// Relationships
	ProjectDecisionLog *ProjectDecisionLog `gorm:"foreignKey:ProjectDecisionLogID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser      *User               `gorm:"foreignKey:UpdatedByID"`
}

func (ProjectDecisionLogDocument) TableName() string {
	return "project_decision_log_documents"
}
