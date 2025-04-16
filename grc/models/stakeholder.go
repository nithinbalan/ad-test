package models

import "time"

type Stakeholder struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn     string     `json:"title_en" binding:"required"`
	TitleAr     string     `json:"title_ar" binding:"required"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null" binding:"omitempty"`
	CreatedBYID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID *uint      `json:"updated_by_id" gorm:"default:null" binding:"omitempty"`

	// Relationships
	CreatedByUser *User `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID"`
}

func (Stakeholder) TableName() string {
	return "stakeholders"
}
