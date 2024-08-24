package models

import (
	"time"
)

type AuditWorkpaper struct {
	ID           uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditID      uint       `json:"audit_id" gorm:"not null"`
	SourceID     uint       `json:"source_id" gorm:"not null"`
	Title        string     `json:"title" gorm:"type:nvarchar(255);not null"`
	Description  *string    `json:"description" gorm:"type:nvarchar(MAX);default:null"`
	URL          *string    `json:"url" gorm:"type:varchar(1000);default:null"`
	Name         *string    `json:"name" gorm:"type:nvarchar(255);default:null"`
	Size         *int       `json:"size" gorm:"default:null"`
	Token        *string    `json:"token" gorm:"type:varchar(255);default:null"`
	OriginalName *string    `json:"original_name" gorm:"type:nvarchar(255);default:null"`
	CreatedAt    time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt    *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID  uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID  *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	Audit         *Audit  `gorm:"foreignKey:AuditID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	Source        *Source `gorm:"foreignKey:SourceID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser *User   `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User   `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by AuditWorkpaper to `audit_workpapers`
func (AuditWorkpaper) TableName() string {
	return "audit_workpapers"
}
