package models

import "time"

type ConsultationRequestComment struct {
	ID                          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	ConsultationRequestID       uint   `json:"consultation_request_id" gorm:"not null"`
	ConsultationRequestStatusID *uint  `json:"consultation_request_status_id" gorm:"not null"`
	Comment                     string `json:"comment" gorm:"type:nvarchar(MAX);not null"`

	// Audit Fields
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	CreatedByID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	UpdatedByID *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	ConsultationRequest       *ConsultationRequest       `gorm:"foreignKey:ConsultationRequestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ConsultationRequestStatus *ConsultationRequestStatus `gorm:"foreignKey:ConsultationRequestStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser             *User                      `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser             *User                      `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Documents                 []*ConsultationRequestCommentDocument `json:"documents" gorm:"foreignKey:ConsultationRequestCommentID"`
}

// TableName overrides the default table name for GORM
func (ConsultationRequestComment) TableName() string {
	return "consultation_request_comments"
}
