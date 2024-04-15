package models

import (
	"time"
)

// StrategicProfileValue represents the values associated with a strategic profile
type StrategicProfileValue struct {
	ID                 uint       `json:"id" gorm:"primarykey;autoIncrement" binding:"-"`
	StrategicProfileID uint       `json:"strategic_profile_id" gorm:"not null" binding:"required"`
	TitleEn            string     `json:"title_en" gorm:"size:100;not null" binding:"required"`
	TitleAr            string     `json:"title_ar" gorm:"size:45;not null" binding:"required"`
	DescriptionEn      *string    `json:"description_en" gorm:"size:500;default:null" binding:"omitempty"`
	DescriptionAr      *string    `json:"description_ar" gorm:"size:500;default:null" binding:"omitempty"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null" binding:"-"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty" gorm:"type:datetime; default:null" binding:"omitempty"`

	// Relationship - ForeignKey in GORM should point to the field in the parent struct that this ID corresponds to.
	StrategicProfile StrategicProfile `gorm:"foreignKey:StrategicProfileID; constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (StrategicProfileValue) TableName() string {
	return "strategic_profile_values"
}
