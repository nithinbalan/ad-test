package models

import (
	"time"
)

type ProcessActivity struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ProcessID         uint       `json:"process_id" gorm:"not null" binding:"required"`
	TitleEn           string     `json:"title_en" gorm:"type:nvarchar(500);not null" binding:"required"`
	TitleAr           string     `json:"title_ar" gorm:"type:nvarchar(500);not null" binding:"required"`
	DescriptionEn     *string    `json:"description_en,omitempty" gorm:"type:nvarchar(MAX);default:null"`
	DescriptionAr     *string    `json:"description_ar,omitempty" gorm:"type:nvarchar(MAX);default:null"`
	Sequence          int        `json:"sequence" gorm:"not null"`
	EstimatedDuration *string    `json:"estimated_duration,omitempty" gorm:"type:nvarchar(100);default:null"`
	ExpectedOutput    string     `json:"expected_output" gorm:"type:nvarchar(MAX);not null"`
	CreatedBYID       uint       `json:"created_by_id" gorm:"not null" binding:"required"`
	UpdatedBYID       *uint      `json:"updated_by_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" gorm:"type:datetime"`

	// Relationships
	Process       *Process               `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser *User                  `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser *User                  `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Users         []*ProcessActivityUser `gorm:"foreignKey:ProcessActivityID"`
}

func (ProcessActivity) TableName() string {
	return "process_activities"
}
