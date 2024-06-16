package models

import (
	"time"
)

type StrategicObjective struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn     string     `json:"title_en" gorm:"type:varchar(MAX);not null" binding:"required"`
	TitleAr     string     `json:"title_ar" gorm:"type:nvarchar(MAX);not null" binding:"required"`
	Status      string     `json:"status" gorm:"type:varchar(10);not null;default:'active'" binding:"required,oneof=active inactive"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null" binding:"omitempty"`
	CreatedByID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBy   *uint      `json:"updated_by_id" gorm:"default:null" binding:"omitempty"`

	// Relationships
	CreatedByUser *User `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by User to `strategic_objectives`
func (StrategicObjective) TableName() string {
	return "strategic_objectives"
}
