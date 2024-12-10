package models

import (
	"time"
)

type DocumentChangeRequestInitiate struct {
	ID                                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DocumentID                            uint       `json:"document_id" gorm:"not null"`
	DocumentChangeRequestInitiateTypeID   uint       `json:"document_change_request_initiate_type_id" gorm:"not null"`
	ReasonForChange                       *string    `json:"reason_for_change" gorm:"type:nvarchar(MAX)"`
	DocumentChangeRequestInitiateStatusID uint       `json:"document_change_request_initiate_statuses_id" gorm:"not null"`
	CreatedBYID                           uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                           *uint      `json:"updated_by_id"`
	CreatedAt                             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                             *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relations
	Document                            *Document                            `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentChangeRequestInitiateType   *DocumentChangeRequestInitiateType   `gorm:"foreignKey:DocumentChangeRequestInitiateTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentChangeRequestInitiateStatus *DocumentChangeRequestInitiateStatus `gorm:"foreignKey:DocumentChangeRequestInitiateStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                       *User                                `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser                       *User                                `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Documents []*DocumentChangeRequestInitiateDocument `json:"documents" gorm:"foreignKey:DocumentChangeRequestInitiateID"`
}

func (DocumentChangeRequestInitiate) TableName() string {
	return "document_change_request_initiates"
}
