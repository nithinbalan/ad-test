package models

import "time"

type IncidentPostReviewRecommendation struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentPostReviewID uint      `json:"incident_post_review_id" binding:"required" gorm:"not null"`
	Recommendation       string    `json:"recommendation" gorm:"type:nvarchar(500);not null"`
	CreatedByID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	IncidentPostReview IncidentPostReview `gorm:"foreignKey:IncidentPostReviewID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      User               `gorm:"foreignKey:CreatedByID"`
}

// TableName overrides the table name used by IncidentPostReviewRecommendation to `incident_post_review_recommendations`
func (IncidentPostReviewRecommendation) TableName() string {
	return "incident_post_review_recommendations"
}
