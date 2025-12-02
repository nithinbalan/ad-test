package models

import "time"

type FocusArea struct {
	ID            uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn       string     `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr       string     `json:"title_ar" gorm:"type:nvarchar(255);not null"`
	DescriptionEn *string    `json:"description_en" gorm:"type:nvarchar(max)"`
	DescriptionAr *string    `json:"description_ar" gorm:"type:nvarchar(max)"`
	CreatedBYID   uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID   *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt     time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	Status        string     `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active', 'inactive')"`

	CreatedByUser *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (FocusArea) TableName() string {
	return "focus_areas"
}
