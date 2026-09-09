package models

import (
	"time"
)

type ProjectChangeRequestDocument struct {
	ID                     uint       `json:"id" gorm:"primary_key;auto_increment"`
	ProjectChangeRequestID uint       `json:"project_change_request_id" gorm:"not null"`
	URL                    string     `json:"url" gorm:"type:varchar(1000)"`
	Name                   string     `json:"name" binding:"required" gorm:"type:nvarchar(255);not null"`
	Size                   int        `json:"size" gorm:"not null"`
	Token                  string     `json:"token" gorm:"type:varchar(max);not null"`
	OriginalName           *string    `json:"original_name" gorm:"type:nvarchar(255)"`
	Type                   string     `json:"type" gorm:"type:varchar(255);not null"`
	CreatedByID            uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID            *uint      `json:"updated_by_id" gorm:"null"`
	CreatedAt              time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt              *time.Time `json:"updated_at" gorm:"type:datetime;null"`

	// Relationships
	ProjectChangeRequest *ProjectChangeRequest `gorm:"foreignKey:ProjectChangeRequestID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser        *User                 `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser        *User                 `gorm:"foreignKey:UpdatedByID"`
}

func (ProjectChangeRequestDocument) TableName() string {
	return "project_change_request_documents"
}
