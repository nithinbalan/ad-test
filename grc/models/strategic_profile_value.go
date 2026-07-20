package models

import (
	"time"
)

// StrategicProfileValue represents the values associated with a strategic profile
type StrategicProfileValue struct {
	ID                 uint       `json:"id" gorm:"primarykey;autoIncrement" binding:"-"`
	StrategicProfileID uint       `json:"strategic_profile_id" gorm:"not null" binding:"required"`
	TitleEn            string     `json:"title_en" gorm:"type:nvarchar(max);not null" binding:"required"`
	TitleAr            string     `json:"title_ar" gorm:"type:nvarchar(max);not null" binding:"required"`
	DescriptionEn      *string    `json:"description_en" gorm:"type:nvarchar(max);default:null" binding:"omitempty"`
	DescriptionAr      *string    `json:"description_ar" gorm:"type:nvarchar(max);default:null" binding:"omitempty"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null" binding:"-"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null" binding:"omitempty"`

	StrategicProfile StrategicProfile `gorm:"foreignKey:StrategicProfileID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (StrategicProfileValue) TableName() string {
	return "strategic_profile_values"
}
