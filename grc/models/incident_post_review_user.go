package models

import "time"

type IncidentPostReviewUser struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentPostReviewID uint      `json:"incident_post_review_id" binding:"required" gorm:"not null"`
	UserID               uint      `json:"user_id" binding:"required" gorm:"not null"`
	CreatedBy            uint      `json:"created_by" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	IncidentPostReview IncidentPostReview `gorm:"foreignkey:IncidentPostReviewID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User               User               `gorm:"foreignkey:UserID"`
	CreatedByUser      User               `gorm:"foreignkey:CreatedBy"`
}

// TableName overrides the table name used by IncidentPostReviewUser to `incident_post_review_users`
func (IncidentPostReviewUser) TableName() string {
	return "incident_post_review_users"
}
