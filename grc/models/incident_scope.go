package models

import (
	"time"
)

type IncidentScope struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn     string    `json:"title_en" binding:"required" gorm:"type:nvarchar(max);not null"`
	TitleAr     string    `json:"title_ar" binding:"required" gorm:"type:nvarchar(max);not null"`
	Type        string    `json:"type" binding:"required,max=25,oneof=single_system multiple_systems whole_organization" gorm:"type:varchar(25);not null"`
	Description *string   `json:"description" gorm:"type:nvarchar(max);null"`
	Status      string    `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
	CreatedBYID uint      `json:"created_by_id"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	CreatedByUser *User `gorm:"foreignKey:CreatedBYID"`
}

// TableName overrides the table name used by IncidentScope to `incident_scopes`
func (IncidentScope) TableName() string {
	return "incident_scopes"
}
