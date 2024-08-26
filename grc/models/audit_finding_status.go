package models

import (
	"time"
)

type AuditFindingStatus struct {
	ID          uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn     string     `json:"title_en" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	TitleAr     string     `json:"title_ar" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	Type        string     `json:"type" binding:"required,max=10" gorm:"type:varchar(10);not null"`
	ColorCode   string     `json:"color_code" binding:"required,max=20" gorm:"type:varchar(20);not null"`
	Status      string     `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime"`
	CreatedByID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	CreatedByUser *User `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by AuditFindingStatus to `audit_finding_statuses`
func (AuditFindingStatus) TableName() string {
	return "audit_finding_statuses"
}
