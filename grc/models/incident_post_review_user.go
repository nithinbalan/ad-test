package models

import "time"

type IncidentPostReviewUser struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentPostReviewID uint      `json:"incident_post_review_id" binding:"required" gorm:"not null"`
	UserID               uint      `json:"user_id" binding:"required" gorm:"not null"`
	CreatedByID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	IncidentPostReview IncidentPostReview `gorm:"foreignKey:IncidentPostReviewID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User               User               `gorm:"foreignKey:UserID"`
	CreatedByUser      User               `gorm:"foreignKey:CreatedByID"`
}

// TableName overrides the table name used by IncidentPostReviewUser to `incident_post_review_users`
func (IncidentPostReviewUser) TableName() string {
	return "incident_post_review_users"
}
