package models

import "time"

type DocumentReview struct {
	ID               uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	DocumentID       uint       `json:"document_id" binding:"required" gorm:"not null"`
	ReviewDate       time.Time  `json:"review_date" binding:"required" gorm:"type:date;not null"`
	ISChangeRequired bool       `json:"is_change_required" gorm:"type:tinyint(1);not null;default:false"`
	Description      *string    `json:"description" gorm:"type:nvarchar(MAX);default:null"`
	CreatedBYID      uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt        time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedBYID      *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt        *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	Document      *Document `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser *User     `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser *User     `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	SupportDocuments      []*DocumentReviewSupportDocument `gorm:"foreignKey:DocumentReviewID"`
	DocumentChangeRequest *DocumentChangeRequest           `gorm:"foreignKey:DocumentReviewID"`
}

func (DocumentReview) TableName() string {
	return "document_reviews"
}
