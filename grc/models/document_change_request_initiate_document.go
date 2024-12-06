package models

import (
	"time"
)

type DocumentChangeRequestInitiateDocument struct {
	ID                              uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DocumentChangeRequestInitiateID uint       `json:"document_change_request_initiate_id" gorm:"not null"`
	URL                             string     `json:"url" gorm:"type:varchar(45)"`
	Name                            string     `json:"name" gorm:"type:varchar(100);not null"`
	Size                            int        `json:"size" gorm:"not null"`
	Token                           string     `json:"token" gorm:"type:varchar(100);not null"`
	OriginalName                    *string    `json:"original_name" binding:"omitempty" gorm:"type:nvarchar(255)"`
	Type                            string     `json:"type" gorm:"type:varchar(100);not null"`
	CreatedAt                       time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                       *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relations
	DocumentChangeRequestInitiate *DocumentChangeRequestInitiate `gorm:"foreignKey:DocumentChangeRequestInitiateID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (DocumentChangeRequestInitiateDocument) TableName() string {
	return "document_change_request_initiate_documents"
}
