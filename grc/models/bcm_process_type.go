package models

import (
	"time"
)

type BCMProcessType struct {
	ID          uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn     string     `json:"title_en" binding:"required" gorm:"type:nvarchar(max);not null"`
	TitleAr     string     `json:"title_ar" binding:"required" gorm:"type:nvarchar(max);not null"`
	CreatedBYID uint       `json:"created_by_id"`
	UpdatedBYID *uint      `json:"updated_by_id"`
	CreatedAt   time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	Status      string     `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active', 'inactive')"`

	// Relationships
	CreatedByUser *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BCMProcessType) TableName() string {
	return "bcm_process_types"
}
