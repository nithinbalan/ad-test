package models

import (
	"time"
)

type DocumentAuditLog struct {
	ID                     uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DocumentID             uint      `json:"document_id" binding:"required" gorm:"not null"`
	DocumentVersionID      uint      `json:"document_version_id" binding:"required" gorm:"not null"`
	UserID                 uint      `json:"user_id" binding:"required" gorm:"not null"`
	DocumentAuditLogTypeID uint      `json:"document_audit_log_type_id" binding:"required" gorm:"not null"`
	CreatedAt              time.Time `json:"created_at" gorm:"type:datetime;not null"`
	CreatedBYID            uint      `json:"created_by_id" gorm:"not null"`

	// Relations
	Document             *Document             `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentVersion      *DocumentVersion      `gorm:"foreignKey:DocumentVersionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User                 *User                 `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentAuditLogType *DocumentAuditLogType `gorm:"foreignKey:DocumentAuditLogTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser        *User                 `gorm:"foreignKey:CreatedBYID;"`
}

func (DocumentAuditLog) TableName() string {
	return "document_audit_logs"
}
