package models

import (
	"time"
)

type PolicyReport struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DocumentID  uint      `json:"document_id" gorm:"not null"`
	Title       string    `json:"title" gorm:"type:nvarchar(255);not null"`
	StartDate   time.Time `json:"start_date" gorm:"type:date;not null"`
	EndDate     time.Time `json:"end_date" gorm:"type:date;not null"`
	CreatedByID uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Document      *Document `gorm:"foreignKey:DocumentID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION;"`
	CreatedByUser *User     `gorm:"foreignKey:CreatedByID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION;"`
}

// TableName overrides the table name used by GORM
func (PolicyReport) TableName() string {
	return "policy_reports"
}
