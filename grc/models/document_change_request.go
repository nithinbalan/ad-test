package models

import (
	"time"
)

type DocumentChangeRequest struct {
	ID                              uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber                 string     `json:"reference_number" gorm:"type:varchar(45);not null"`
	DocumentID                      uint       `json:"document_id" gorm:"not null"`
	DocumentChangeRequestInitiateID *uint      `json:"document_change_request_initiate_id" gorm:"default:null"`
	DocumentChangeRequestTypeID     uint       `json:"document_change_request_type_id" gorm:"not null"`
	DocumentChangeRequestStatusID   uint       `json:"document_change_request_status_id" gorm:"not null"`
	ReasonForChange                 *string    `json:"reason_for_change" gorm:"type:nvarchar(max);default:null"`
	Version                         float64    `json:"version" gorm:"type:decimal(4,2);not null"`
	SubmittedByID                   *uint      `json:"submitted_by_id" gorm:"default:null"`
	SubmittedAt                     *time.Time `json:"submitted_at" gorm:"type:datetime;default:null"`
	PublishedByID                   *uint      `json:"published_by_id" gorm:"default:null"`
	PublishedAt                     *time.Time `json:"published_at" gorm:"type:datetime;default:null"`
	NextReviewLevel                 *int       `json:"next_review_level" gorm:"default:null"`
	IsLocked                        bool       `json:"is_locked" gorm:"type:tinyint(1);not null;default:0"`
	LastCheckoutAt                  *time.Time `json:"last_checkout_at" gorm:"type:date;default:null"`
	LastCheckoutByID                *uint      `json:"last_checkout_by_id" gorm:"default:null"`
	LastCheckinAt                   *time.Time `json:"last_checkin_at" gorm:"type:date;default:null"`
	LastCheckinByID                 *uint      `json:"last_checkin_by_id" gorm:"default:null"`
	CreatedBYID                     uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                     *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                       time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                       *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	Document                      *Document                      `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentChangeRequestInitiate *DocumentChangeRequestInitiate `gorm:"foreignKey:DocumentChangeRequestInitiateID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentChangeRequestType     *DocumentChangeRequestType     `gorm:"foreignKey:DocumentChangeRequestTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentChangeRequestStatus   *DocumentChangeRequestStatus   `gorm:"foreignKey:DocumentChangeRequestStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SubmittedByUser               *User                          `gorm:"foreignKey:SubmittedByID"`
	PublishedByUser               *User                          `gorm:"foreignKey:PublishedByID"`
	LastCheckoutByUser            *User                          `gorm:"foreignKey:LastCheckoutByID"`
	LastCheckinByUser             *User                          `gorm:"foreignKey:LastCheckinByID"`
	CreatedByUser                 *User                          `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser                 *User                          `gorm:"foreignKey:UpdatedBYID"`

	ChangeRequestDocument *DocumentChangeRequestDocument `json:"change_request_document" gorm:"foreignKey:DocumentChangeRequestID"`
}

func (DocumentChangeRequest) TableName() string {
	return "document_change_requests"
}
