package models

import "time"

type DocumentReviewSupportDocument struct {
	ID               uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	DocumentReviewID uint       `json:"document_review_id" binding:"required" gorm:"not null"`
	URL              string     `json:"url" binding:"required" gorm:"type:varchar(1000);not null"`
	Name             string     `json:"name" binding:"required" gorm:"type:nvarchar(255);not null"`
	Size             int        `json:"size" binding:"required" gorm:"not null"`
	Token            string     `json:"token" binding:"required" gorm:"type:varchar(255);not null"`
	OriginalName     *string    `json:"original_name" gorm:"type:nvarchar(255);default:null"`
	Type             string     `json:"type" binding:"required" gorm:"type:varchar(255);not null"`
	CreatedAt        time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt        *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	DocumentReview *DocumentReview `gorm:"foreignKey:DocumentReviewID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (DocumentReviewSupportDocument) TableName() string {
	return "document_review_support_documents"
}
