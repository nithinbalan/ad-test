package models

import "time"

type IncidentPostReviewLesson struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentPostReviewID uint      `json:"incident_post_review_id" binding:"required" gorm:"not null"`
	Lesson               string    `json:"lesson" gorm:"type:nvarchar(500);not null"`
	CreatedByID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	IncidentPostReview IncidentPostReview `gorm:"foreignKey:IncidentPostReviewID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      User               `gorm:"foreignKey:CreatedByID"`
}

// TableName overrides the table name used by IncidentPostReviewLesson to `incident_post_review_lessons`
func (IncidentPostReviewLesson) TableName() string {
	return "incident_post_review_lessons"
}
