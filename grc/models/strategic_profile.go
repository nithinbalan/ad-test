package models

import "time"

type StrategicProfile struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn      string     `json:"title_en" gorm:"size:255;not null" binding:"required"`
	TitleAr      string     `json:"title_ar" gorm:"size:255;default:null" binding:"omitempty"`
	Address      string     `json:"address" gorm:"type:varchar(MAX);not null" binding:"required"`
	Location     *string    `json:"location,omitempty" gorm:"size:255 default:null" binding:"omitempty"`
	Email        *string    `json:"email,omitempty" gorm:"size:100 email default:null" binding:"omitempty,email"`
	Industry     *string    `json:"industry,omitempty" gorm:"size:100 default:null" binding:"omitempty"`
	Phone        *string    `json:"phone,omitempty" gorm:"size:45 default:null" binding:"omitempty"`
	Linkedin     *string    `json:"linkedin,omitempty" gorm:"size:100 default:null" binding:"omitempty"`
	X            *string    `json:"x,omitempty" gorm:"size:100 default:null" binding:"omitempty"` // Consider a more descriptive field name
	Facebook     *string    `json:"facebook,omitempty" gorm:"size:100 default:null" binding:"omitempty"`
	Instagram    *string    `json:"instagram,omitempty" gorm:"size:100 default:null" binding:"omitempty"`
	Youtube      *string    `json:"youtube,omitempty" gorm:"size:100 default:null" binding:"omitempty"`
	AboutEn      *string    `json:"about_en,omitempty" gorm:"type:varchar(MAX) default:null" binding:"omitempty"`
	AboutAr      *string    `json:"about_ar,omitempty" gorm:"type:varchar(MAX)" binding:"omitempty"`
	VisionEn     *string    `json:"vision_en,omitempty" gorm:"size:500 default:null" binding:"omitempty,max=500"`
	VisionAr     *string    `json:"vision_ar,omitempty" gorm:"size:500 default:null" binding:"omitempty,max=500"`
	MissionEn    *string    `json:"mission_en,omitempty" gorm:"size:500 default:null" binding:"omitempty,max=500"`
	MissionAr    *string    `json:"mission_ar,omitempty" gorm:"size:500 default:null" binding:"omitempty,max=500"`
	ProfileImage *string    `json:"profile_image,omitempty" gorm:"size:500 default:null" binding:"omitempty,max=500"`
	CreatedAt    time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty" gorm:"type:datetime default:null" binding:"omitempty"`
	CreatedByID  uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID  *uint      `json:"updated_by_id,omitempty" gorm:"default:null" binding:"omitempty"`

	// Relationships
	CreatedByUser *User `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedByID"` // Pointer for nullable relationship

	StrategicProfileValues []StrategicProfileValue `json:"strategic_profile_values" gorm:"foreignKey:StrategicProfileID"`
}
