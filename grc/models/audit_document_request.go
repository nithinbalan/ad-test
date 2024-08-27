package models

import (
	"time"
)

type AuditDocumentRequest struct {
	ID                           uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditID                      uint       `json:"audit_id" gorm:"not null"`
	AuditDocumentRequestStatusID uint       `json:"audit_document_request_status_id" gorm:"not null"`
	ResponsibleUnitID            uint       `json:"responsible_unit_id" binding:"required" gorm:"not null"`
	Title                        string     `json:"title" gorm:"type:nvarchar(255);not null"`
	Description                  *string    `json:"description" gorm:"type:nvarchar(MAX);default:null"`
	RequestedAt                  time.Time  `json:"requested_at" gorm:"type:datetime;not null"`
	CreatedAt                    time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                    *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID                  uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                  *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	Audit                      *Audit                               `gorm:"foreignKey:AuditID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditDocumentRequestStatus *AuditDocumentRequestStatus          `gorm:"foreignKey:AuditDocumentRequestStatusID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ResponsibleUnit            *OrganizationUnit                    `gorm:"foreignKey:ResponsibleUnitID"`
	AuditableItems             []*AuditDocumentRequestAuditableItem `gorm:"foreignKey:AuditDocumentRequestID"`
	CreatedByUser              *User                                `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser              *User                                `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by AuditDocumentRequest to `audit_document_requests`
func (AuditDocumentRequest) TableName() string {
	return "audit_document_requests"
}
