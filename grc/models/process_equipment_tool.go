package models

import (
	"time"
)

type ProcessEquipmentTool struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn     string     `json:"title_en" gorm:"type:nvarchar(500);not null"`
	TitleAr     string     `json:"title_ar" gorm:"type:nvarchar(500);not null"`
	Type        string     `json:"type" gorm:"type:varchar(20);not null;default:'other'"`
	Status      string     `json:"status" gorm:"type:varchar(10);not null;default:'active'"`
	CreatedBYID *uint      `json:"created_by_id" gorm:"default:null"`
	UpdatedBYID *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	CreatedByUser *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProcessEquipmentTool) TableName() string {
	return "process_equipment_tools"
}
