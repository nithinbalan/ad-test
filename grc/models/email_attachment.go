package models

import "time"

type EmailAttachment struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	EmailID     uint      `json:"email_id" gorm:"not null;index"`
	StorageKey  string    `json:"storage_key" gorm:"type:nvarchar(512);not null"`
	FileName    string    `json:"file_name" gorm:"type:nvarchar(255);not null"`
	ContentType *string   `json:"content_type" gorm:"type:nvarchar(128);default:null"`
	SizeBytes   *int64    `json:"size_bytes" gorm:"default:null"`
	SortOrder   int       `json:"sort_order" gorm:"not null;default:0"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (EmailAttachment) TableName() string {
	return "email_attachments"
}
