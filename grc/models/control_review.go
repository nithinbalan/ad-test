package models

import (
	"time"
)

type ControlReview struct {
	ID          uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ControlID   uint       `json:"control_id" gorm:"not null"`
	Comment     string     `json:"comment" binding:"required" gorm:"type:nvarchar(max);not null"`
	ReviewedAt  time.Time  `json:"reviewed_at" gorm:"type:datetime;not null"`
	ReviewedBy  uint       `json:"reviewed_by" gorm:"tnot null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	UpdatedByID *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	Control        *Control                 `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Documents      []*ControlReviewDocument `json:"documents" gorm:"foreignKey:ControlReviewID"`
	ReviewedByUser *User                    `gorm:"foreignKey:ReviewedBy"`
	UpdatedByUser  *User                    `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by ControlReview to `control_reviews`
func (ControlReview) TableName() string {
	return "control_reviews"
}
