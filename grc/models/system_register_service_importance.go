package models

import (
	"time"
)

type SystemRegisterServiceImportance struct {
	ID          uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn     string     `json:"title_en" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	TitleAr     string     `json:"title_ar" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	CreatedBYID uint       `json:"created_by_id"`
	UpdatedBYID *uint      `json:"updated_by_id"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	Status      string     `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active', 'inactive')"`

	// Relationships
	CreatedByUser *User `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by SystemRegisterServiceImportance to `system_register_service_importances`
func (SystemRegisterServiceImportance) TableName() string {
	return "system_register_service_importances"
}
