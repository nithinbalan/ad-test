package models

import (
	"time"
)

type AuditScheduleStatusUpdateDocument struct {
	ID                          uint       `json:"id" gorm:"primary_key;auto_increment"`
	AuditScheduleStatusUpdateID uint       `json:"audit_schedule_status_update_id" gorm:"not null"`
	URL                         string     `json:"url" gorm:"type:varchar(1000)"`
	Name                        string     `json:"name"  binding:"required" gorm:"type:nvarchar(255)"`
	Size                        int        `json:"size"`
	Token                       string     `json:"token" gorm:"type:varchar(255)"`
	OriginalName                string     `json:"original_name" gorm:"type:nvarchar(255)"`
	Type                        string     `json:"type" gorm:"type:varchar(255)"`
	CreatedBYID                 uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                 *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                   *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	AuditScheduleStatusUpdate *AuditScheduleStatusUpdate `gorm:"foreignKey:AuditScheduleStatusUpdateID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	CreatedByUser             *User                      `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	UpdatedByUser             *User                      `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
}

// TableName overrides the table name used by AuditScheduleStatusUpdateDocument to `audit_schedule_status_update_documents`
func (AuditScheduleStatusUpdateDocument) TableName() string {
	return "audit_schedule_status_update_documents"
}
