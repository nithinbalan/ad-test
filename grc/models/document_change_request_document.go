package models

import (
	"time"
)

type DocumentChangeRequestDocument struct {
	ID                      uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DocumentChangeRequestID uint       `json:"document_change_request_id" gorm:"not null"`
	URL                     string     `json:"url" gorm:"type:varchar(1000);not null"`
	Name                    string     `json:"name" gorm:"type:nvarchar(255);not null"`
	Size                    int        `json:"size" gorm:"not null"`
	Token                   string     `json:"token" gorm:"type:varchar(255);not null"`
	OriginalName            *string    `json:"original_name" gorm:"type:nvarchar(255);default:null"`
	Type                    string     `json:"type" gorm:"type:varchar(255);not null"`
	CreatedAt               time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt               *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	DocumentChangeRequest *DocumentChangeRequest `gorm:"foreignKey:DocumentChangeRequestID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (DocumentChangeRequestDocument) TableName() string {
	return "document_change_request_documents"
}
