package models

import (
	"time"
)

type DocumentChangeRequest struct {
	ID                            uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DocumentsID                   uint       `json:"documents_id" gorm:"not null"`
	DocumentChangeRequestTypeID   uint       `json:"document_change_request_type_id" gorm:"not null"`
	ReasonForChange               *string    `json:"reason_for_change" gorm:"type:nvarchar(MAX)"`
	DocumentChangeRequestStatusID uint       `json:"document_change_request_statuses_id" gorm:"not null"`
	CreatedBYID                   uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                   *uint      `json:"updated_by_id"`
	CreatedAt                     time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                     *time.Time `json:"updated_at" gorm:"type:datetime"`

	//Relations
	Document                    *Document                    `gorm:"foreignKey:DocumentsID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentChangeRequestType   *DocumentChangeRequestType   `gorm:"foreignKey:DocumentChangeRequestTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentChangeRequestStatus *DocumentChangeRequestStatus `gorm:"foreignKey:DocumentChangeRequestStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser               *User                        `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser               *User                        `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Documents []*DocumentChangeRequestSupportingDocument `json:"documents" gorm:"foreignKey:DocumentChangeRequestID"`
}

func (DocumentChangeRequest) TableName() string {
	return "document_change_requests"
}
